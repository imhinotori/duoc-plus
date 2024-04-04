package grades

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/imhinotori/duoc-plus/internal/helper"
	"net/http"
	"net/url"
)

type Service struct {
	Config *config.Config
	Duoc   *duoc.Client
}

func New(cfg *config.Config, duoc *duoc.Client) *Service {
	return &Service{
		Config: cfg,
		Duoc:   duoc,
	}
}

func (s Service) Grades(claims jwt.MapClaims) ([]common.Grades, error) {
	endpoint := "/notas_v1.0/v1/notasAlumno"

	query := url.Values{}
	query.Set("codAlumno", claims["student_code"].(string))

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims["api_bearer"].(string))

	if err != nil {
		return []common.Grades{}, err
	}

	if code != http.StatusOK {
		return []common.Grades{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData []common.DuocGradesCourses

	if err = json.Unmarshal(response, &responseData); err != nil {
		return []common.Grades{}, err
	}

	log.Debug("Getting grades data", "username", claims["username"])

	grades := make([]common.Grades, 0, len(responseData))

	for i, data := range responseData {
		grades = append(grades, convertDuocGradesToGrades(data))

		log.Debug("Getting grades data", "username", claims["username"], "course", i, "courseName", data.DegreeName)
	}

	return grades, nil
}

func convertDuocGradesToGrades(original common.DuocGradesCourses) common.Grades {
	grades := common.Grades{
		Name: original.DegreeName,
		Code: original.DegreeCode,
	}

	for _, subject := range original.Subjects {
		newSubject := common.Subject{
			Code:     subject.SubjectCode,
			Name:     subject.Name,
			Average:  helper.ConvertToFloat64(subject.Average),
			Partials: make([]float64, 0, len(subject.PartialGrades)),
			Exams:    make([]float64, 0, len(subject.ExamsGrades)),
		}

		for _, grade := range subject.PartialGrades {
			if grade.Text == "Nota Final" || grade.Text == "Nota Presentación" {
				continue
			}

			newSubject.Partials = append(newSubject.Partials, helper.ConvertToFloat64(grade.Grade))
		}

		for _, grade := range subject.ExamsGrades {
			if grade.Text == "Nota Final" || grade.Text == "Nota Presentación" {
				continue
			}
			newSubject.Exams = append(newSubject.Exams, helper.ConvertToFloat64(grade.Grade))
		}

		grades.Subjects = append(grades.Subjects, newSubject)
	}

	return grades
}
