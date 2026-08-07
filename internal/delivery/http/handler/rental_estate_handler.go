package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	usecase_rental_estate "github.com/ballinwza/scraper-dashboard-be/internal/usecase/rental_estate"
	"github.com/gin-gonic/gin"
)

type RentalEstateHandler struct {
	usecase usecase_rental_estate.IRentalEstateUsecase
}

func NewRentalEstateHandler(usecase usecase_rental_estate.IRentalEstateUsecase) *RentalEstateHandler {
	return &RentalEstateHandler{
		usecase: usecase,
	}
}

// RentalEstates
// @Summary ดึงรายการข้อมูลอสังหาฯ (Pagination, Filtering, Sorting)
// @Tags Rental Estates
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sort order (asc/desc)"
// @Success 200 {object} rentalEstatePaginatedResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rental/estates [get]
func (h *RentalEstateHandler) RentalEstates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	filter := domain.RentalEstateFilter{
		Page:     page,
		Limit:    limit,
		Search:   c.Query("search"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		SortBy:   c.DefaultQuery("sort_by", "created_at"),
		Order:    c.DefaultQuery("order", "desc"),
	}

	result, err := h.usecase.FetchRentalEstates(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RentalEstateById
// @Summary ดึงรายละเอียดเชิงลึกของอสังหาริมทรัพย์รายชิ้น
// @Tags Rental Estates
// @Produce json
// @Param id path string true "Estate ID"
// @Success 200 {object} rentalEstateResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rental/estate/{id} [get]
func (h *RentalEstateHandler) RentalEstateById(c *gin.Context) {
	id := c.Param("id")

	item, err := h.usecase.GetRentalEstateItem(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rental estate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ExportRentalEstates
// @Summary ส่งออกข้อมูลอสังหาฯ ที่เลือกหรือค้นหาออกมาเป็นไฟล์ CSV
// @Tags Rental Estates
// @Produce text/csv
// @Param search query string false "Search keyword"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {file} binary
// @Failure 500 {object} errorResponse
// @Router /rental/estates/export [get]
func (h *RentalEstateHandler) RentalEstateExportCSV(c *gin.Context) {
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	filter := domain.RentalEstateFilter{
		Search:   c.Query("search"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		SortBy:   c.DefaultQuery("sort_by", "created_at"),
		Order:    c.DefaultQuery("order", "desc"),
		// Limit เป็น 0 เพื่อดึงข้อมูลทั้งหมดตามเงื่อนไขค้นหา
	}

	items, err := h.usecase.ExportCSVRentalEstate(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ตั้งค่า Header สำหรับการดาวน์โหลดไฟล์ CSV
	c.Header("Content-Disposition", "attachment; filename=rental_estates_export.csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")

	writer := csv.NewWriter(c.Writer)

	// เขียน UTF-8 BOM เพื่อให้เปิดภาษาไทยใน Excel ได้ถูกต้อง
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// Header Columns
	writer.Write([]string{"ID", "Title", "Price", "Location", "Bedrooms", "Bathrooms", "Source URL", "Created At"})

	// Rows
	for _, item := range items {
		record := []string{
			item.ID.String(),
			item.Title,
			fmt.Sprintf("%.2f", item.Price),
			item.Location,
			fmt.Sprintf("%d", item.Bedrooms),
			fmt.Sprintf("%d", item.Bathrooms),
			item.SourceURL,
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	writer.Flush()
}

// DeleteRentalEstate
// @Summary ลบรายการอสังหาริมทรัพย์ที่ไม่อัปเดตหรือซ้ำซ้อน
// @Tags Rental Estates
// @Produce json
// @Param id path string true "Estate ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rental/estate/{id} [delete]
func (h *RentalEstateHandler) DeleteRentalEstateById(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.DeleteRentalEstateById(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rental estate not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rental estate deleted successfully"})
}
