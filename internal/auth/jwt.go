package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imhinotori/duoc-plus/internal/common"
	"time"
)

func (s Service) generateTokenPair(uniqueSessionId string, usr *common.User) (*string, *time.Time, *string, *time.Time, error) {
	expireTime := time.Now().Add(time.Hour * 7)

	claims := &common.JWTClaims{
		ID:       uniqueSessionId,
		Username: usr.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.Config.JWT.Key))
	if err != nil {
		return nil, nil, nil, nil, errors.New("error generating access token")

	}

	refreshTokenExpireTime := time.Now().Add(time.Hour * 24)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["id"] = uniqueSessionId
	refreshTokenClaims["exp"] = refreshTokenExpireTime.Unix()

	rt, err := refreshToken.SignedString([]byte(s.Config.JWT.Key))
	if err != nil {
		return nil, nil, nil, nil, errors.New("error generating refresh token")
	}

	return &t, &expireTime, &rt, &refreshTokenExpireTime, nil
}
