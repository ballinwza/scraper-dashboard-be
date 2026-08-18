package ai_estate_rag_grpc

import (
	"context"
	"crypto/x509"
	"fmt"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type IAiEstateRagRepository interface {
	GetAnswer(ctx context.Context, question string) (string, error)
	UploadChunkDocument(ctx context.Context, filename string, fileBytes []byte) (*grpc_api.UploadPdfResponse, error)

	// RAG
	SearchSimilar(
		ctx context.Context,
		userId, chatbotId, query string,
		top_k *int,
		knowledge_file_id *string,
	) (*grpc_api.RagResponseDTO, error)

	// Knowledge File
	UploadFileStramMultiTenant(ctx context.Context, req domain_ai_estate_rag.UploadFileMultiTenantReq) (*grpc_api.UploadFileStreamResponse, error)
	GetKnowledgeFile(ctx context.Context, id, userId string) (*grpc_api.GetKnowledgeFileResponse, error)
	ListKnowledgeFiles(ctx context.Context, chatbot_id, userId string, limit, offset int) (*grpc_api.ListKnowledgeFilesResponse, error)
	DeleteKnowledgeFile(ctx context.Context, chatbot_id, userId string) (*grpc_api.DeleteKnowledgeFileResponse, error)

	// Chatbot
	CreateMultiTenantChatbot(
		ctx context.Context,
		userId, name, description, systemPrompt string,
	) (*grpc_api.CreateMultiTenantChatbotResponse, error)
	GetMultiTenantChatbot(ctx context.Context, id, userId string) (*grpc_api.GetMultiTenantChatbotResponse, error)
	ListMultiTenantChatbots(
		ctx context.Context,
		userId string,
		pageSize, pageToken int,
	) (*grpc_api.ListMultiTenantChatbotsResponse, error)
	UpdateMultiTenantChatbot(
		ctx context.Context,
		id, userId, name, description, systemPrompt string,
	) (*grpc_api.UpdateMultiTenantChatbotResponse, error)
	DeleteMultiTenantChatbot(ctx context.Context, id, userId string) (*grpc_api.DeleteMultiTenantChatbotResponse, error)

	// Chat Session
	CreateChatSession(
		ctx context.Context,
		userId, chatbotId, sessionTitle string,
	) (*grpc_api.CreateChatSessionResponse, error)
	GetChatSession(ctx context.Context, id, userId string) (*grpc_api.GetChatSessionResponse, error)
	ListChatSessions(
		ctx context.Context,
		userId, chatbotId string,
		pageSize, pageToken int,
	) (*grpc_api.ListChatSessionsResponse, error)
	AddChatMessage(
		ctx context.Context,
		sessionId, userId, content string,
		role domain_ai_estate_rag.MessageRole,
	) (*grpc_api.AddChatMessageResponse, error)
	DeleteChatSession(ctx context.Context, id, userId string) error
}

type aiEstateRagRepository struct {
	chat           grpc_api.ChatGRPCClient
	knowledge_file grpc_api.KnowledgeFileServiceClient
	chat_session   grpc_api.ChatSessionServiceClient
	rag            grpc_api.RagServiceClient
	chatbot        grpc_api.ChatbotServiceClient
}

func NewAiEstateRagRepository(targetAddr string) (IAiEstateRagRepository, func(), error) {
	var transportCreds credentials.TransportCredentials

	// เช็คว่าเป็น Local (localhost / 127.0.0.1) หรือเป็น Production (Cloud Run)
	if helper.IsLocalIP(targetAddr) {
		// Local: ใช้ insecure credentials
		transportCreds = insecure.NewCredentials()
	} else {
		// Production (Cloud Run): บังคับใช้ TLS/SSL (Port 443)
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read system root certs: %w", err)
		}
		transportCreds = credentials.NewClientTLSFromCert(systemRoots, "")
	}

	// สร้าง Connection
	conn, err := grpc.NewClient(
		targetAddr,
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gRPC on %s: %w", targetAddr, err)
	}

	cleanup := func() {
		conn.Close()
	}

	client := grpc_api.NewChatGRPCClient(conn)
	knowledge := grpc_api.NewKnowledgeFileServiceClient(conn)
	chat_session := grpc_api.NewChatSessionServiceClient(conn)
	rag := grpc_api.NewRagServiceClient(conn)
	chatbot := grpc_api.NewChatbotServiceClient(conn)

	return &aiEstateRagRepository{
		chat:           client,
		knowledge_file: knowledge,
		chat_session:   chat_session,
		rag:            rag,
		chatbot:        chatbot,
	}, cleanup, nil
}
