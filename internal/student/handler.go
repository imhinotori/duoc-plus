package student

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
	StudentData(ctx iris.Context)
}

type Handler struct {
	Service *Service
}

func (h Handler) Start(app *iris.Application, verificationMiddleware context.Handler) {
	party := app.Party("/student")
	if h.Service.Config.General.Cache {
		party.Use(cache.Handler(time.Duration(h.Service.Config.General.CacheTime) * time.Second))
	}
	party.Use(verificationMiddleware)
	party.Get("/", h.StudentData)
}

// StudentData
// @Description Get student information
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.User	"Successfully retrieved student information"
// @Failure 400 {string} string "Error getting student information."
// @Router /student [get]
func (h Handler) StudentData(ctx iris.Context) {
	claims := jwt.Get(ctx).(*auth.Claims)

	studentData, err := h.Service.StudentData(claims)
	if err != nil {
		_ = ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		log.Debug("Error getting student data", "error", err)
		return
	}

	_ = ctx.StopWithJSON(iris.StatusOK, studentData)
}
