package http

import (
	"time"

	"github.com/ballinwza/scraper-dashboard-be/config"
	handler "github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	middleware "github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	cfg config.Config,
	scraperHandler *handler.ScraperHandler,
	rentalEstateHandler *handler.RentalEstateHandler,
	userHandler *handler.UserHandler,
	ragHandler *handler.RagHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RequestLogger())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", cfg.FEUri},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
			authGroup.POST("/logout", middleware.JWTRefreshMiddleware(cfg.JwtRefreshSecret), userHandler.Logout)
			authGroup.POST("/refresh", middleware.JWTRefreshMiddleware(cfg.JwtRefreshSecret), userHandler.Refresh)
		}

		userGroup := api.Group("/user")
		userGroup.Use(middleware.JWTAuthMiddleware(cfg.JwtAccessSecret))
		{
			userGroup.GET("/username", userHandler.GetUser)
		}

		scraperGroup := api.Group("/scraper")
		scraperGroup.Use(middleware.JWTAuthMiddleware(cfg.JwtAccessSecret))
		{
			scraperGroup.POST("/rental-estate", scraperHandler.ScraperRentalEstate)
		}

		rentalEstateGroup := api.Group("/rental")
		rentalEstateGroup.Use(middleware.JWTAuthMiddleware(cfg.JwtAccessSecret))
		{
			rentalEstateGroup.GET("/estates", rentalEstateHandler.RentalEstates)
			rentalEstateGroup.GET("/estates/export", rentalEstateHandler.RentalEstateExportCSV)
			rentalEstateGroup.GET("/estate/:id", rentalEstateHandler.RentalEstateById)
			rentalEstateGroup.DELETE("/estate/:id", rentalEstateHandler.DeleteRentalEstateById)
		}

		ragGroup := api.Group("/rag")
		ragGroup.Use(middleware.JWTAuthMiddleware(cfg.JwtAccessSecret))
		{
			ragGroup.POST("/qna", ragHandler.AskQuestion)
			ragGroup.POST("/upload", ragHandler.RecieverOfUplodFile)
			ragGroup.POST("/multi/upload", ragHandler.UploadFileMultiTenant)
		}

	}

	return r
}
