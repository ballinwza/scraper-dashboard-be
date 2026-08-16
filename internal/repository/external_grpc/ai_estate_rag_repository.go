package external_grpc

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"strings"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type IAiEstateRagRepository interface {
	GetAnswer(ctx context.Context, question string) (string, error)
	UploadChunkDocument(ctx context.Context, filename string, fileBytes []byte) (*grpc_api.UploadPdfResponse, error)
	UploadFileStramMultiTenant(ctx context.Context, req domain.UploadFileMultiTenantReq) (*grpc_api.UploadFileStreamResponse, error)
}

type aiEstateRagRepository struct {
	client grpc_api.ChatGRPCClient
}

func NewAiEstateRagRepository(targetAddr string) (IAiEstateRagRepository, func(), error) {
	var transportCreds credentials.TransportCredentials

	// เช็คว่าเป็น Local (localhost / 127.0.0.1) หรือเป็น Production (Cloud Run)
	if helper.IsLocalIP(targetAddr) {
		// Local: ใช้ insecure credentials
		transportCreds = insecure.NewCredentials()
	} else {
		// Production (Cloud Run): บังคับใช้ TLS/SSL (Port 443)
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read system root certs: %w", err)
		}
		transportCreds = credentials.NewClientTLSFromCert(systemRoots, "")
	}

	// สร้าง Connection
	conn, err := grpc.NewClient(
		targetAddr,
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gRPC on %s: %w", targetAddr, err)
	}

	cleanup := func() {
		conn.Close()
	}

	client := grpc_api.NewChatGRPCClient(conn)

	return &aiEstateRagRepository{
		client: client,
	}, cleanup, nil
}

func (c *aiEstateRagRepository) GetAnswer(ctx context.Context, question string) (string, error) {
	res, err := c.client.Query(ctx, &grpc_api.ChatRequest{
		Question: question,
	})
	if err != nil {
		return "", fmt.Errorf("grpc request failed: %w", err)
	}

	return res.GetMessage(), nil
}

// UploadDocument ทำหน้าที่หั่นไฟล์ []byte ออกเป็น Chunks แล้ว Stream ส่งเข้า gRPC Server
func (c *aiEstateRagRepository) UploadChunkDocument(ctx context.Context, filename string, fileBytes []byte) (*grpc_api.UploadPdfResponse, error) {
	// 1. เปิด Client Stream ไปยัง gRPC Server
	stream, err := c.client.UploadPdf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open upload stream: %w", err)
	}

	// 2. ส่ง Metadata เป็น Chunk แรก (ระบุชื่อไฟล์)
	err = stream.Send(&grpc_api.UploadPdfRequest{
		Data: &grpc_api.UploadPdfRequest_Metadata{
			Metadata: &grpc_api.PdfMetadata{
				Filename: filename,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send metadata chunk: %w", err)
	}

	// 3. หั่น []byte ออกเป็น Chunks (ชิ้นละ 64 KB) แล้วทยอยส่งผ่าน Stream
	const chunkSize = 64 * 1024
	buffer := bytes.NewReader(fileBytes)
	chunkBuf := make([]byte, chunkSize)

	for {
		n, err := buffer.Read(chunkBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file buffer: %w", err)
		}

		// ส่งเนื้อไฟล์ส่วนย่อย (chunk_data) ผ่าน stream
		err = stream.Send(&grpc_api.UploadPdfRequest{
			Data: &grpc_api.UploadPdfRequest_ChunkData{
				ChunkData: chunkBuf[:n],
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send chunk data: %w", err)
		}
	}

	// 4. ปิดการส่ง Stream และรอรับ Response สรุปผลจาก Server
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to receive upload response: %w", err)
	}

	return res, nil
}

// UploadDocument ทำหน้าที่หั่นไฟล์ []byte ออกเป็น Chunks แล้ว Stream ส่งเข้า gRPC Server
func (c *aiEstateRagRepository) UploadFileStramMultiTenant(ctx context.Context, req domain.UploadFileMultiTenantReq) (*grpc_api.UploadFileStreamResponse, error) {
	if !isValidFileType(req.FileType) {
		return nil, fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	// 1. เรียก gRPC Method ให้ตรงกับ MultiTenant
	stream, err := c.client.UploadFileStramMultiTenant(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open upload stream: %w", err)
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
		return nil, fmt.Errorf("failed to send metadata chunk: %w", err)
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
			return nil, fmt.Errorf("error reading file buffer: %w", err)
		}

		// ส่ง Chunk Data โดยใช้ UploadFileStreamRequest_ChunkData ให้ถูกต้อง
		err = stream.Send(&grpc_api.UploadFileStreamRequest{
			Payload: &grpc_api.UploadFileStreamRequest_ChunkData{
				ChunkData: chunkBuf[:n],
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send chunk data: %w", err)
		}
	}

	// 4. ปิด Stream และรับ Response (จะได้ *grpc_api.UploadFileStreamResponse)
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("failed to receive upload response: %w", err)
	}

	return res, nil
}

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
