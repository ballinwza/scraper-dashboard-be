package usecase_rental_estate

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/mongodb"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/scraper"
)

const collectionName = "rental_estate"

type IRentalEstateUsecase interface {
	ScraperRentalEstate(ctx context.Context, targetURL string, startPage, maxPage int) error
	GetRentalEstateItem(ctx context.Context, id string) (*domain.RentalEstate, error)
	DeleteRentalEstateById(ctx context.Context, id string) error
	FetchRentalEstates(ctx context.Context, filter domain.RentalEstateFilter) (*domain.PaginatedResponse[domain.RentalEstate], error)
	ExportCSVRentalEstate(ctx context.Context, filter domain.RentalEstateFilter) ([]domain.RentalEstate, error)
}

type rentalEstateUsecase struct {
	scraperRepo scraper.IDotpropertyScraperRepository
	mongodbRepo mongodb.IMongoDBGenericRepository[domain.RentalEstate]
}

func NewScraperRentalEstateUsecase(
	scraperRepo scraper.IDotpropertyScraperRepository,
	mongodbRepo mongodb.IMongoDBGenericRepository[domain.RentalEstate],
) IRentalEstateUsecase {
	return &rentalEstateUsecase{
		scraperRepo: scraperRepo,
		mongodbRepo: mongodbRepo,
	}
}
