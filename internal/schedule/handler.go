package schedule

import (
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type Provider interface {
	Schedule(ctx iris.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *iris.Application, verificationMiddleware context.Handler) {
	party := app.Party("/schedule")
	party.Use(verificationMiddleware)
	party.Get("/", h.Schedule)
}

// Schedule
// @Description Get user schedule
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.CareerSchedule	"Successfully retrieved schedule"
// @Failure 400 {string} string "Error getting schedule."
// @Router /schedule [get]
func (h Handler) Schedule(ctx iris.Context) {
	claims := jwt.Get(ctx).(*auth.Claims)

	schedule, err := h.Service.Schedule(claims)
	if err != nil {
		_ = ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		return
	}

	_ = ctx.StopWithJSON(iris.StatusOK, schedule)
}
