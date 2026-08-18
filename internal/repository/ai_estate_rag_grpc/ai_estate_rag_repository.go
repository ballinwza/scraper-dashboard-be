package ai_estate_rag_grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
)

func (c *aiEstateRagRepository) GetAnswer(ctx context.Context, question string) (string, error) {
	res, err := c.chat.Query(ctx, &grpc_api.ChatRequest{
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
	stream, err := c.chat.UploadPdf(ctx)
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
