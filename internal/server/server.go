package server

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/kataras/iris/v12"
	"github.com/quic-go/quic-go/http3"
	"net"
	"strconv"
)

type Server struct {
	Application   *iris.Application
	Configuration *config.Config
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
	app.RegisterView(iris.Blocks("./web/views", ".html").Reload(cfg.General.Debug))

	server := &Server{
		Application:   app,
		Configuration: cfg,
	}

	return server, nil
}

func (s *Server) Run() error {
	addr := net.JoinHostPort(s.Configuration.HTTP.Address, strconv.Itoa(s.Configuration.HTTP.Port))

	if s.Configuration.HTTP.SSL {
		log.Info("SSL enabled.")

		return s.Application.Run(iris.Raw(func() error {
			return http3.ListenAndServe(addr, s.Configuration.JWT.PublicKey, s.Configuration.JWT.PrivateKey, s.Application)
		}), iris.WithOptimizations)
	}

	log.Warn("SSL disabled, this is not recommended.")

	return s.Application.Run(iris.Addr(addr))
}
