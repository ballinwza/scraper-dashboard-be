package ai_estate_rag_grpc

import (
	"context"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

func (c *aiEstateRagRepository) CreateChatSession(
	ctx context.Context,
	userId, chatbotId, sessionTitle string,
) (*grpc_api.CreateChatSessionResponse, error) {
	req := &grpc_api.CreateChatSessionRequest{
		UserId:       userId,
		ChatbotId:    chatbotId,
		SessionTitle: sessionTitle,
	}

	res, err := c.chat_session.CreateChatSession(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) GetChatSession(
	ctx context.Context,
	id, userId string,
) (*grpc_api.GetChatSessionResponse, error) {
	req := &grpc_api.GetChatSessionRequest{
		Id:     id,
		UserId: userId,
	}

	res, err := c.chat_session.GetChatSession(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) ListChatSessions(
	ctx context.Context,
	userId, chatbotId string,
	pageSize, pageToken int,
) (*grpc_api.ListChatSessionsResponse, error) {
	req := &grpc_api.ListChatSessionsRequest{
		UserId:    userId,
		ChatbotId: chatbotId,
		PageSize:  int32(pageSize),
		PageToken: int32(pageToken),
	}

	res, err := c.chat_session.ListChatSessions(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) AddChatMessage(
	ctx context.Context,
	sessionId, userId, content string,
	role domain_ai_estate_rag.MessageRole,
) (*grpc_api.AddChatMessageResponse, error) {
	req := &grpc_api.AddChatMessageRequest{
		SessionId: sessionId,
		UserId:    userId,
		Content:   content,
		Role:      grpc_api.MessageRole(role),
	}

	res, err := c.chat_session.AddChatMessage(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) DeleteChatSession(
	ctx context.Context,
	id, userId string,
) error {
	req := &grpc_api.DeleteChatSessionRequest{
		Id:     id,
		UserId: userId,
	}

	_, err := c.chat_session.DeleteChatSession(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return err
	}

	return nil
}
