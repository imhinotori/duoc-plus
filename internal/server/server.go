package server

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/attendance"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/imhinotori/duoc-plus/internal/grades"
	"github.com/imhinotori/duoc-plus/internal/schedule"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/quic-go/quic-go/http3"
	"net"
	"strconv"
)

type Server struct {
	Application   *iris.Application
	Configuration *config.Config
	JWTSigner     *jwt.Signer
	JWTVerifier   *jwt.Verifier
	Duoc          *duoc.Client
	Handlers      *Handlers
}

type Handlers struct {
	authHandler       auth.Handler
	attendanceHandler attendance.Handler
	scheduleHandler   schedule.Handler
	gradesHandler     grades.Handler
}

func New(cfgOpts ...config.Option) (*Server, error) {
	cfg := config.Build(cfgOpts...)

	app := iris.Default()

	if cfg.General.Debug {
		log.SetLevel(log.DebugLevel)
		app.Logger().SetLevel("debug")
		log.Debug("Debug mode enabled")
	} else {
		log.SetLevel(log.InfoLevel)
		app.Logger().SetLevel("info")
	}

	httpClient, err := duoc.NewHost(cfg)

	server := &Server{
		Application:   app,
		Configuration: cfg,
		Duoc:          httpClient,
		Handlers:      &Handlers{},
	}

	err = server.assignJWTFiles()
	if err != nil {
		return nil, err
	}

	server.Handlers.authHandler = auth.Handler{Service: auth.New(cfg, server.JWTSigner, server.JWTVerifier, server.Duoc)}
	server.Handlers.attendanceHandler = attendance.Handler{Service: attendance.New(cfg, server.Duoc)}
	server.Handlers.scheduleHandler = schedule.Handler{Service: schedule.New(cfg, server.Duoc)}
	server.Handlers.gradesHandler = grades.Handler{Service: grades.New(cfg, server.Duoc)}

	return server, nil
}

func (s *Server) Run() error {
	addr := net.JoinHostPort(s.Configuration.HTTP.Address, strconv.Itoa(s.Configuration.HTTP.Port))

	verifyMiddleware := s.JWTVerifier.Verify(func() interface{} {
		return new(auth.Claims)
	})

	log.Warn("SSL disabled, this is not recommended.")

	s.Handlers.authHandler.Start(s.Application, verifyMiddleware)
	s.Handlers.attendanceHandler.Start(s.Application, verifyMiddleware)
	s.Handlers.scheduleHandler.Start(s.Application, verifyMiddleware)
	s.Handlers.gradesHandler.Start(s.Application, verifyMiddleware)

	handleSwagger(s)

	if s.Configuration.HTTP.SSL {
		log.Info("SSL enabled.")

		return s.Application.Run(iris.Raw(func() error {
			return http3.ListenAndServe(addr, s.Configuration.HTTP.SSLCert, s.Configuration.HTTP.SSLKey, s.Application)
		}), iris.WithOptimizations)
	}

	return s.Application.Run(iris.Addr(addr))
}
