package ai_estate_rag_grpc

import (
	"context"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

func (c *aiEstateRagRepository) SearchSimilar(
	ctx context.Context,
	userId, chatbotId, query string,
	top_k *int,
	knowledge_file_id *string,
) (*grpc_api.RagResponseDTO, error) {
	req := &grpc_api.RagSearchSimilarRequestDTO{
		UserId:    userId,
		ChatbotId: chatbotId,
		QueryText: query,
	}

	if top_k != nil {
		k := int32(*top_k)
		req.TopK = &k
	} else {
		DEFAULT_TOP_K := 10
		top_k = &DEFAULT_TOP_K
	}

	if knowledge_file_id != nil {
		req.KnowledgeFileId = knowledge_file_id
	}

	res, err := c.rag.SearchSimilar(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}
