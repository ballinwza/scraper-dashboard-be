package usecase_ai_estate_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/ai_estate_rag_grpc"
)

type IAiEstateRagUsecase interface {
	// RAG
	SearchSimilar(
		ctx context.Context,
		userId, chatbotId, queryText string,
		top_k *int,
		knowledge_file_id *string,
	) (*dto.SearchSimilarResDTO, error)

	// Knowledge File
	GetKnowledgeFile(ctx context.Context, id, userId string) (*domain_ai_estate_rag.GetKnowledgeFileResponse, error)
	ListKnowledgeFiles(
		ctx context.Context,
		chatbotId, userId string,
		limit, offset int,
	) (*domain_ai_estate_rag.ListKnowledgeFilesResponse, error)
	DeleteKnowledgeFile(ctx context.Context, chatbotId, userId string) (*domain_ai_estate_rag.DeleteKnowledgeFileResponse, error)
	MultiTenantUploadFile(ctx context.Context, userId, chatbotId, filename, fileType string, fileBytes []byte) (*domain_ai_estate_rag.UploadFileMultiTenantResponse, error)

	// Chatbot
	CreateMultiTenantChatbot(
		ctx context.Context,
		userId, name, description, systemPrompt string,
	) (*dto.CreateMultiTenantChatbotResDTO, error)
	GetMultiTenantChatbot(ctx context.Context, id, userId string) (*dto.GetMultiTenantChatbotResDTO, error)
	ListMultiTenantChatbots(
		ctx context.Context,
		userId string,
		pageSize, pageToken int,
	) (*dto.ListMultiTenantChatbotsResDTO, error)
	UpdateMultiTenantChatbot(
		ctx context.Context,
		id, userId, name, description, systemPrompt string,
	) (*dto.UpdateMultiTenantChatbotResDTO, error)
	DeleteMultiTenantChatbot(ctx context.Context, id, userId string) (*dto.DeleteMultiTenantChatbotResDTO, error)

	// Chat Session
	CreateChatSession(
		ctx context.Context,
		userId, chatbotId, sessionTitle string,
	) (*dto.CreateChatSessionResDTO, error)
	GetChatSession(ctx context.Context, id, userId string) (*dto.GetChatSessionResDTO, error)
	ListChatSessions(
		ctx context.Context,
		userId, chatbotId string,
		pageSize, pageToken int,
	) (*dto.ListChatSessionsResDTO, error)
	AddChatMessage(
		ctx context.Context,
		sessionId, userId, content string,
		role domain_ai_estate_rag.MessageRole,
	) (*dto.AddChatMessageResDTO, error)
	DeleteChatSession(ctx context.Context, id, userId string) error
}

type aiEstateRagUsecase struct {
	grpc ai_estate_rag_grpc.IAiEstateRagRepository
}

func NewAiEstateRagUsecase(grpc ai_estate_rag_grpc.IAiEstateRagRepository) IAiEstateRagUsecase {
	return &aiEstateRagUsecase{
		grpc: grpc,
	}
}
