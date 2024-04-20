package schedule

import (
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cache/persistence"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Provider interface {
	Schedule(ctx echo.Context) error
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *echo.Echo, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/schedule")
	party.Use(authService.AuthMiddleware)
	party.GET("", h.Schedule)
}

// Schedule
// @Description Get user schedule
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.CareerSchedule	"Successfully retrieved schedule"
// @Failure 400 {string} string "Error getting schedule."
// @Router /schedule [get]
func (h Handler) Schedule(ctx echo.Context) error {
	userSessionId := ctx.Get("user").(*jwt.Token)
	claims := userSessionId.Claims.(*common.JWTClaims)

	usr, err := h.Service.Database.GetUserFromSessionId(claims.ID)
	if err != nil {
		log.Debug("Error getting user from session ID", "error", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	schedule, err := h.Service.Schedule(*usr)
	if err != nil {
		log.Debugf("Error getting schedule: %s", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, schedule)
}
