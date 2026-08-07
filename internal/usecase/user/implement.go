package usecase_user

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/mongodb"
)

const collectionName = "user"

type UserUsecase interface {
	Register(ctx context.Context, username, password, name string) (*domain.User, error)
	Login(ctx context.Context, username, password, userAgent, clientIp string) (*string, *string, error)
	Refresh(ctx context.Context, refreshToken string) (*string, *string, error)
}

type userUsecase struct {
	mongodbRepo mongodb.IMongoDBGenericRepository[domain.User]
	cfg         config.Config
}

func NewAuthUsecase(userRepo mongodb.IMongoDBGenericRepository[domain.User], cfg config.Config) UserUsecase {
	return &userUsecase{
		mongodbRepo: userRepo,
		cfg:         cfg,
	}
}
