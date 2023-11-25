package grades

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
	party := app.Party("/grades")
	party.Use(verificationMiddleware)
	party.Get("/", h.Grades)
}

// Grades
// @Description Get user grades
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {object} common.GradesCourses	"Successfully retrieved grades"
// @Failure 400 {object} string "Error getting grades."
// @Router /grades [get]
func (h Handler) Grades(ctx iris.Context) {
	claims := jwt.Get(ctx).(*auth.Claims)

	attendance, err := h.Service.Grades(claims)
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"message": err.Error(),
		})
		log.Debug("Error getting attendance", "error", err)
		return
	}

	ctx.StopWithJSON(iris.StatusOK, attendance)
}
