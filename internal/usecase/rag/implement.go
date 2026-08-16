package usecase_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/external_grpc"
)

type IRagUsecase interface {
	GenerateAnswer(ctx context.Context, question string) (string, error)
	ReceiverOfUploadFile(ctx context.Context, file []byte, filename string) (*domain.UploadPdfResponse, error)
	MultiTenantUploadFile(ctx context.Context, userId, chatbotId, filename, fileType string, fileBytes []byte) (*dto.UploadFileMultiTenantResDTO, error)
}

type ragUsecase struct {
	grpc external_grpc.IAiEstateRagRepository
}

func NewRagUsecase(grpc external_grpc.IAiEstateRagRepository) IRagUsecase {
	return &ragUsecase{
		grpc: grpc,
	}
}
