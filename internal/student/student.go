package student

import (
	"crypto/sha256"
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
	Caser  cases.Caser
}

func New(cfg *config.Config, duoc *duoc.Client) *Service {
	return &Service{
		Config: cfg,
		Duoc:   duoc,
		Caser:  cases.Title(language.LatinAmericanSpanish),
	}
}

func (s Service) StudentData(claims jwt.MapClaims) (common.User, error) {
	endpoint := "/credencial-virtual_v1.0/v1/datosAlumno"

	query := url.Values{}
	query.Set("alumnoId", strconv.Itoa(int(claims["student_id"].(float64))))

	log.Debug("Getting student data", "studentId", claims["student_id"])
	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims["api_bearer"].(string))

	if err != nil {
		return common.User{}, err
	}

	if code != http.StatusOK {
		return common.User{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData common.DuocStudentData

	if err = json.Unmarshal(response, &responseData); err != nil {
		return common.User{}, err
	}

	log.Debug("Converting student data to new format", "username", claims["username"].(string))

	returnData, err := s.convertDuocStudentDataToStudentData(responseData, claims)
	if err != nil {
		return common.User{}, err
	}

	return returnData, nil
}

func (s Service) convertDuocStudentDataToStudentData(original common.DuocStudentData, claims jwt.MapClaims) (common.User, error) {
	var avatar string

	if original.Avatar != "" {
		avatar = original.Avatar
	} else {
		avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%x", sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(claims["email"].(string))))))
	}

	NewStudentData := common.User{
		FullName: s.Caser.String(original.NombreCompleto),
		Rut:      original.Rut,
		Avatar:   avatar,
	}

	careers := make([]common.Career, 0, len(original.Carreras))

	for _, duocCareer := range original.Carreras {
		careers = append(careers, common.Career{
			School:     duocCareer.Escuela,
			CareerName: duocCareer.NomCarrera,
			CareerCode: duocCareer.CodCarrera,
			Campus:     duocCareer.Sede,
		})
	}

	NewStudentData.Careers = careers

	return NewStudentData, nil
}

type Career struct {
	School     string `json:"school"`
	CareerName string `json:"career_name"`
	CareerCode string `json:"career_code"`
	Campus     string `json:"campus"`
}
