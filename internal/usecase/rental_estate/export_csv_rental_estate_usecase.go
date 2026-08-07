package usecase_rental_estate

import (
	"context"
	"regexp"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// GetExportData ดึงข้อมูลรายการอสังหาฯ ทั้งหมดตามเงื่อนไข Filter เพื่อนำไป Export ออกเป็นไฟล์ (ไม่จำกัด Pagination Limit)
func (u *rentalEstateUsecase) ExportCSVRentalEstate(ctx context.Context, filter domain.RentalEstateFilter) ([]domain.RentalEstate, error) {
	// 1. Build BSON Filter (ใช้ Filter เดียวกับฟังก์ชันค้นหาทั่วไป)
	mongoFilter := bson.M{}

	// Search (Text Search บน Title หรือ Location)
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

	// 2. Build Mongo FindOptions (ไม่ใส่ SetLimit/SetSkip เพื่อดึงข้อมูลทั้งหมดสำหรับ Export)
	findOptions := options.Find()

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

	// 3. Call Generic Repository
	// ไม่ตั้งค่า Limit และ Offset ให้กับ FindPaginated หรือส่ง FindOptions ที่ไม่มี Limit
	items, _, err := u.mongodbRepo.FindPaginated(ctx, collectionName, mongoFilter, findOptions)
	if err != nil {
		logger.Log.Error("GetExportData: failed to fetch export data from repository",
			zap.Any("filter", filter),
			zap.Error(err),
		)
		return nil, domain.ErrInternalServer
	}

	// ป้องกันการส่งคืน nil ให้กับ Handler
	if items == nil {
		items = []domain.RentalEstate{}
	}

	logger.Log.Info("GetExportData: successfully fetched export data",
		zap.Int("total_records", len(items)),
	)

	return items, nil
}
