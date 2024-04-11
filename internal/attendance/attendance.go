package attendance

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
	"regexp"
	"strings"
	"sync"
	"time"
)

type Service struct {
	Config *config.Config
	Duoc   *duoc.Client
}

var monthsMap = map[string]string{
	"ENERO":      "January",
	"FEBRERO":    "February",
	"MARZO":      "March",
	"ABRIL":      "April",
	"MAYO":       "May",
	"JUNIO":      "June",
	"JULIO":      "July",
	"AGOSTO":     "August",
	"SEPTIEMBRE": "September",
	"OCTUBRE":    "October",
	"NOVIEMBRE":  "November",
	"DICIEMBRE":  "December",
}

func New(cfg *config.Config, duoc *duoc.Client) *Service {
	return &Service{
		Config: cfg,
		Duoc:   duoc,
	}
}

func (s Service) Attendance(claims jwt.MapClaims) ([]common.Attendance, error) {
	endpoint := "/asistencia_v1.0/v1/asistenciaCompleta"

	query := url.Values{}
	query.Set("codAlumno", claims["student_code"].(string))

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims["api_bearer"].(string))

	if err != nil {
		return []common.Attendance{}, err
	}

	if code == http.StatusNoContent {
		return []common.Attendance{}, common.NoContentError
	} else if code != http.StatusOK {
		return []common.Attendance{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData []common.DuocAttendance

	if err = json.Unmarshal(response, &responseData); err != nil {
		return []common.Attendance{}, err
	}

	log.Debug("Getting attendance data", "username", claims["username"].(string))

	returnData := make([]common.Attendance, len(responseData))

	for i, duocAttendance := range responseData {
		returnData[i] = convertDuocAttendanceToAttendance(duocAttendance)
	}

	return returnData, nil
}

func convertDuocAttendanceToAttendance(original common.DuocAttendance) common.Attendance {
	attendance := common.Attendance{
		DegreeName: original.DegreeName,
		DegreeCode: original.DegreeCode,
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, duocSubjectAttendance := range original.SubjectAttendance {
		wg.Add(1)
		go func(subjectAttendance common.DuocSubjectAttendance) {
			defer wg.Done()

			newSubjectAttendance := common.SubjectAttendance{
				Name:            subjectAttendance.Header.SubjectName,
				Code:            subjectAttendance.Header.SubjectCode,
				ClassesHeld:     helper.ConvertToInt(subjectAttendance.Header.ClassesHeld),
				AssistedClasses: helper.ConvertToInt(subjectAttendance.Header.AssistedClasses),
				Percentage:      helper.ConvertToFloat64(subjectAttendance.Header.Percentage),
			}

			for _, detail := range subjectAttendance.Details {
				date, err := convertDuocDateToDate(detail.Date)
				if err != nil {
					log.Debug("Error converting date", "error", err)
					continue
				}

				newSubjectAttendance.Details = append(newSubjectAttendance.Details, common.AttendanceDetail{
					Date:       date,
					Attendance: helper.ConvertToInt(detail.Attendance),
				})
			}

			mu.Lock()
			attendance.SubjectAttendance = append(attendance.SubjectAttendance, newSubjectAttendance)
			mu.Unlock()
		}(duocSubjectAttendance)
	}

	wg.Wait()
	return attendance
}

func convertDuocDateToDate(originalDate string) (time.Time, error) {
	// Eliminar el día de la semana y las comas
	re := regexp.MustCompile(`^[A-ZÁÉÍÓÚÑ]+\s*,\s*`)
	result := re.ReplaceAllString(originalDate, "")

	// Reemplazar los "de" por un espacio
	result = strings.Replace(result, "de", " ", -1)

	result = strings.TrimSpace(result)

	for spanishMonth, englishMonth := range monthsMap {
		result = strings.Replace(result, spanishMonth, englishMonth, -1)
	}

	return time.Parse("02 January 2006", result)
}
