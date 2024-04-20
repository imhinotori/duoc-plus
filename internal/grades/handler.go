package grades

import (
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Provider interface {
	Attendance(ctx echo.Context) error
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *echo.Echo, authService *auth.Service) {
	party := app.Group("/grades")
	party.Use(authService.AuthMiddleware)
	party.GET("", h.Grades)

}

// Grades
// @Description Get user grades
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.Grades	"Successfully retrieved grades"
// @Failure 400 {string} string "Error getting grades."
// @Router /grades [get]
func (h Handler) Grades(ctx echo.Context) error {
	userSessionId := ctx.Get("user").(*jwt.Token)
	claims := userSessionId.Claims.(*common.JWTClaims)

	usr, err := h.Service.Database.GetUserFromSessionId(claims.ID)
	if err != nil {
		log.Debug("Error getting user from session ID", "error", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	grades, err := h.Service.Grades(*usr)
	if err != nil {
		log.Debug("Error getting grades", "error", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, grades)
}
