package usecase_user

import (
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func (u *userUsecase) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func (u *userUsecase) checkRefreshToken(user domain.User, incomingTokenHash string) error {
// 	now := time.Now()
// 	var matchedToken *domain.RefreshTokenEmbedded

// 	// 1. ค้นหา Token ชิ้นที่ตรงกับที่ส่งมาจาก Client ก่อน
// 	for _, token := range user.RefreshTokens {
// 		if token.TokenHash == incomingTokenHash {
// 			matchedToken = &token
// 			break
// 		}
// 	}

// 	// 2. ถ้าหา Token ไม่เจอใน Array
// 	if matchedToken == nil {
// 		return errors.New("token not found or revoked")
// 	}

// 	// 3. เช็คว่าหมดอายุหรือยังด้วย .Before()
// 	// ถ้าเวลาหมดอายุ (ExpiresAt) มาถึงก่อน/น้อยกว่า เวลาปัจจุบัน (now) แปลว่าหมดอายุแล้ว
// 	if matchedToken.ExpiresAt.Before(now) {
// 		return errors.New("token expired")
// 	}

// 	return nil
// }
