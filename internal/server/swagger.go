package server

import (
	_ "github.com/imhinotori/duoc-plus/docs"
	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/swaggerFiles"
)

func handleSwagger(s *Server) {
	swaggerUI := swagger.Handler(swaggerFiles.Handler,
		swagger.URL("/swagger/swagger.json"),
		swagger.DeepLinking(true),
		swagger.Prefix("/swagger"),
	)

	s.Application.Get("/swagger", swaggerUI)
	s.Application.Get("/swagger/{any:path}", swaggerUI)
}
