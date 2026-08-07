package handler

import (
	"time"

	"github.com/ballinwza/scraper-dashboard-be/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// triggerScraperRentalEstateRequest DTO
type triggerScraperRentalEstateRequest struct {
	TargetURL string `json:"target_url" validate:"required" example:"https://www.dotproperty.co.th/en/condos-for-rent/bangkok"`
	StartPage int    `json:"start_page" example:"1"`
	MaxPage   int    `json:"max_page" example:"1"`
}

type rentalEstateResponse struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	DatePosted   string             `json:"date_posted"`
	FormalName   string             `json:"formal_name"`
	PropertyType string             `json:"property_type"`
	Price        float64            `json:"price"`
	Bedrooms     int                `json:"bedrooms"`
	Bathrooms    int                `json:"bathrooms"`
	AreaSqM      float64            `json:"area_sqm"`
	ImageURL     string             `json:"image_url"`
	Location     string             `json:"location"`
	SourceURL    string             `json:"source_url"`
	Latitude     float64            `json:"latitude"`
	Longitude    float64            `json:"longitude"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// PaginationMeta รายละเอียด Metadata ข้อมูลการแบ่งหน้า
type rentalEstatePaginationMetaResponse struct {
	CurrentPage int   `json:"current_page" example:"1"`
	Limit       int   `json:"limit" example:"10"`
	TotalItems  int64 `json:"total_items" example:"100"`
	TotalPages  int   `json:"total_pages" example:"10"`
	HasNext     bool  `json:"has_next" example:"true"`
	HasPrev     bool  `json:"has_prev" example:"false"`
}

type rentalEstatePaginatedResponse struct {
	Data       []domain.RentalEstate              `json:"data"`
	Pagination rentalEstatePaginationMetaResponse `json:"pagination"`
}
