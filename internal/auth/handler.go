package auth

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type Provider interface {
	Authenticate(ctx iris.Context)
	RefreshToken(ctx iris.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *iris.Application, verificationMiddleware context.Handler) {
	party := app.Party("/auth")
	party.Post("/", h.Authenticate)
	party.Post("/refresh", h.RefreshToken).Use(verificationMiddleware)
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
func (h Handler) Authenticate(ctx iris.Context) {
	var creds Credentials

	if err := ctx.ReadBody(&creds); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": "Error reading body",
		})
		return
	}

	tokens, err := h.Service.Authenticate(creds)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		return
	}

	ctx.StopWithJSON(iris.StatusOK, tokens)
}

// RefreshToken
// @Description Refresh Token
// @Accept  json
// @Produce  json
// @Param   refresh_token     query    string     true        "Refresh token"
// @Success 200 {object} common.AuthenticationResponse	"ok"
// @Failure 400 {object} string "Error reading body."
// @Router /auth/refresh [post]
func (h Handler) RefreshToken(ctx iris.Context) {
	claims := jwt.Get(ctx).(*Claims)

	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ReadBody(&data); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": "Error reading body",
		})
		return
	}

	tokens, err := h.Service.RefreshToken(claims)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		return
	}

	ctx.StopWithJSON(iris.StatusOK, tokens)
}
