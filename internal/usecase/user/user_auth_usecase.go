package usecase_user

import (
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func (u *userUsecase) generateAccessToken(accessSecret, username, role string, expireTimeMins int) (string, error) {
	claims := domain.JWTAccessClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTimeMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(accessSecret))
}

func (u *userUsecase) generateRefreshToken(refreshSecret, username string, expireTimeMins int) (string, error) {
	claims := domain.RefreshClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTimeMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}
