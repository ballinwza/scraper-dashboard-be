package usecase_rental_estate

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

func (u *rentalEstateUsecase) GetRentalEstateItem(ctx context.Context, id string) (*domain.RentalEstate, error) {
	item, err := u.mongodbRepo.GetByID(ctx, collectionName, id)
	if err != nil {
		logger.Log.Error("Failed to GetRentalEstateItem by ID", zap.Error(err))
		return nil, err
	}

	return item, nil
}
