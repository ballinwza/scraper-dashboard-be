package external_grpc

import (
	"context"
	"crypto/x509"
	"fmt"

	grpc_api "github.com/ballinwza/scraper-dashboard-be/internal/delivery/grpc/api"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type IAiEstateRagRepository interface {
	GetAnswer(ctx context.Context, question string) (string, error)
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
