package grpc

import (
	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	grpc_handler "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/handler"
	usecase_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rag"
	"google.golang.org/grpc"
)

func NewGRPCServer(aiUsecase usecase_rag.IRagUsecase) *grpc.Server {
	grpcServer := grpc.NewServer()

	// Register Handlers

	aiHandler := grpc_handler.NewAiEstateRagHandler(aiUsecase)
	grpc_api.RegisterChatGRPCServer(grpcServer, aiHandler)
	return grpcServer
}
