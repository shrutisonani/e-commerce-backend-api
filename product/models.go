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

	Variants []ProductVariant `json:"variants,omitempty"`
	Images   []ProductImages  `json:"images,omitempty"`
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

// PRODUCT VARIANTS TABLE
type ProductVariant struct {
	ID           int     `db:"id"            json:"id"`
	ProductID    int     `db:"product_id"    json:"product_id"`
	Size         string  `db:"size"          json:"size"`
	Color        string  `db:"color"         json:"color"`
	ColorCode    string  `db:"color_code"    json:"color_code"`
	MRP          float64 `db:"mrp"           json:"mrp"`
	SellingPrice float64 `db:"selling_price" json:"selling_price"`
	DiscountPct  float64 `db:"discount_pct"  json:"discount_pct"`
	StockQty     int     `db:"stock_qty"     json:"stock_qty"`
}

// CREATE VARIANT REQUEST
type CreateVariantRequest struct {
	ProductID    int     `json:"product_id" validate:"required"`
	Size         string  `json:"size" validate:"required"`
	Color        string  `json:"color" validate:"required"`
	ColorCode    string  `json:"color_code" validate:"required"`
	MRP          float64 `json:"mrp" validate:"required,gt=0"`
	SellingPrice float64 `json:"selling_price" validate:"required,gt=0"`
	StockQty     int     `json:"stock_qty" validate:"required"`
}

// UPDATE VARIANT REQUEST
type UpdateVariantRequest struct {
	ProductID    *int     `json:"product_id" validate:"required"`
	Size         *string  `json:"size" validate:"required"`
	Color        *string  `json:"color" validate:"required"`
	ColorCode    *string  `json:"color_code" validate:"required"`
	MRP          *float64 `json:"mrp" validate:"gt=0"`
	SellingPrice *float64 `json:"selling_price" validate:"gt=0"`
	StockQty     *int     `json:"stock_qty" validate:"required"`
}

// PRODUCT IMAGES TABLE
type ProductImages struct {
	ID           int    `db:"id"            json:"id"`
	ProductID    int    `db:"product_id"    json:"product_id"`
	ImageURL     string `db:"image_url"     json:"image_url"`
	VariantID    int    `db:"variant_id"    json:"variant_id"`
	AllText      string `db:"all_text"      json:"all_text"`
	DisplayOrder int    `db:"display_order" json:"display_order"`
}

// CREATE IMAGE REQUEST
type CreateImageRequest struct {
	ProductID    int    `json:"product_id"    validate:"required"`
	ImageURL     string `json:"image_url"     validate:"required"`
	VariantID    int    `json:"variant_id"`
	AllText      string `json:"all_text"`
	DisplayOrder int    `json:"display_order"`
}
