package server

import (
	"github.com/imhinotori/duoc-plus/docs"
)

func (s *Server) registerSwagger() {
	docs.SwaggerInfo.BasePath = "/"

	//s.Application.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
