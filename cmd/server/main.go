package main

import (
	"context"
	"fmt"
	"log"
	"time"

	config "github.com/ballinwza/scraper-dashboard-be/config"
	_ "github.com/ballinwza/scraper-dashboard-be/docs"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/mongodb"
	"github.com/ballinwza/scraper-dashboard-be/internal/repository/scraper"
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

	logger.Info("🚀 Starting Scraper Backend Application...",
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

	// # Dependencies Injection setup
	mongoRealEstateRepo := mongodb.NewMongoGenericRepository[domain.RentalEstate](db)
	mongoUserRepo := mongodb.NewMongoGenericRepository[domain.User](db)
	scraperRealEstateRepo := scraper.NewDotpropertyScraperRepository()

	// Usecases
	rentalEstateUsecase := usecase_rental_estate.NewScraperRentalEstateUsecase(scraperRealEstateRepo, mongoRealEstateRepo)
	authUsecase := usecase_user.NewAuthUsecase(mongoUserRepo, cfg)

	// Handlers
	scraperHandler := handler.NewScraperHandler(rentalEstateUsecase)
	rentalHandler := handler.NewRentalEstateHandler(rentalEstateUsecase)
	authHandler := handler.NewUserHandler(authUsecase, cfg)

	// # Initialize Router
	r := http.SetupRouter(
		cfg.JwtAccessSecret,
		scraperHandler,
		rentalHandler,
		authHandler,
	)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🚀 Server running on http://localhost%s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
