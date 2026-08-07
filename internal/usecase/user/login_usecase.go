package usecase_user

import (
	"context"
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *userUsecase) Login(ctx context.Context, username, password, userAgent, clientIp string) (*string, *string, error) {
	now := helper.NowUTC()

	// Verify PasswordHash
	filter := bson.M{
		"username": username,
	}

	// user, err := u.mongodbRepo.FindOneByFilter(ctx, collectionName, filter)
	// if err != nil {
	// 	return nil, nil, domain.ErrInvalidCredentials
	// }

	// if !u.checkPasswordHash(password, user.PasswordHash) {
	// 	return nil, nil, domain.ErrInvalidCredentials
	// }

	// Generate new AccessToken & Refresh Token
	accessToken, err := u.generateAccessToken(u.cfg.JwtAccessSecret, username, "user", u.cfg.JwtAccessExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	refreshToken, err := u.generateRefreshToken(u.cfg.JwtRefreshSecret, username, u.cfg.JwtRefreshExpirationMins)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}

	// updateDate := domain.User{
	// 	RefreshTokens: domain.RefreshTokenEmbedded{
	// 		TokenHash: helper.HashTokenSHA256(refreshToken),
	// 		UserAgent: userAgent,
	// 		ClientIP:  clientIp,
	// 		ExpiresAt: now.Add(time.Duration(u.cfg.JwtRefreshExpirationMins) * time.Minute),
	// 	},
	// 	UpdatedAt: now,
	// }

	// err = u.mongodbRepo.Update(ctx, collectionName, user.ID, updateDate)
	// if err != nil {
	// 	return nil, nil, domain.ErrInternalServer
	// }

	// filter := bson.M{"_id": userID}

	// ใช้ Aggregation Pipeline ใน Update Command เพื่อกรอง IP เก่าออก แล้ว concat array ใหม่เข้าไป
	newExpireAt := now.Add(time.Duration(u.cfg.JwtAccessExpirationMins) * time.Minute)
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$set", Value: bson.M{
			"updated_at": now,
			"refresh_tokens": bson.M{
				"$concatArrays": bson.A{
					// กรองเอาเฉพาะ token ที่ ip_address ไม่ตรงกับ ip ปัจจุบันไว้
					bson.M{
						"$filter": bson.M{
							"input": bson.M{"$ifNull": bson.A{"$refresh_tokens", bson.A{}}},
							"as":    "item",
							"cond":  bson.M{"$ne": bson.A{"$$item.client_ip", clientIp}},
						},
					},
					bson.A{
						bson.M{
							"token_hash": helper.HashTokenSHA256(refreshToken),
							"client_ip":  clientIp,
							"user_agent": userAgent,
							"expires_at": newExpireAt,
						},
					},
				},
			},
		}}},
	}

	// opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	_, err = u.mongodbRepo.FindOneAndUpdate(ctx, collectionName, filter, pipeline)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	return &accessToken, &refreshToken, nil

}
