package server

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
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

	server := &Server{
		Application:   app,
		Configuration: cfg,
	}

	err := server.assignJWTFiles()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) Run() error {
	addr := net.JoinHostPort(s.Configuration.HTTP.Address, strconv.Itoa(s.Configuration.HTTP.Port))

	if s.Configuration.HTTP.SSL {
		log.Info("SSL enabled.")

		return s.Application.Run(iris.Raw(func() error {
			return http3.ListenAndServe(addr, s.Configuration.HTTP.SSLCert, s.Configuration.HTTP.SSLCert, s.Application)
		}), iris.WithOptimizations)
	}

	log.Warn("SSL disabled, this is not recommended.")

	return s.Application.Run(iris.Addr(addr))
}
