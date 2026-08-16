package usecase_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
)

func (u *ragUsecase) ReceiverOfUploadFile(ctx context.Context, fileBytes []byte, filename string) (*domain.UploadPdfResponse, error) {
	res, _ := u.grpc.UploadChunkDocument(ctx, filename, fileBytes)

	result := domain.UploadPdfResponse{
		FileId:  res.FileId,
		Success: res.Success,
		Message: res.Message,
	}

	return &result, nil
}
