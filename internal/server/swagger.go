package server

import (
	_ "github.com/imhinotori/duoc-plus/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) registerSwagger() {
	s.Application.GET("/swagger/*", echoSwagger.WrapHandler)
}
