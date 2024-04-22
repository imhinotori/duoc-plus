package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Provider interface {
	Authenticate(ctx echo.Context) error
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *echo.Echo) {
	party := app.Group("/auth")
	party.POST("", h.Authenticate)
}

// Authenticate
// @Description Authenticate
// @Accept  json
// @Produce  json
// @Param   username     query    string     true        "Username of the user"
// @Param   password      query    string     true        "Password of the user"
// @Success 200 {object} common.AuthenticationResponse	"ok"
// @Failure 400 {object} string "Error reading body."
// @Router /auth [post]
func (h Handler) Authenticate(ctx echo.Context) error {
	var credentials Credentials

	if err := ctx.Bind(&credentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "Error reading body",
		})

	}

	user, err := h.Service.Authenticate(credentials)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	uniqueSessionId := h.Service.IDGenerator()

	expireTime := time.Now().Add(time.Hour * 7)

	claims := &common.JWTClaims{
		ID:       uniqueSessionId,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	err = h.Service.saveAccountDetails(user, uniqueSessionId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Error saving account details",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.Service.Config.JWT.Key))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Error generating token",
		})

	}

	return ctx.JSON(http.StatusOK, map[string]any{"access_token": t, "expires_on": expireTime.Format(time.RFC3339)})
}
