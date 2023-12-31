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
		Username:            credentials.Username,
		Email:               strings.Replace(credentials.Username, "@duocuc.cl", "", -1),
		StudentCode:         responseData.CodAlumno,
		StudentId:           responseData.IDAlumno,
		DuocApiBearerToken:  ssoResponseData.AccessToken,
		DuocApiRefreshToken: ssoResponseData.RefreshToken,
	}

	refreshClaims := jwt.Claims{Subject: strconv.Itoa(claims.StudentId)}

	log.Debug("User authenticated successfully, returning Tokens", "username", credentials.Username)

	return s.GenerateTokenPair(claims, refreshClaims)
}

func (s Service) RefreshToken(claims *Claims) (jwt.TokenPair, error) {
	log.Debug("Trying to refresh user token", "username", claims.Username)
	endpoint := "/auth/realms/WEB_APPS_PRD/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", s.Config.Duoc.ClientId)
	data.Set("client_secret", s.Config.Duoc.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", claims.DuocApiRefreshToken)

	response, code, err := s.Duoc.Request(s.Config.Duoc.SSOURL+endpoint, "POST", []byte(data.Encode()), nil)

	if err != nil {
		return jwt.TokenPair{}, err
	}

	if code != iris.StatusOK {
		log.Debug("Error refreshing user tokens", "username", claims.Username, "error", string(response), "data", data)
		return jwt.TokenPair{}, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var ssoResponseData ssoAuthResponse

	if err = json.Unmarshal(response, &ssoResponseData); err != nil {
		return jwt.TokenPair{}, err
	}

	log.Debug("Successfully refreshed user tokens", "username", claims.Username)

	claims.DuocApiBearerToken = ssoResponseData.AccessToken
	claims.DuocApiRefreshToken = ssoResponseData.RefreshToken

	refreshClaims := jwt.Claims{Subject: strconv.Itoa(claims.StudentId)} // Assuming StudentId is the appropriate field

	log.Debug("User tokens refreshed successfully", "username", claims.Username)

	return s.GenerateTokenPair(*claims, refreshClaims)
}

func (s Service) GenerateTokenPair(claims Claims, refreshClaims jwt.Claims) (jwt.TokenPair, error) {
	log.Debug("Generating user token", "username", claims.Username)
	return s.Signer.NewTokenPair(claims, refreshClaims, 120*time.Hour) // TODO!
}
