package auth

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Config   *config.Config
	Signer   *jwt.Signer
	Verifier *jwt.Verifier
	Duoc     *duoc.Client
}

func New(cfg *config.Config, signer *jwt.Signer, verifier *jwt.Verifier, duoc *duoc.Client) *Service {
	return &Service{
		Config:   cfg,
		Signer:   signer,
		Verifier: verifier,
		Duoc:     duoc,
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username            string `json:"username"`
	Email               string `json:"email"`     // Its username + @duocuc.cl
	StudentCode         string `json:"codAlumno"` // It's probably an int, but well.
	StudentId           int    `json:"idAlumno"`  // Why two ids (?) I don't know.
	DuocApiBearerToken  string `json:"api_bearer"`
	DuocApiRefreshToken string `json:"refresh_token"`
}

func (s Service) Authenticate(credentials Credentials) (jwt.TokenPair, error) {
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
		return jwt.TokenPair{}, err
	}

	if code != iris.StatusOK {
		return jwt.TokenPair{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var ssoResponseData ssoAuthResponse

	if err2 := json.Unmarshal(response, &ssoResponseData); err2 != nil {
		return jwt.TokenPair{}, err2
	}

	log.Debug("Successfully authenticated User, Getting some general data", "username", credentials.Username)

	// Get some data....

	endpoint = "/vivo_v1.0/v1/matriculaVigente"

	query := url.Values{}
	query.Set("nombreUsuario", credentials.Username)

	response, code, err = s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, ssoResponseData.AccessToken)

	if err != nil {
		return jwt.TokenPair{}, err
	}

	if code != iris.StatusOK {
		return jwt.TokenPair{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData struct {
		CodAlumno string `json:"codAlumno"`
		IDAlumno  int    `json:"idAlumno"`
		Rut       string `json:"rut"`
		RutDV     string `json:"rut_dv"`
	}

	if err2 := json.Unmarshal(response, &responseData); err2 != nil {
		return jwt.TokenPair{}, err2
	}

	log.Debug("Successfully got some general data", "username", credentials.Username)

	claims := Claims{
		Username:           credentials.Username,
		Email:              strings.Replace(credentials.Username, "@duocuc.cl", "", -1),
		StudentCode:        responseData.CodAlumno,
		StudentId:          responseData.IDAlumno,
		DuocApiBearerToken: ssoResponseData.AccessToken,
	}

	refreshClaims := jwt.Claims{Subject: strconv.Itoa(responseData.IDAlumno)}

	log.Debug("User authenticated successfully, returning Tokens", "username", credentials.Username)

	return s.Signer.NewTokenPair(claims, refreshClaims, 120*time.Hour) // TODO!

}
