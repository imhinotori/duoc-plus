package schedule

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (s Service) Schedule(claims jwt.MapClaims) ([]common.CareerSchedule, error) {
	endpoint := "/horario_v1.0/v1/horario"

	query := url.Values{}
	query.Set("alumnoId", strconv.Itoa(int(claims["student_id"].(float64))))

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims["api_bearer"].(string))

	if err != nil {
		return []common.CareerSchedule{}, err
	}

	if code != http.StatusOK {
		return []common.CareerSchedule{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData []common.DuocSchedule

	if err = json.Unmarshal(response, &responseData); err != nil {
		return []common.CareerSchedule{}, err
	}

	log.Debug("Getting schedule data", "username", claims["username"].(string))

	caser := cases.Title(language.LatinAmericanSpanish)
	schedule := convertDuocScheduleToSchedule(responseData, caser)

	return schedule, nil
}

func convertDuocScheduleToSchedule(s []common.DuocSchedule, caser cases.Caser) []common.CareerSchedule {
	schedule := make([]common.CareerSchedule, len(s))

	for i, career := range s {
		schedule[i].CareerName = strings.Replace(caser.String(career.NameCareer), "  ", " ", -1)
		schedule[i].Schedule = common.Week{}

		for _, day := range career.Days {
			course := convertDuocCourseToCourse(day.Courses, caser)

			switch day.Day {
			case "1":
				schedule[i].Schedule.Monday = course
			case "2":
				schedule[i].Schedule.Tuesday = course
			case "3":
				schedule[i].Schedule.Wednesday = course
			case "4":
				schedule[i].Schedule.Thursday = course
			case "5":
				schedule[i].Schedule.Friday = course
			case "6":
				schedule[i].Schedule.Saturday = course
			case "7":
				schedule[i].Schedule.Sunday = course
			}

		}
	}

	return schedule

}

func convertDuocCourseToCourse(c []common.DuocCourse, caser cases.Caser) []common.Course {
	course := make([]common.Course, 0, len(c))

	for _, item := range c {
		course = append(course, common.Course{
			SubjectName: strings.Replace(caser.String(item.Name), "  ", " ", -1),
			SubjectCode: item.Code,
			Instructor:  strings.Replace(caser.String(item.Instructor), "  ", " ", -1),
			StartTime:   item.StartTime,
			EndTime:     item.EndTime,
			Classroom:   item.Classroom,
		})
	}

	return course
}
