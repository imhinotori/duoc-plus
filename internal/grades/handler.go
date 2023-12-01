package grades

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/cache"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"
)

type Provider interface {
	Attendance(ctx iris.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *iris.Application, verificationMiddleware context.Handler) {
	party := app.Party("/grades")
	if h.Service.Config.General.Cache {
		party.Use(cache.Handler(time.Duration(h.Service.Config.General.CacheTime) * time.Second))
	}
	party.Use(verificationMiddleware)
	party.Get("/", h.Grades)
}

// Grades
// @Description Get user grades
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.Grades	"Successfully retrieved grades"
// @Failure 400 {string} string "Error getting grades."
// @Router /grades [get]
func (h Handler) Grades(ctx iris.Context) {
	claims := jwt.Get(ctx).(*auth.Claims)

	grades, err := h.Service.Grades(claims)
	if err != nil {
		_ = ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		log.Debug("Error getting grades", "error", err)
		return
	}

	_ = ctx.StopWithJSON(iris.StatusOK, grades)
}
