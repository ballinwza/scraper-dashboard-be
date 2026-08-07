package usecase_user

import (
	"context"
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *userUsecase) Refresh(ctx context.Context, refreshToken string) (*string, *string, error) {
	claims := &domain.RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(u.cfg.JwtRefreshSecret), nil // ใช้ refreshSecret
	})
	if err != nil || !token.Valid {
		return nil, nil, domain.ErrInvalidToken
	}

	// Check User
	tokenHash := helper.HashTokenSHA256(refreshToken)
	now := helper.NowUTC()
	filter := bson.M{
		"username": claims.Username,
		"refresh_tokens": bson.M{
			"$elemMatch": bson.M{
				"token_hash": tokenHash,
				"expires_at": bson.M{"$gt": now}, // $gt = Greater Than ( expires_at ต้องมากกว่าเวลาปัจจุบัน)
			},
		},
	}

	newRefreshToken, err := u.generateRefreshToken(u.cfg.JwtRefreshSecret, claims.Username, u.cfg.JwtRefreshExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	newRefreshTokenHashed := helper.HashTokenSHA256(newRefreshToken)
	newExpireAt := now.Add(time.Duration(u.cfg.JwtRefreshExpirationMins) * time.Minute)
	updateData := bson.M{
		"$set": bson.M{
			"refresh_tokens.$.token_hash": newRefreshTokenHashed,
			"refresh_tokens.$.expires_at": newExpireAt,
		},
	}

	user, err := u.mongodbRepo.FindOneAndUpdate(ctx, collectionName, filter, updateData)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Generate new AccessToken & Refresh Token
	newAccessToken, err := u.generateAccessToken(u.cfg.JwtAccessSecret, user.Username, "user", u.cfg.JwtAccessExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	return &newAccessToken, &newRefreshToken, nil
}
