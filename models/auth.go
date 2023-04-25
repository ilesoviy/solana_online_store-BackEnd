package models

import "github.com/golang-jwt/jwt"

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom_claims"`
}
