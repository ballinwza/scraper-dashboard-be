package usecase_rental_estate

import (
	"context"
	"errors"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

// DeleteByID ดึงข้อมูลก่อนเพื่อตรวจสอบสิทธิ์/การมีอยู่ แล้วทำการลบรายการออกจาก MongoDB ตาม ID
func (u *rentalEstateUsecase) DeleteRentalEstateById(ctx context.Context, id string) error {
	// 1. Validate Parameter
	if id == "" {
		logger.Log.Warn("DeleteByID: empty ID provided")
		return domain.ErrInvalidInput
	}

	// 2. Check existence in DB before deleting (Optional - แต่แนะนำสำหรับ Clean Architecture)
	existingItem, err := u.mongodbRepo.GetByID(ctx, collectionName, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			logger.Log.Warn("DeleteByID: item not found", zap.String("id", id))
			return domain.ErrNotFound
		}
		logger.Log.Error("DeleteByID: failed to check item existence", zap.String("id", id), zap.Error(err))
		return domain.ErrInternalServer
	}

	if existingItem == nil {
		return domain.ErrNotFound
	}

	// 3. Execute Delete in Repository
	err = u.mongodbRepo.DeleteById(ctx, collectionName, id)
	if err != nil {
		logger.Log.Error("DeleteByID: failed to delete item from repository", zap.String("id", id), zap.Error(err))
		return domain.ErrInternalServer
	}

	logger.Log.Info("DeleteByID: successfully deleted rental estate", zap.String("id", id))
	return nil
}
