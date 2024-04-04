package schedule

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"net/http"
	"time"
)

type Provider interface {
	Schedule(ctx *gin.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *gin.Engine, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/schedule")
	party.Use(authService.AuthMiddleware.MiddlewareFunc())
	{
		if storage != nil && storage[0] != nil {
			party.GET("/", cache.CachePage(storage[0], time.Minute, func(c *gin.Context) {
				h.Schedule(c)
			}))
		} else {
			party.GET("/", h.Schedule)
		}
	}

}

// Schedule
// @Description Get user schedule
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.CareerSchedule	"Successfully retrieved schedule"
// @Failure 400 {string} string "Error getting schedule."
// @Router /schedule [get]
func (h Handler) Schedule(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	schedule, err := h.Service.Schedule(claims)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, schedule)
}
