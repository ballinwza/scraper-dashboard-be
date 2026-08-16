package usecase_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
)

func (u *ragUsecase) MultiTenantUploadFile(ctx context.Context, userId, chatbotId, filename, fileType string, fileBytes []byte) (*dto.UploadFileMultiTenantResDTO, error) {
	mongo_data := domain.UploadFileMultiTenantReq{
		UserId:    userId,
		ChatbotId: chatbotId,
		FileType:  fileType,
		Filename:  filename,
		FileBytes: fileBytes,
	}
	res, _ := u.grpc.UploadFileStramMultiTenant(ctx, mongo_data)

	result := dto.UploadFileMultiTenantResDTO{
		FileId:      res.FileId,
		Status:      res.Status.String(),
		TotalChunks: int(res.TotalChunks),
		TotalBytes:  int(res.TotalBytes),
		Message:     res.Message,
	}

	return &result, nil
}
