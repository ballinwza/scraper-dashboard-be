package usecase_rental_estate

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (u *rentalEstateUsecase) ScraperRentalEstate(ctx context.Context, targetURL string, startPage, maxPage int) error {
	if startPage == 0 {
		startPage = 1
	}

	if maxPage == 0 {
		maxPage = 1
	}

	// Scrapping
	items, err := u.scraperRepo.ScrapeMainPage(targetURL, startPage, maxPage)
	if err != nil {
		logger.Log.Error("Failed to run scraper job in usecase", zap.Error(err))
		return domain.ErrInternalServer
	}

	if len(items) == 0 {
		return domain.ErrInvalidInput
	}

	// Upsert
	now := helper.NowUTC()
	bulkModels := make([]domain.BulkUpsert, 0, len(items))

	for _, item := range items {
		if item.SourceURL == "" {
			continue
		}

		if item.CreatedAt.IsZero() {
			item.CreatedAt = now
		}
		item.UpdatedAt = now

		filter := bson.M{"source_url": item.SourceURL}

		update := bson.M{
			"$set": item,
		}

		bulkModels = append(bulkModels, domain.BulkUpsert{
			Filter: filter,
			Update: update,
		})
	}

	u.mongodbRepo.BulkUpsert(ctx, collectionName, bulkModels)

	logger.Log.Info("UpsertItems complete", zap.Int("total", len(items)))
	return nil
}
