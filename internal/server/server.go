package server

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/attendance"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/database"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/imhinotori/duoc-plus/internal/grades"
	"github.com/imhinotori/duoc-plus/internal/schedule"
	"github.com/imhinotori/duoc-plus/internal/student"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"strconv"
)

type Server struct {
	Application   *echo.Echo
	Configuration *config.Config
	Duoc          *duoc.Client
	Database      *database.Database
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
	app := echo.New()

	if cfg.General.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	httpClient, err := duoc.NewHost(cfg)
	if err != nil {
		return nil, err
	}

	db, err := database.New(cfg)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Application:   app,
		Configuration: cfg,
		Duoc:          httpClient,
		Database:      db,
		Handlers:      &Handlers{},
	}

	server.Application.Use(middleware.Logger())
	server.Application.Use(middleware.Recover())

	server.registerSwagger()

	server.Handlers.authHandler = auth.Handler{Service: auth.New(cfg, db, server.Duoc)}
	server.Handlers.attendanceHandler = attendance.Handler{Service: attendance.New(cfg, db, server.Duoc)}
	server.Handlers.scheduleHandler = schedule.Handler{Service: schedule.New(cfg, db, server.Duoc)}
	server.Handlers.gradesHandler = grades.Handler{Service: grades.New(cfg, db, server.Duoc)}
	server.Handlers.studentHandler = student.Handler{Service: student.New(cfg, db, server.Duoc)}

	server.Application.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	return server, nil
}

func (s *Server) Run() error {
	addr := net.JoinHostPort(s.Configuration.HTTP.Address, strconv.Itoa(s.Configuration.HTTP.Port))

	s.Handlers.authHandler.Start(s.Application)
	s.Handlers.attendanceHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.scheduleHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.gradesHandler.Start(s.Application, s.Handlers.authHandler.Service)
	s.Handlers.studentHandler.Start(s.Application, s.Handlers.authHandler.Service)

	//handleSwagger(s)

	if s.Configuration.HTTP.SSL {
		return s.Application.StartTLS(addr, s.Configuration.HTTP.SSLCert, s.Configuration.HTTP.SSLKey)
	}
	log.Warn("SSL disabled, this is not recommended.")
	return s.Application.Start(addr)
}
