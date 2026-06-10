package brand

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Db *sqlx.DB
}

// BRANDS TABLE
type Brands struct {
	Id          int       `json:"id"`	
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	LogoURL     string    `json:"logo_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CREATE BRAND REQUEST
type CreateBrandRequest struct {
	Name    string `json:"name" validate:"required"`
	Slug    string `json:"slug" validate:"required"`
	LogoURL string `json:"logo_url"`
}

// UPDATE BRAND REQUEST
type UpdateBrandRequest struct {
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
	LogoURL  *string `json:"logo_url"`
	IsActive *bool   `json:"is_active"`
}