package student

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"net/http"
	"time"
)

type Provider interface {
	StudentData(ctx *gin.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *gin.Engine, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/student")
	party.Use(authService.AuthMiddleware.MiddlewareFunc())
	{
		if storage != nil && storage[0] != nil {
			party.GET("/", cache.CachePage(storage[0], time.Minute, func(c *gin.Context) {
				h.StudentData(c)
			}))
		} else {
			party.GET("/", h.StudentData)
		}
	}

}

// StudentData
// @Description Get student information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.User	"Successfully retrieved student information"
// @Failure 400 {string} string "Error getting student information."
// @Router /student [get]
func (h Handler) StudentData(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	studentData, err := h.Service.StudentData(claims)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		log.Debug("Error getting student data", "error", err)
		return
	}

	ctx.JSON(http.StatusOK, studentData)
}
