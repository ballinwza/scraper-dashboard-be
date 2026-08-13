package usecase_user

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *userUsecase) Logout(ctx context.Context, username, refreshHashToken string) error {
	now := helper.NowUTC()
	// 1. สร้าง Filter สำหรับค้นหาเอกสารที่มี token_hash ใน refresh_tokens
	filter := bson.M{
		"username":                  username,
		"refresh_tokens.token_hash": refreshHashToken,
	}

	// 2. สร้าง Update Operator ใช้ $pull เพื่อลบ token_hash ออกจาก Array และอัปเดต updated_at
	update := bson.M{
		"$pull": bson.M{
			"refresh_tokens": bson.M{
				"$or": []bson.M{
					{"token_hash": refreshHashToken},
					{"expires_at": bson.M{"$lte": now}},
				},
			},
		},
		"$set": bson.M{
			"updated_at": now,
		},
	}

	// 3. ส่ง filter และ update ไปให้ Generic Repository ดำเนินการ
	return u.mongodbRepo.UpdateOne(ctx, collectionName, filter, update)
}
