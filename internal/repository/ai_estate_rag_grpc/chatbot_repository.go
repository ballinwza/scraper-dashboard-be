package ai_estate_rag_grpc

import (
	"context"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

func (c *aiEstateRagRepository) CreateMultiTenantChatbot(
	ctx context.Context,
	userId, name, description, systemPrompt string,
) (*grpc_api.CreateMultiTenantChatbotResponse, error) {
	req := &grpc_api.CreateMultiTenantChatbotRequest{
		UserId:       userId,
		Name:         name,
		Description:  description,
		SystemPrompt: systemPrompt,
	}

	res, err := c.chatbot.CreateMultiTenantChatbot(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) GetMultiTenantChatbot(
	ctx context.Context,
	id, userId string,
) (*grpc_api.GetMultiTenantChatbotResponse, error) {
	req := &grpc_api.GetMultiTenantChatbotRequest{
		Id:     id,
		UserId: userId,
	}

	res, err := c.chatbot.GetMultiTenantChatbot(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) ListMultiTenantChatbots(
	ctx context.Context,
	userId string,
	pageSize, pageToken int,
) (*grpc_api.ListMultiTenantChatbotsResponse, error) {
	req := &grpc_api.ListMultiTenantChatbotsRequest{
		UserId:    userId,
		PageSize:  int32(pageSize),
		PageToken: int32(pageToken),
	}

	res, err := c.chatbot.ListMultiTenantChatbots(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) UpdateMultiTenantChatbot(
	ctx context.Context,
	id, userId, name, description, systemPrompt string,
) (*grpc_api.UpdateMultiTenantChatbotResponse, error) {
	req := &grpc_api.UpdateMultiTenantChatbotRequest{
		Id:           id,
		UserId:       userId,
		Name:         name,
		Description:  description,
		SystemPrompt: systemPrompt,
	}

	res, err := c.chatbot.UpdateMultiTenantChatbot(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) DeleteMultiTenantChatbot(
	ctx context.Context,
	id, userId string,
) (*grpc_api.DeleteMultiTenantChatbotResponse, error) {
	req := &grpc_api.DeleteMultiTenantChatbotRequest{
		Id:     id,
		UserId: userId,
	}

	res, err := c.chatbot.DeleteMultiTenantChatbot(ctx, req)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}
