package usecase_ai_estate_rag

import (
	"context"

	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/ai_estate_rag_grpc"
)

type IAiEstateRagUsecase interface {
	GetKnowledgeFile(ctx context.Context, id, userId string) (*domain_ai_estate_rag.GetKnowledgeFileResponse, error)
	ListKnowledgeFiles(
		ctx context.Context,
		chatbotId, userId string,
		limit, offset int,
	) (*domain_ai_estate_rag.ListKnowledgeFilesResponse, error)
	DeleteKnowledgeFile(ctx context.Context, chatbotId, userId string) (*domain_ai_estate_rag.DeleteKnowledgeFileResponse, error)
	MultiTenantUploadFile(ctx context.Context, userId, chatbotId, filename, fileType string, fileBytes []byte) (*domain_ai_estate_rag.UploadFileMultiTenantResponse, error)
}

type aiEstateRagUsecase struct {
	grpc ai_estate_rag_grpc.IAiEstateRagRepository
}

func NewAiEstateRagUsecase(grpc ai_estate_rag_grpc.IAiEstateRagRepository) IAiEstateRagUsecase {
	return &aiEstateRagUsecase{
		grpc: grpc,
	}
}
