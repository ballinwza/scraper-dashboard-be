package usecase_rag

import (
	"context"

	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/ai_estate_rag_grpc"
)

type IRagUsecase interface {
	GenerateAnswer(ctx context.Context, question string) (string, error)
	ReceiverOfUploadFile(ctx context.Context, file []byte, filename string) (*domain_ai_estate_rag.UploadPdfResponse, error)
}

type ragUsecase struct {
	grpc ai_estate_rag_grpc.IAiEstateRagRepository
}

func NewRagUsecase(grpc ai_estate_rag_grpc.IAiEstateRagRepository) IRagUsecase {
	return &ragUsecase{
		grpc: grpc,
	}
}
