package grades

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/auth"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/kataras/iris/v12"
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

func (s Service) Grades(claims *auth.Claims) ([]common.DuocGradesCourses, error) {
	endpoint := "/notas_v1.0/v1/notasAlumno"

	query := url.Values{}
	query.Set("codAlumno", claims.StudentCode)

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims.DuocApiBearerToken)

	if err != nil {
		return []common.DuocGradesCourses{}, err
	}

	if code != iris.StatusOK {
		return []common.DuocGradesCourses{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData []common.DuocGradesCourses

	if err = json.Unmarshal(response, &responseData); err != nil {
		return responseData, err
	}

	log.Debug("Getting grades data", "username", claims.Username)

	return responseData, nil
}
