package server

import (
	"github.com/imhinotori/duoc-plus/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (s *Server) registerSwagger() {
	docs.SwaggerInfo.BasePath = "/"

	s.Application.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
