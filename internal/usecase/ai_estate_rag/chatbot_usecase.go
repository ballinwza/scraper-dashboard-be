package usecase_ai_estate_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
)

func (u *aiEstateRagUsecase) CreateMultiTenantChatbot(
	ctx context.Context,
	userId, name, description, systemPrompt string,
) (*dto.CreateMultiTenantChatbotResDTO, error) {
	res, err := u.grpc.CreateMultiTenantChatbot(ctx, userId, name, description, systemPrompt)
	if err != nil {
		return nil, err
	}

	return &dto.CreateMultiTenantChatbotResDTO{
		Chatbot: dto.ChatbotBlueprintResDTO{
			ID:           res.Chatbot.Id,
			UserID:       res.Chatbot.UserId,
			Name:         res.Chatbot.Name,
			Description:  res.Chatbot.Description,
			SystemPrompt: res.Chatbot.SystemPrompt,
			CreatedAt:    res.Chatbot.CreatedAt.AsTime(),
			UpdatedAt:    res.Chatbot.UpdatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) GetMultiTenantChatbot(
	ctx context.Context,
	id, userId string,
) (*dto.GetMultiTenantChatbotResDTO, error) {
	res, err := u.grpc.GetMultiTenantChatbot(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &dto.GetMultiTenantChatbotResDTO{
		Chatbot: dto.ChatbotBlueprintResDTO{
			ID:           res.Chatbot.Id,
			UserID:       res.Chatbot.UserId,
			Name:         res.Chatbot.Name,
			Description:  res.Chatbot.Description,
			SystemPrompt: res.Chatbot.SystemPrompt,
			CreatedAt:    res.Chatbot.CreatedAt.AsTime(),
			UpdatedAt:    res.Chatbot.UpdatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) ListMultiTenantChatbots(
	ctx context.Context,
	userId string,
	pageSize, pageToken int,
) (*dto.ListMultiTenantChatbotsResDTO, error) {
	res, err := u.grpc.ListMultiTenantChatbots(ctx, userId, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	chatbots := make([]dto.ChatbotBlueprintResDTO, len(res.Chatbots))
	for i, cb := range res.Chatbots {
		chatbots[i] = dto.ChatbotBlueprintResDTO{
			ID:           cb.Id,
			UserID:       cb.UserId,
			Name:         cb.Name,
			Description:  cb.Description,
			SystemPrompt: cb.SystemPrompt,
			CreatedAt:    cb.CreatedAt.AsTime(),
			UpdatedAt:    cb.UpdatedAt.AsTime(),
		}
	}

	return &dto.ListMultiTenantChatbotsResDTO{
		Chatbots:      chatbots,
		NextPageToken: res.NextPageToken,
		TotalCount:    res.TotalCount,
	}, nil
}

func (u *aiEstateRagUsecase) UpdateMultiTenantChatbot(
	ctx context.Context,
	id, userId, name, description, systemPrompt string,
) (*dto.UpdateMultiTenantChatbotResDTO, error) {
	res, err := u.grpc.UpdateMultiTenantChatbot(ctx, id, userId, name, description, systemPrompt)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateMultiTenantChatbotResDTO{
		Chatbot: dto.ChatbotBlueprintResDTO{
			ID:           res.Chatbot.Id,
			UserID:       res.Chatbot.UserId,
			Name:         res.Chatbot.Name,
			Description:  res.Chatbot.Description,
			SystemPrompt: res.Chatbot.SystemPrompt,
			CreatedAt:    res.Chatbot.CreatedAt.AsTime(),
			UpdatedAt:    res.Chatbot.UpdatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) DeleteMultiTenantChatbot(
	ctx context.Context,
	id, userId string,
) (*dto.DeleteMultiTenantChatbotResDTO, error) {
	res, err := u.grpc.DeleteMultiTenantChatbot(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteMultiTenantChatbotResDTO{
		Success: res.Success,
		Message: res.Message,
	}, nil
}
