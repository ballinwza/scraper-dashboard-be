package usecase_ai_estate_rag

import (
	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
)

// MapKnowledgeFileFromProto แปลงจาก gRPC (Protobuf) เป็น Domain Struct
func mapKnowledgeFileFromProto(res *grpc_api.KnowledgeFile) domain_ai_estate_rag.KnowledgeFile {
	chunks := make([]domain_ai_estate_rag.Chunk, len(res.Chunks))
	for i, pb := range res.Chunks {
		chunks[i] = mapChunkFromProto(pb)
	}

	return domain_ai_estate_rag.KnowledgeFile{
		ID:           res.Id,
		UserID:       res.UserId,
		ChatbotID:    res.ChatbotId,
		Filename:     res.Filename,
		FileType:     res.FileType,
		FileSizeByes: res.FileSizeBytes,
		Status:       domain_ai_estate_rag.FileStatus(res.Status),
		TotalChunks:  res.TotalChunks,
		Chunks:       chunks,
		TotalPage:    res.TotalPage,
		TextContent:  res.TextContent,
		ErrorMessage: res.ErrorMessage,
		CreatedAt:    res.CreatedAt.AsTime(),
		UpdatedAt:    res.UpdatedAt.AsTime(),
	}
}

func mapChunkFromProto(pb *grpc_api.Chunk) domain_ai_estate_rag.Chunk {
	if pb == nil {
		return domain_ai_estate_rag.Chunk{}
	}

	return domain_ai_estate_rag.Chunk{
		VectorID:    pb.VectorId,
		ChunkIndex:  pb.ChunkIndex,
		TextContent: pb.TextContent,
		PageNumber:  pb.PageNumber,
		TokenCount:  pb.TokenCount,
	}
}

func mapListKnowledgeFileFromProto(res *grpc_api.ListKnowledgeFilesResponse) domain_ai_estate_rag.ListKnowledgeFilesResponse {
	files := make([]domain_ai_estate_rag.KnowledgeFile, len(res.Files))
	for i, pb := range res.Files {
		files[i] = mapKnowledgeFileFromProto(pb)
	}

	return domain_ai_estate_rag.ListKnowledgeFilesResponse{
		Files:      files,
		TotalCount: res.TotalCount,
	}
}
