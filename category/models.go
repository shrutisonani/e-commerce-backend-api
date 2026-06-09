package category

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Db *sqlx.DB
}

// CATEGORIES TABLE
type Categories struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	ImageURL    string    `json:"image_url"`
	ParentID    *int      `json:"parent_id"` // Nullable for top-level categories
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CREATE CATEGORY REQUEST
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug" validate:"required"`
	ImageURL    string `json:"image_url"`
	ParentID    *int   `json:"parent_id"` // Optional for top-level categories
}

// UPDATE CATEGORY REQUEST
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	ImageURL    *string `json:"image_url"`
	ParentID    *int    `json:"parent_id"`
	IsActive    *bool   `json:"is_active"`
}
