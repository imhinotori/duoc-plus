package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Provider interface {
	Authenticate(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *gin.Engine) {
	app.POST("/login", h.Service.AuthMiddleware.LoginHandler)
	party := app.Group("/auth")
	party.GET("/refresh_token", h.Service.AuthMiddleware.RefreshHandler)
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
func (h Handler) Authenticate(ctx *gin.Context) {
	var creds Credentials

	if err := ctx.Bind(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error reading body",
		})
		return

	}

	tokens, err := h.Service.Authenticate(creds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
