package schedule

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
	"strconv"
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

func (s Service) Schedule(claims *auth.Claims) ([]common.DuocSchedule, error) {
	endpoint := "/horario_v1.0/v1/horario"

	query := url.Values{}
	query.Set("alumnoId", strconv.Itoa(claims.StudentId))

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, claims.DuocApiBearerToken)

	if err != nil {
		return []common.DuocSchedule{}, err
	}

	if code != iris.StatusOK {
		return []common.DuocSchedule{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData []common.DuocSchedule

	if err = json.Unmarshal(response, &responseData); err != nil {
		return []common.DuocSchedule{}, err
	}

	log.Debug("Getting schedule data", "username", claims.Username)

	return responseData, nil
}
