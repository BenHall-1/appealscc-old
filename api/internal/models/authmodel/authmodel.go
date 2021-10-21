package authmodel

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Id          string `json:"Id"`
	Email       string `json:"Email"`
	GlobalAdmin bool   `json:"GlobalAdmin"`
	PremiumType int    `json:"PremiumType"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}
