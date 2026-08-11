package http

import (
	handler "github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	middleware "github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	jwtSecret string,
	scraperHandler *handler.ScraperHandler,
	rentalEstateHandler *handler.RentalEstateHandler,
	authHandler *handler.UserHandler,
	ragHandler *handler.RagHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RequestLogger())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
		}

		scraperGroup := api.Group("/scraper")
		scraperGroup.Use(middleware.JWTAuthMiddleware(jwtSecret))
		{
			scraperGroup.POST("/rental-estate", scraperHandler.ScraperRentalEstate)
		}

		rentalEstateGroup := api.Group("/rental")
		rentalEstateGroup.Use(middleware.JWTAuthMiddleware(jwtSecret))
		{
			rentalEstateGroup.GET("/estates", rentalEstateHandler.RentalEstates)
			rentalEstateGroup.GET("/estates/export", rentalEstateHandler.RentalEstateExportCSV)
			rentalEstateGroup.GET("/estate/:id", rentalEstateHandler.RentalEstateById)
			rentalEstateGroup.DELETE("/estate/:id", rentalEstateHandler.DeleteRentalEstateById)
		}

		ragGroup := api.Group("/rag")
		ragGroup.Use(middleware.JWTAuthMiddleware(jwtSecret))
		{
			ragGroup.POST("/qna", ragHandler.AskQuestion)
		}

	}

	return r
}
