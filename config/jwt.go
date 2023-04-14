package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("INIPOKONYARAHASIAakuDanKAMU")

type JWTClaim struct {
	UserID string
	jwt.RegisteredClaims
}
