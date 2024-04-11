package server

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/imhinotori/duoc-plus/internal/attendance"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/imhinotori/duoc-plus/internal/grades"
	"github.com/imhinotori/duoc-plus/internal/schedule"
	"github.com/imhinotori/duoc-plus/internal/student"
	"net"
	"strconv"
)

type Server struct {
	Application   *gin.Engine
	Configuration *config.Config
	Duoc          *duoc.Client
	Handlers      *Handlers
}

type Handlers struct {
	authHandler       auth.Handler
	attendanceHandler attendance.Handler
	scheduleHandler   schedule.Handler
	gradesHandler     grades.Handler
	studentHandler    student.Handler
}

func New(cfg *config.Config) (*Server, error) {
	app := gin.Default()

	if cfg.General.Debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.Debug("Debug mode enabled")
	} else {
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	httpClient, err := duoc.NewHost(cfg)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Application:   app,
		Configuration: cfg,
		Duoc:          httpClient,
		Handlers:      &Handlers{},
	}

	server.Handlers.authHandler = auth.Handler{Service: auth.New(cfg, server.Duoc)}
	server.Handlers.attendanceHandler = attendance.Handler{Service: attendance.New(cfg, server.Duoc)}
	server.Handlers.scheduleHandler = schedule.Handler{Service: schedule.New(cfg, server.Duoc)}
	server.Handlers.gradesHandler = grades.Handler{Service: grades.New(cfg, server.Duoc)}
	server.Handlers.studentHandler = student.Handler{Service: student.New(cfg, server.Duoc)}

	return server, nil
}

func (s *Server) Run() error {
	addr := net.JoinHostPort(s.Configuration.HTTP.Address, strconv.Itoa(s.Configuration.HTTP.Port))

	log.Warn("SSL disabled, this is not recommended.")

	s.Handlers.authHandler.Start(s.Application)
	s.Handlers.attendanceHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.scheduleHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.gradesHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.studentHandler.Start(s.Application, s.Handlers.authHandler.Service)

	//handleSwagger(s)

	if s.Configuration.HTTP.SSL {
		log.Info("SSL enabled.")

		return s.Application.RunTLS(addr, s.Configuration.HTTP.SSLCert, s.Configuration.HTTP.SSLKey)
	}

	return s.Application.Run(addr)
}
