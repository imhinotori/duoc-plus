package auth

import (
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/common"
	"github.com/imhinotori/duoc-plus/internal/config"
	"github.com/imhinotori/duoc-plus/internal/database"
	"github.com/imhinotori/duoc-plus/internal/duoc"
	"github.com/jaevor/go-nanoid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type Service struct {
	Config         *config.Config
	Duoc           *duoc.Client
	Database       *database.Database
	AuthMiddleware echo.MiddlewareFunc
	IDGenerator    func() string
}

func New(cfg *config.Config, db *database.Database, duoc *duoc.Client) *Service {
	idGenerator, err := nanoid.Standard(21)
	if err != nil {
		log.Error("error setting up ID Generator", "error", err)
		return nil
	}

	service := &Service{
		Config:      cfg,
		Database:    db,
		IDGenerator: idGenerator,
		Duoc:        duoc,
	}

	jwtMiddlewareConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(common.JWTClaims)
		},
		SigningKey: []byte(cfg.JWT.Key),
	}

	service.AuthMiddleware = echojwt.WithConfig(jwtMiddlewareConfig)

	return service
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Service) Authenticate(credentials Credentials) (*common.User, error) {
	ssoData, err := s.authenticationRequest(credentials)
	if err != nil {
		return nil, err

	}

	enrollmentData, err := s.accountDetailsRequest(credentials.Username, ssoData.AccessToken)
	if err != nil {
		return nil, err
	}

	usr := &common.User{
		Email:                 credentials.Username,
		Rut:                   enrollmentData.Rut + "-" + enrollmentData.RutDV,
		Username:              strings.Replace(credentials.Username, "@duocuc.cl", "", -1),
		StudentCode:           enrollmentData.StudentCode,
		StudentId:             enrollmentData.StudentId,
		AccessToken:           ssoData.AccessToken,
		AccessTokenExpiresIn:  ssoData.ExpiresIn,
		RefreshToken:          ssoData.RefreshToken,
		RefreshTokenExpiresIn: ssoData.RefreshExpiresIn,
	}

	return usr, nil
}

func (s Service) saveAccountDetails(account *common.User, id string, expirationTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	jsonUsr, err := json.Marshal(account)
	if err != nil {
		return err
	}

	s.Database.Users.Set(ctx, id, jsonUsr, expirationTime)

	return nil

}

func (s Service) getAccountDetails(uniqueId string) (common.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	data, err := s.Database.Users.Get(ctx, uniqueId).Result()
	if err != nil {
		return common.User{}, err
	}

	var usr common.User
	err = json.Unmarshal([]byte(data), &usr)
	if err != nil {
		return common.User{}, err
	}

	return usr, nil
}
