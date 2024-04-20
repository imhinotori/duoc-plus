package common

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
