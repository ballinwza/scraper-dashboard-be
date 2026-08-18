package ai_estate_rag_grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	domain_ai_estate_rag "github.com/ballinwza/scraper-dashboard-be/internal/domain/ai_estate_rag"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.uber.org/zap"
)

// UploadDocument ทำหน้าที่หั่นไฟล์ []byte ออกเป็น Chunks แล้ว Stream ส่งเข้า gRPC Server
func (c *aiEstateRagRepository) UploadFileStramMultiTenant(ctx context.Context, req domain_ai_estate_rag.UploadFileMultiTenantReq) (*grpc_api.UploadFileStreamResponse, error) {
	if !isValidFileType(req.FileType) {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(fmt.Errorf("unsupported file type: %s", req.FileType)),
		)
		return nil, fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	// 1. เรียก gRPC Method ให้ตรงกับ MultiTenant
	stream, err := c.knowledge_file.CreateKnowledgeFile(ctx)
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	// 2. ส่ง Metadata Chunk แรกให้ถูกต้องตาม Struct ของ UploadFileStreamRequest
	err = stream.Send(&grpc_api.UploadFileStreamRequest{
		Payload: &grpc_api.UploadFileStreamRequest_Metadata{
			Metadata: &grpc_api.FileMetadata{
				UserId:    req.UserId,
				ChatbotId: req.ChatbotId,
				Filename:  req.Filename,
				FileType:  req.FileType,
			},
		},
	})
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	// 3. หั่น req.FileBytes ออกเป็น Chunks (ชิ้นละ 64 KB) แล้วทยอยส่งผ่าน Stream
	const chunkSize = 64 * 1024
	buffer := bytes.NewReader(req.FileBytes) // ใช้ req.FileBytes (หรือ field ที่เก็บ bytes ใน DTO)
	chunkBuf := make([]byte, chunkSize)

	for {
		n, err := buffer.Read(chunkBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(
				domain.ErrInternalServer.Error(),
				zap.Error(err),
			)
			return nil, err
		}

		// ส่ง Chunk Data โดยใช้ UploadFileStreamRequest_ChunkData ให้ถูกต้อง
		err = stream.Send(&grpc_api.UploadFileStreamRequest{
			Payload: &grpc_api.UploadFileStreamRequest_ChunkData{
				ChunkData: chunkBuf[:n],
			},
		})
		if err != nil {
			logger.Error(
				domain.ErrInternalServer.Error(),
				zap.Error(err),
			)
			return nil, err
		}
	}

	// 4. ปิด Stream และรับ Response (จะได้ *grpc_api.UploadFileStreamResponse)
	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) GetKnowledgeFile(ctx context.Context, id, userId string) (*grpc_api.GetKnowledgeFileResponse, error) {
	res, err := c.knowledge_file.GetKnowledgeFile(ctx, &grpc_api.GetKnowledgeFileRequest{
		Id:     id,
		UserId: userId,
	})
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) ListKnowledgeFiles(ctx context.Context, chatbot_id, userId string, limit, offset int) (*grpc_api.ListKnowledgeFilesResponse, error) {
	res, err := c.knowledge_file.ListKnowledgeFiles(ctx, &grpc_api.ListKnowledgeFilesRequest{
		ChatbotId: chatbot_id,
		UserId:    userId,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *aiEstateRagRepository) DeleteKnowledgeFile(ctx context.Context, chatbot_id, userId string) (*grpc_api.DeleteKnowledgeFileResponse, error) {
	res, err := c.knowledge_file.DeleteKnowledgeFile(ctx, &grpc_api.DeleteKnowledgeFileRequest{
		ChatbotId: chatbot_id,
		UserId:    userId,
	})
	if err != nil {
		logger.Error(
			domain.ErrInternalServer.Error(),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

// Helper
func isValidFileType(fileType string) bool {
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/webp":      true,
		"image/gif":       true,
	}
	return allowedTypes[strings.ToLower(fileType)]
}
