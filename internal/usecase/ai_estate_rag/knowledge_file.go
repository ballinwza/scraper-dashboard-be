package usecase_ai_estate_rag

import (
	"context"

	domain "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
)

func (u *aiEstateRagUsecase) MultiTenantUploadFile(ctx context.Context, userId, chatbotId, filename, fileType string, fileBytes []byte) (*domain_ai_estate_rag.UploadFileMultiTenantResponse, error) {
	mongo_data := domain.UploadFileMultiTenantReq{
		UserId:    userId,
		ChatbotId: chatbotId,
		FileType:  fileType,
		Filename:  filename,
		FileBytes: fileBytes,
	}
	res, _ := u.grpc.UploadFileStramMultiTenant(ctx, mongo_data)

	result := domain_ai_estate_rag.UploadFileMultiTenantResponse{
		FileID:      res.FileId,
		Status:      domain_ai_estate_rag.FileStatus(res.Status),
		TotalChunks: int32(res.TotalChunks),
		TotalBytes:  int64(res.TotalBytes),
		Message:     res.Message,
	}

	return &result, nil
}

func (u *aiEstateRagUsecase) GetKnowledgeFile(ctx context.Context, id, userId string) (*domain_ai_estate_rag.GetKnowledgeFileResponse, error) {
	res, err := u.grpc.GetKnowledgeFile(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	knowledge := mapKnowledgeFileFromProto(res.File)
	result := domain_ai_estate_rag.GetKnowledgeFileResponse{
		File: knowledge,
	}
	return &result, nil
}

func (u *aiEstateRagUsecase) ListKnowledgeFiles(
	ctx context.Context,
	chatbotId, userId string,
	limit, offset int,
) (*domain_ai_estate_rag.ListKnowledgeFilesResponse, error) {
	res, err := u.grpc.ListKnowledgeFiles(ctx, chatbotId, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	listKnowledge := mapListKnowledgeFileFromProto(res)

	return &listKnowledge, nil
}

func (u *aiEstateRagUsecase) DeleteKnowledgeFile(ctx context.Context, chatbotId, userId string) (*domain_ai_estate_rag.DeleteKnowledgeFileResponse, error) {
	res, err := u.grpc.DeleteKnowledgeFile(ctx, chatbotId, userId)
	if err != nil {
		return nil, err
	}

	result := domain_ai_estate_rag.DeleteKnowledgeFileResponse{
		Success: res.Success,
		Message: res.Message,
	}
	return &result, nil
}
