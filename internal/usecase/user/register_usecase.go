package usecase_user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUsecase) Register(
	ctx context.Context,
	username, password, name string,
) (*domain.User, error) {
	filter := domain.User{
		Username: username,
	}

	fmt.Println("filter : ", filter)

	existingUser, _ := u.mongodbRepo.FindOneByFilter(ctx, collectionName, filter)
	if existingUser != nil {
		return nil, errors.New("Username already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	// indexes := []mongo.IndexModel{
	// 	// 1. Unique Index สำหรับ Username (ป้องกันชื่อซ้ำในระดับ DB)
	// 	{
	// 		Keys:    bson.D{{Key: "username", Value: 1}},
	// 		Options: options.Index().SetUnique(true).SetName("idx_username_unique"),
	// 	},
	// 	// 2. Unique Index สำหรับ Email
	// 	{
	// 		Keys:    bson.D{{Key: "email", Value: 1}},
	// 		Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
	// 	},
	// 	// 3. Index สำหรับค้นหา Refresh Token Hash ภายใน Array
	// 	{
	// 		Keys:    bson.D{{Key: "refresh_tokens.token_hash", Value: 1}},
	// 		Options: options.Index().SetName("idx_refresh_token_hash"),
	// 	},
	// }

	passwordHashed := string(hashedPassword)
	role := "user"
	now := helper.NowUTC()
	user := &domain.User{
		Username:     username,
		PasswordHash: passwordHashed,
		Name:         name,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}

	if err := u.mongodbRepo.Create(ctx, collectionName, user); err != nil {
		return nil, err
	}

	return user, nil
}
