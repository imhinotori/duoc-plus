package attendance

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
	Attendance(ctx *gin.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *gin.Engine, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/attendance")
	party.Use(authService.AuthMiddleware.MiddlewareFunc())
	{
		if storage != nil && storage[0] != nil {
			party.GET("/", cache.CachePage(storage[0], time.Minute, func(c *gin.Context) {
				h.Attendance(c)
			}))
		} else {
			party.GET("/", h.Attendance)
		}
	}

}

// Attendance
// @Description Get user attendance
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.Attendance	"Successfully retrieved attendance"
// @Failure 400 {string} string "Error getting attendance."
// @Router /attendance [get]
func (h Handler) Attendance(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	attendance, err := h.Service.Attendance(claims)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		log.Debug("Error getting attendance", "error", err)
		return
	}

	ctx.JSON(http.StatusOK, attendance)
}
