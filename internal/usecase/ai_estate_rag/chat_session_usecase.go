package usecase_ai_estate_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
)

func (u *aiEstateRagUsecase) CreateChatSession(
	ctx context.Context,
	userId, chatbotId, sessionTitle string,
) (*dto.CreateChatSessionResDTO, error) {
	res, err := u.grpc.CreateChatSession(ctx, userId, chatbotId, sessionTitle)
	if err != nil {
		return nil, err
	}

	messages := make([]dto.ChatMessageResDTO, len(res.Session.Messages))
	for i, msg := range res.Session.Messages {
		messages[i] = dto.ChatMessageResDTO{
			Role:      domain_ai_estate_rag.MessageRole(msg.Role),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.AsTime(),
		}
	}

	return &dto.CreateChatSessionResDTO{
		Session: dto.ChatSessionResDTO{
			ID:           res.Session.Id,
			UserID:       res.Session.UserId,
			ChatbotID:    res.Session.ChatbotId,
			SessionTitle: res.Session.SessionTitle,
			Messages:     messages,
			CreatedAt:    res.Session.CreatedAt.AsTime(),
			UpdatedAt:    res.Session.UpdatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) GetChatSession(
	ctx context.Context,
	id, userId string,
) (*dto.GetChatSessionResDTO, error) {
	res, err := u.grpc.GetChatSession(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	messages := make([]dto.ChatMessageResDTO, len(res.Session.Messages))
	for i, msg := range res.Session.Messages {
		messages[i] = dto.ChatMessageResDTO{
			Role:      domain_ai_estate_rag.MessageRole(msg.Role),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.AsTime(),
		}
	}

	return &dto.GetChatSessionResDTO{
		Session: dto.ChatSessionResDTO{
			ID:           res.Session.Id,
			UserID:       res.Session.UserId,
			ChatbotID:    res.Session.ChatbotId,
			SessionTitle: res.Session.SessionTitle,
			Messages:     messages,
			CreatedAt:    res.Session.CreatedAt.AsTime(),
			UpdatedAt:    res.Session.UpdatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) ListChatSessions(
	ctx context.Context,
	userId, chatbotId string,
	pageSize, pageToken int,
) (*dto.ListChatSessionsResDTO, error) {
	res, err := u.grpc.ListChatSessions(ctx, userId, chatbotId, pageSize, pageToken)
	if err != nil {
		return nil, err
	}

	sessions := make([]dto.ChatSessionResDTO, len(res.Sessions))
	for i, sess := range res.Sessions {
		messages := make([]dto.ChatMessageResDTO, len(sess.Messages))
		for j, msg := range sess.Messages {
			messages[j] = dto.ChatMessageResDTO{
				Role:      domain_ai_estate_rag.MessageRole(msg.Role),
				Content:   msg.Content,
				CreatedAt: msg.CreatedAt.AsTime(),
			}
		}

		sessions[i] = dto.ChatSessionResDTO{
			ID:           sess.Id,
			UserID:       sess.UserId,
			ChatbotID:    sess.ChatbotId,
			SessionTitle: sess.SessionTitle,
			Messages:     messages,
			CreatedAt:    sess.CreatedAt.AsTime(),
			UpdatedAt:    sess.UpdatedAt.AsTime(),
		}
	}

	return &dto.ListChatSessionsResDTO{
		Sessions:   sessions,
		TotalCount: res.TotalCount,
	}, nil
}

func (u *aiEstateRagUsecase) AddChatMessage(
	ctx context.Context,
	sessionId, userId, content string,
	role domain_ai_estate_rag.MessageRole,
) (*dto.AddChatMessageResDTO, error) {
	res, err := u.grpc.AddChatMessage(ctx, sessionId, userId, content, role)
	if err != nil {
		return nil, err
	}

	return &dto.AddChatMessageResDTO{
		Message: dto.ChatMessageResDTO{
			Role:      domain_ai_estate_rag.MessageRole(res.Message.Role),
			Content:   res.Message.Content,
			CreatedAt: res.Message.CreatedAt.AsTime(),
		},
	}, nil
}

func (u *aiEstateRagUsecase) DeleteChatSession(
	ctx context.Context,
	id, userId string,
) error {
	return u.grpc.DeleteChatSession(ctx, id, userId)
}
