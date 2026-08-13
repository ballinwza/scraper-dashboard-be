package usecase_user

import (
	"context"
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *userUsecase) Refresh(ctx context.Context, username, refreshHashedToken string) (*string, *string, error) {
	now := helper.NowUTC()
	filter := bson.M{
		"username": username,
		"refresh_tokens": bson.M{
			"$elemMatch": bson.M{
				"token_hash": refreshHashedToken,
				"expires_at": bson.M{"$gt": now},
			},
		},
	}

	newRefreshToken, err := u.generateRefreshToken(u.cfg.JwtRefreshSecret, username, u.cfg.JwtRefreshExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	newRefreshTokenHashed := helper.HashTokenSHA256(newRefreshToken)
	newExpireAt := now.Add(time.Duration(u.cfg.JwtRefreshExpirationMins) * time.Minute)
	updateData := bson.M{
		"$set": bson.M{
			"refresh_tokens.$.token_hash": newRefreshTokenHashed,
			"refresh_tokens.$.expires_at": newExpireAt,
			"updated_at":                  now,
		},
	}

	user, err := u.mongodbRepo.FindOneAndUpdate(ctx, collectionName, filter, updateData)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	newAccessToken, err := u.generateAccessToken(u.cfg.JwtAccessSecret, user.Username, "user", u.cfg.JwtAccessExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	return &newAccessToken, &newRefreshToken, nil
}
