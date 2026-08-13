package usecase_user

import (
	"context"
	"errors"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (u *userUsecase) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.M{
		"username": username,
	}

	user, err := u.mongodbRepo.FindOneByFilter(ctx, collectionName, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
