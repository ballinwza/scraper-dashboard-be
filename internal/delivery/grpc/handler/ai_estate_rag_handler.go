package grpc_handler

import (
	"context"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	usecase_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AiEstateRagHandler struct {
	grpc_api.UnimplementedChatGRPCServer
	ragUsecase usecase_rag.IRagUsecase
}

var _ grpc_api.ChatGRPCServer = (*AiEstateRagHandler)(nil)

func NewAiEstateRagHandler(aiUsecase usecase_rag.IRagUsecase) *AiEstateRagHandler {
	return &AiEstateRagHandler{
		ragUsecase: aiUsecase,
	}
}

func (h *AiEstateRagHandler) Query(ctx context.Context, req *grpc_api.ChatRequest) (*grpc_api.ChatResponse, error) {
	question := req.GetQuestion()

	message, err := h.ragUsecase.GenerateAnswer(ctx, question)
	if err != nil {
		// จัดการ Error ส่งกลับไปในรูปแบบ gRPC Status Code
		return nil, status.Errorf(codes.Internal, "failed to process greeting: %v", err)
	}

	return &grpc_api.ChatResponse{
		Message: message,
	}, nil
}
