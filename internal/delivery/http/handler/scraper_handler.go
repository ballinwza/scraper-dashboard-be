package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	usecase_rental_estate "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rental_estate"
	"github.com/ballinwza/scraper-dashboard-be/pkg/validator"
)

type ScraperHandler struct {
	usecase usecase_rental_estate.IRentalEstateUsecase
}

func NewScraperHandler(usecase usecase_rental_estate.IRentalEstateUsecase) *ScraperHandler {
	return &ScraperHandler{
		usecase: usecase,
	}
}

// ScraperRentalEstate
// @Summary Trigger Scraper Job
// @Description Scrapping Manualy by Target URL
// @Tags Scraper Controls
// @Accept json
// @Produce json
// @Param request body triggerScraperRentalEstateRequest true "Target URL Payload"
// @Success 202 {object} apiResponse "Scraper job started successfully"
// @Failure 400 {object} errorResponse "Invalid request payload"
// @Failure 406 {object} errorResponse "Avaliable for admin only"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Failed to start scraper job"
// @Router /scraper/rental-estate [post]
func (h *ScraperHandler) ScraperRentalEstate(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, ok := role.(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Avaliable for admin only"})
		return
	}

	var req triggerScraperRentalEstateRequest
	validator := validator.New()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
		return
	}

	err := h.usecase.ScraperRentalEstate(c.Request.Context(), req.TargetURL, req.StartPage, req.MaxPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start scraper job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Scraper job started successfully",
	})
}
