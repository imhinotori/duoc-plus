package auth

import (
	"github.com/golang-jwt/jwt/v5"
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
	party.POST("/refresh", h.Refresh)
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
	expirationTime := expireTime.Sub(time.Now())

	err = h.Service.saveAccountDetails(user, uniqueSessionId, expirationTime)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Error saving account details",
		})
	}

	t, tokenType, rt, refreshTokenExpireTime, err := h.Service.generateTokenPair(uniqueSessionId, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Error generating token pair",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{"access_token": t, "expires_on": tokenType.Format(time.RFC3339), "refresh_token": rt, "refresh_expires_on": refreshTokenExpireTime.Format(time.RFC3339)})
}

// Refresh
// @Description Refresh
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.AuthenticationResponse	"ok"
// @Failure 400 {object} string "Error reading body."
// @Router /auth/refresh [get]
func (h Handler) Refresh(ctx echo.Context) error {
	tokenRequestBody := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	err := ctx.Bind(&tokenRequestBody)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{"message": "Error reading body"})
	}

	t, err := jwt.Parse(tokenRequestBody.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(h.Service.Config.JWT.Key), nil
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid token"})
	}

	if c, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		id := c["id"].(string)
		if id == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid token"})
		}

		user, err := h.Service.getAccountDetails(id)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid token"})
		}

		token, tokenRefreshTime, refreshToken, refreshTokenExpireTime, err := h.Service.generateTokenPair(id, &user)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]any{
				"message": "Error generating token pair",
			})
		}

		err = h.Service.saveAccountDetails(&user, id, tokenRefreshTime.Sub(time.Now()))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]any{
				"message": "Error saving account details",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]any{"access_token": token, "expires_on": tokenRefreshTime.Format(time.RFC3339), "refresh_token": refreshToken, "refresh_expires_on": refreshTokenExpireTime.Format(time.RFC3339)})
	}

	return ctx.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid token"})
}
