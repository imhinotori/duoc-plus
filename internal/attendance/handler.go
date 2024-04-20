package attendance

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Provider interface {
	Attendance(ctx *gin.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *echo.Echo, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/attendance")
	log.Infof("Registering attendance routes: %+v", authService.AuthMiddleware)
	party.Use(authService.AuthMiddleware)
	party.GET("", h.Attendance)
}

// Attendance
// @Description Get user attendance
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.Attendance	"Successfully retrieved attendance"
// @Success 204 {object} common.Attendance	"No content"
// @Failure 400 {string} string "Error getting attendance."
// @Router /attendance [get]
func (h Handler) Attendance(ctx echo.Context) error {
	userSessionId := ctx.Get("user").(*jwt.Token)
	claims := userSessionId.Claims.(*common.JWTClaims)

	usr, err := h.Service.Database.GetUserFromSessionId(claims.ID)
	if err != nil {
		log.Debug("Error getting user from session ID", "error", err)
		return ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	attendance, err := h.Service.Attendance(*usr)
	if err != nil {
		if errors.Is(err, common.NoContentError) {
			return ctx.JSON(http.StatusNoContent, map[string]interface{}{
				"message": "No content",
			})
		}
		log.Debug("Error getting attendance", "error", err)

		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Error getting attendance",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, attendance)
}
