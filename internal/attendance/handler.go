package attendance

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type Provider interface {
	Attendance(ctx iris.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *iris.Application, verificationMiddleware context.Handler) {
	party := app.Party("/attendance")
	party.Use(verificationMiddleware)
	party.Get("/", h.Attendance)
}

func (h Handler) Attendance(ctx iris.Context) {
	claims := jwt.Get(ctx).(*auth.Claims)

	attendance, err := h.Service.Attendance(claims)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		log.Debug("Error getting attendance", "error", err)
		return
	}

	ctx.StopWithJSON(iris.StatusOK, attendance)
}
