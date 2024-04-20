package student

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/database"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Service struct {
	Config   *config.Config
	Duoc     *duoc.Client
	Database *database.Database
	Caser    cases.Caser
}

func New(cfg *config.Config, db *database.Database, duoc *duoc.Client) *Service {
	return &Service{
		Config:   cfg,
		Duoc:     duoc,
		Database: db,
		Caser:    cases.Title(language.LatinAmericanSpanish),
	}
}

func (s Service) StudentData(usr common.User) (common.User, error) {
	endpoint := "/credencial-virtual_v1.0/v1/datosAlumno"

	query := url.Values{}
	query.Set("alumnoId", strconv.Itoa(usr.StudentId))

	log.Debug("Getting student data", "studentId", usr.StudentId)
	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, usr.AccessToken)

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

	log.Debug("Converting student data to new format", "username", usr.Username)

	returnData, err := s.convertDuocStudentDataToStudentData(responseData, usr)
	if err != nil {
		return common.User{}, err
	}

	return returnData, nil
}

func (s Service) convertDuocStudentDataToStudentData(original common.DuocStudentData, usr common.User) (common.User, error) {
	var avatar string

	if original.Avatar != "" {
		avatar = original.Avatar
	} else {
		avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%x", sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(usr.Email)))))
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
