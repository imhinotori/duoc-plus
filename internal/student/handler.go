package student

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
	StudentData(ctx echo.Context) error
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *echo.Echo, authService *auth.Service, storage ...persistence.CacheStore) {
	party := app.Group("/student")
	party.Use(authService.AuthMiddleware)
	party.GET("", h.StudentData)

}

// StudentData
// @Description Get student information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.User	"Successfully retrieved student information"
// @Failure 400 {string} string "Error getting student information."
// @Router /student [get]
func (h Handler) StudentData(ctx echo.Context) error {
	userSessionId := ctx.Get("user").(*jwt.Token)
	claims := userSessionId.Claims.(*common.JWTClaims)

	usr, err := h.Service.Database.GetUserFromSessionId(claims.ID)
	if err != nil {
		log.Debug("Error getting user from session ID", "error", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	studentData, err := h.Service.StudentData(*usr)
	if err != nil {
		log.Debug("Error getting student data", "error", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, studentData)
}
