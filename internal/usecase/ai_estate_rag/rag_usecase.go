package usecase_ai_estate_rag

import (
	"context"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/dto"
)

func (u *aiEstateRagUsecase) SearchSimilar(
	ctx context.Context,
	userId, chatbotId, queryText string,
	top_k *int,
	knowledge_file_id *string,
) (*dto.SearchSimilarResDTO, error) {
	res, err := u.grpc.SearchSimilar(ctx, userId, chatbotId, queryText, top_k, knowledge_file_id)
	if err != nil {
		return nil, err
	}

	sources := make([]dto.SearchVectorRecordItemDTO, len(res.Sources))
	for i, item := range res.Sources {
		var values []float32
		if item.Record != nil && item.Record.Values != nil {
			values = item.Record.Values
		}

		var metadata dto.MetadataVectorRecordDTO
		if item.Record != nil && item.Record.Metadata != nil {
			metadata = dto.MetadataVectorRecordDTO{
				UserID:      item.Record.Metadata.UserId,
				ChatbotID:   item.Record.Metadata.ChatbotId,
				FileID:      item.Record.Metadata.FileId,
				ChunkIndex:  item.Record.Metadata.ChunkIndex,
				TextContent: item.Record.Metadata.TextContent,
				PageNumber:  item.Record.Metadata.PageNumber,
				Filename:    item.Record.Metadata.Filename,
			}
		}

		var recordID string
		if item.Record != nil {
			recordID = item.Record.Id
		}

		sources[i] = dto.SearchVectorRecordItemDTO{
			Score: item.Score,
			Record: dto.VectorRecordDTO{
				ID:       recordID,
				Values:   values,
				Metadata: metadata,
			},
		}
	}

	return &dto.SearchSimilarResDTO{
		AnswerMessage: res.AnswerMessage,
		Sources:       sources,
	}, nil
}
