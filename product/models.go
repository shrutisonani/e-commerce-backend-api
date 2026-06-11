package product

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Db *sqlx.DB
}

// PRODUCTS TABLE
type Products struct {
	Id               int       `json:"id"`
	BrandId          int       `json:"brand_id"`
	CategoryId       int       `json:"category_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	GenderTag        string    `json:"gender_tag"`
	Material         string    `json:"material"`
	CareInstructions string    `json:"care_instructions"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CREATE PRODUCT REQUEST
type CreateProductRequest struct {
	BrandId          int    `json:"brand_id" validate:"required"`
	CategoryId       int    `json:"category_id" validate:"required"`
	Title            string `json:"title" validate:"required"`
	Description      string `json:"description"`
	GenderTag        string `json:"gender_tag"`
	Material         string `json:"material"`
	CareInstructions string `json:"care_instructions"`
}

// UPDATE PRODUCT REQUEST
type UpdateProductRequest struct {
	BrandId          *int    `json:"brand_id"`
	CategoryId       *int    `json:"category_id"`
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	GenderTag        *string `json:"gender_tag"`
	Material         *string `json:"material"`
	CareInstructions *string `json:"care_instructions"`
	IsActive         *bool   `json:"is_active"`
}
