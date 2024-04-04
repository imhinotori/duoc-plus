package grades

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
	party := app.Group("/grades")
	party.Use(authService.AuthMiddleware.MiddlewareFunc())
	{
		if storage != nil && storage[0] != nil {
			party.GET("/", cache.CachePage(storage[0], time.Minute, func(c *gin.Context) {
				h.Grades(c)
			}))
		} else {
			party.GET("/", h.Grades)
		}
	}

}

// Grades
// @Description Get user grades
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.Grades	"Successfully retrieved grades"
// @Failure 400 {string} string "Error getting grades."
// @Router /grades [get]
func (h Handler) Grades(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)

	grades, err := h.Service.Grades(claims)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		log.Debug("Error getting grades", "error", err)
		return
	}

	ctx.JSON(http.StatusOK, grades)
}
