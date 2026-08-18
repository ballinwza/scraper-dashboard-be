package usecase_rag

import (
	"context"

	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
)

func (u *ragUsecase) ReceiverOfUploadFile(ctx context.Context, fileBytes []byte, filename string) (*domain_ai_estate_rag.UploadPdfResponse, error) {
	res, _ := u.grpc.UploadChunkDocument(ctx, filename, fileBytes)

	result := domain_ai_estate_rag.UploadPdfResponse{
		FileId:  res.FileId,
		Success: res.Success,
		Message: res.Message,
	}

	return &result, nil
}
