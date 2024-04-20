package auth

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"net/http"
	"net/url"
)

type enrollmentResponseData struct {
	StudentCode string `json:"codAlumno"`
	StudentId   int    `json:"idAlumno"`
	Rut         string `json:"rut"`
	RutDV       string `json:"rut_dv"`
}

func (s Service) authenticationRequest(credentials Credentials) (*ssoAuthResponse, error) {
	log.Debug("Trying to authenticate user", "username", credentials.Username)
	endpoint := "/auth/realms/WEB_APPS_PRD/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", s.Config.Duoc.ClientId)
	data.Set("client_secret", s.Config.Duoc.ClientSecret)
	data.Set("grant_type", "password")
	data.Set("username", credentials.Username)
	data.Set("password", credentials.Password)

	response, code, err := s.Duoc.Request(s.Config.Duoc.SSOURL+endpoint, "POST", []byte(data.Encode()), nil)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var ssoResponseData ssoAuthResponse

	if err = json.Unmarshal(response, &ssoResponseData); err != nil {
		return nil, err
	}

	return &ssoResponseData, nil
}

func (s Service) accountDetailsRequest(username, token string) (*enrollmentResponseData, error) {
	endpoint := "/vivo_v1.0/v1/matriculaVigente"

	query := url.Values{}
	query.Set("nombreUsuario", username)

	response, code, err := s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, token)

	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData enrollmentResponseData

	if err = json.Unmarshal(response, &responseData); err != nil {
		return nil, err
	}

	return &responseData, nil
}
