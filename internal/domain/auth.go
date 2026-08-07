package domain

import "github.com/golang-jwt/jwt/v5"

type JWTAccessClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
