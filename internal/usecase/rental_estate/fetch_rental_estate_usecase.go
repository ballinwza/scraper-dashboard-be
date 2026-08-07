package usecase_rental_estate

import (
	"context"
	"math"
	"regexp"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// FetchRentalEstates ดึงรายการข้อมูลอสังหาฯ ตาม Filter และจัดการโครงสร้าง Paginated Response
func (u *rentalEstateUsecase) FetchRentalEstates(ctx context.Context, filter domain.RentalEstateFilter) (*domain.PaginatedResponse[domain.RentalEstate], error) {
	// 1. Sanitize input
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	// 2. Build BSON Filter
	mongoFilter := bson.M{}

	// Search (Text search บน Title หรือ Location)
	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"title": bson.M{"$regex": regexp.QuoteMeta(filter.Search), "$options": "i"}},
			{"location": bson.M{"$regex": regexp.QuoteMeta(filter.Search), "$options": "i"}},
		}
	}

	// Price Range Filter
	priceFilter := bson.M{}
	if filter.MinPrice > 0 {
		priceFilter["$gte"] = filter.MinPrice
	}
	if filter.MaxPrice > 0 {
		priceFilter["$lte"] = filter.MaxPrice
	}
	if len(priceFilter) > 0 {
		mongoFilter["price"] = priceFilter
	}

	// 3. Build Mongo FindOptions (Pagination & Sorting)
	findOptions := options.Find()
	findOptions.SetLimit(int64(filter.Limit))
	findOptions.SetSkip(int64((filter.Page - 1) * filter.Limit))

	// Sorting
	sortOrder := -1
	if filter.Order == "asc" {
		sortOrder = 1
	}
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	// 4. Call Generic Repository
	// u.rentalEstateRepo คือ IMongoDBGenericRepository[domain.RentalEstate]
	items, totalCount, err := u.mongodbRepo.FindPaginated(ctx, collectionName, mongoFilter, findOptions)
	if err != nil {
		logger.Error("FetchRentalEstates: failed to query repository", zap.Error(err))
		return nil, domain.ErrDatabase
	}

	// 5. Build Response
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	return &domain.PaginatedResponse[domain.RentalEstate]{
		Data: items,
		Pagination: domain.PaginationMeta{
			CurrentPage: filter.Page,
			Limit:       filter.Limit,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
			HasNext:     filter.Page < totalPages,
			HasPrev:     filter.Page > 1,
		},
	}, nil
}
