package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RentalEstate struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	DatePosted   string             `json:"date_posted" bson:"date_posted"`
	FormalName   string             `json:"formal_name" bson:"formal_name"`
	PropertyType string             `json:"property_type" bson:"property_type"`
	Price        float64            `json:"price" bson:"price"`
	Bedrooms     int                `json:"bedrooms" bson:"bedrooms"`
	Bathrooms    int                `json:"bathrooms" bson:"bathrooms"`
	AreaSqM      float64            `json:"area_sqm" bson:"area_sqm"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	Location     string             `json:"location" bson:"location"`
	SourceURL    string             `json:"source_url" bson:"source_url"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type RentalEstateFilter struct {
	Page     int     `json:"page" bson:"page"`
	Limit    int     `json:"limit" bson:"limit"`
	Search   string  `json:"search" bson:"search"`
	MinPrice float64 `json:"min_price" bson:"min_price"`
	MaxPrice float64 `json:"max_price" bson:"max_price"`
	SortBy   string  `json:"sort_by" bson:"sort_by"`
	Order    string  `json:"order" bson:"order"`
}
