package external_grpc

import (
	"context"
	"fmt"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IAiEstateRagRepository interface {
	GetAnswer(ctx context.Context, question string) (string, error)
}

type aiEstateRagRepository struct {
	client grpc_api.ChatGRPCClient
}

// NewAiEstateRagRepository สร้าง gRPC Client Connection ไปยัง port 50052
func NewAiEstateRagRepository(targetAddr string) (IAiEstateRagRepository, func(), error) {
	if targetAddr == "" {
		targetAddr = "localhost:50052" // หรือ "grpc-service:50052" สำหรับ Docker
	}

	// สร้าง Connection ไปยัง gRPC Target
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to external gRPC on %s: %w", targetAddr, err)
	}

	// ฟังก์ชัน Cleanup สำหรับสั่ง Close Connection เมื่อแอปปิด
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
