package auth

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	Config         *config.Config
	Duoc           *duoc.Client
	AuthMiddleware *jwt.GinJWTMiddleware
}

func New(cfg *config.Config, duoc *duoc.Client) *Service {
	service := &Service{
		Config: cfg,
		Duoc:   duoc,
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Duoc UC",
		Key:         []byte(cfg.JWT.Key),
		Timeout:     time.Hour * 5,
		MaxRefresh:  time.Hour * 10,
		IdentityKey: jwt.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*common.User); ok {
				return jwt.MapClaims{
					jwt.IdentityKey: v.Username,
					"username":      v.Username,
					"email":         v.Email,
					"student_code":  v.StudentCode,
					"student_id":    v.StudentId,
					"api_bearer":    v.AccessToken,
					"refresh_token": v.RefreshToken,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Debug("Extracting claims", "claims", claims)
			return &common.User{
				Username: claims[jwt.IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Debug("Authenticating user")
			var creds Credentials
			if err := c.ShouldBind(&creds); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			if usr, err := service.Authenticate(creds); err == nil && usr != nil {
				return usr, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*common.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Errorf("Error creating auth middleware: %s", err)
		return nil
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Errorf("Error initializing auth middleware: %s", errInit)
		return nil
	}

	service.AuthMiddleware = authMiddleware

	return service
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Service) Authenticate(credentials Credentials) (*common.User, error) {
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

	if err2 := json.Unmarshal(response, &ssoResponseData); err2 != nil {
		return nil, err2
	}

	log.Debug("Successfully authenticated User, Getting some general data", "username", credentials.Username)

	// Get some data....

	endpoint = "/vivo_v1.0/v1/matriculaVigente"

	query := url.Values{}
	query.Set("nombreUsuario", credentials.Username)

	response, code, err = s.Duoc.RequestWithQuery(s.Config.Duoc.MobileAPIUrl+endpoint, "GET", nil, query, ssoResponseData.AccessToken)

	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("invalid response structure: %s", string(response))
	}

	var responseData struct {
		CodAlumno string `json:"codAlumno"`
		IDAlumno  int    `json:"idAlumno"`
		Rut       string `json:"rut"`
		RutDV     string `json:"rut_dv"`
	}

	if err2 := json.Unmarshal(response, &responseData); err2 != nil {
		return nil, err2
	}

	log.Debug("Successfully got some general data", "username", credentials.Username)

	usr := &common.User{
		Email:        strings.Replace(credentials.Username, "@duocuc.cl", "", -1),
		Rut:          responseData.Rut + "-" + responseData.RutDV,
		Username:     credentials.Username,
		StudentCode:  responseData.CodAlumno,
		StudentId:    responseData.IDAlumno,
		AccessToken:  ssoResponseData.AccessToken,
		RefreshToken: ssoResponseData.RefreshToken,
	}

	return usr, nil
}
