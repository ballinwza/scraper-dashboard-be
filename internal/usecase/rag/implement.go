package usecase_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/repository/external_grpc"
)

type IRagUsecase interface {
	GenerateAnswer(ctx context.Context, question string) (string, error)
}

type ragUsecase struct {
	grpc external_grpc.IAiEstateRagRepository
}

func NewRagUsecase(grpc external_grpc.IAiEstateRagRepository) IRagUsecase {
	return &ragUsecase{
		grpc: grpc,
	}
}
