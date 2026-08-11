package main

import (
	"context"
	"fmt"
	"log"
	"time"

	config "github.com/ballinwza/scraper-dashboard-be/config"
	_ "github.com/ballinwza/scraper-dashboard-be/docs"
	http "github.com/ballinwza/scraper-dashboard-be/internal/delivery"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/external_grpc"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/mongodb"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/scraper"
	usecase_rag "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rag"
	usecase_rental_estate "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rental_estate"
	usecase_user "github.com/ballinwza/scraper-dashboard-be/internal/usecase/user"
	"github.com/ballinwza/scraper-dashboard-be/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @title           Scraper & Data Aggregator API
// @version         1.0
// @description     API Service สำหรับจัดการข้อมูล Scraping และแสดงผลบน Dashboard
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// # Initialize ENV
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("⚙️ App initialized in [%s] mode", cfg.Environment)

	// # Initialize Logger
	logger.InitLogger(cfg.Environment)
	defer logger.Sync()

	logger.Info(
		"🚀 Starting Scraper Backend Application...",
		zap.String("env", cfg.Environment),
		zap.String("port", cfg.ServerPort),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatal("MongoDB ping failed", zap.Error(err))
	}

	db := client.Database("scraper_dashboard")

	// # GRPC
	grpcConn, grpcCleanup, err := external_grpc.NewAiEstateRagRepository(cfg.AiEstateRagUri)
	if err != nil {
		logger.Fatal("Failed to initialize external gRPC", zap.Error(err))
	}
	defer grpcCleanup()

	// # Dependencies Injection setup
	mongoRealEstateRepo := mongodb.NewMongoGenericRepository[domain.RentalEstate](db)
	mongoUserRepo := mongodb.NewMongoGenericRepository[domain.User](db)
	scraperRealEstateRepo := scraper.NewDotpropertyScraperRepository()

	// Usecases
	rentalEstateUsecase := usecase_rental_estate.NewScraperRentalEstateUsecase(scraperRealEstateRepo, mongoRealEstateRepo)
	authUsecase := usecase_user.NewAuthUsecase(mongoUserRepo, cfg)
	ragUsecase := usecase_rag.NewRagUsecase(grpcConn)

	// Handlers
	scraperHandler := handler.NewScraperHandler(rentalEstateUsecase)
	rentalHandler := handler.NewRentalEstateHandler(rentalEstateUsecase)
	authHandler := handler.NewUserHandler(authUsecase, cfg)
	ragHandler := handler.NewRagHandler(ragUsecase, cfg)

	// ==========================================
	// ✨ Start gRPC Server (Non-blocking via Goroutine)
	// ==========================================
	// TODO: เก็บไว้เป็นตัวอย่าง GRPC
	// go func() {
	// 	grpcPort := ":50051" // หรือดึงมาจาก config เช่น cfg.GRPCPort
	// 	lis, err := net.Listen("tcp", grpcPort)
	// 	if err != nil {
	// 		logger.Fatal("Failed to listen for gRPC", zap.Error(err))
	// 	}

	// 	// สร้าง gRPC Server (นำ Usecase ไปใช้ใน gRPC Handler ได้ตามต้องการ)
	// 	grpcServer := appGrpc.NewGRPCServer(ragUsecase)

	// 	log.Printf("📡 gRPC Server running on port %s", grpcPort)
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		logger.Fatal("Failed to serve gRPC", zap.Error(err))
	// 	}
	// }()

	// ==========================================
	// # Initialize & Start HTTP Router (Main Thread)
	// ==========================================
	r := http.SetupRouter(
		cfg.JwtAccessSecret,
		scraperHandler,
		rentalHandler,
		authHandler,
		ragHandler,
	)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server running on http://localhost%s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
