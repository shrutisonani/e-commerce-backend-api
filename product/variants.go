package product

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewProductVariant(db *sqlx.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

func (repository *Repository) CreateProductVariant(c *gin.Context) {
	input := CreateVariantRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(400)
		return
	}

	// check if the product exists and is active
	ProductID := 0
	err := repository.Db.Get(
		&ProductID,
		`
		SELECT id
		FROM products
		WHERE id = ? AND is_active = true
		`,
		input.ProductID,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product not found or inactive",
			})
			return
		}

		c.JSON(404, gin.H{"error": "Failed to validate product"})
		c.AbortWithStatus(404)
		return
	}

	// check if the selling price is greater than MRP
	if input.SellingPrice > input.MRP {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Selling price cannot be greater than MRP",
		})
		return
	}

	// calculate discount percentage
	discount := ((input.MRP - input.SellingPrice) / input.MRP) * 100

	// insert the product variant into the database
	query := `
		INSERT INTO product_variants
		(
			product_id,
			size,
			color,
			color_code,
			mrp,
			selling_price,
			discount_pct,
			stock_qty
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := repository.Db.Exec(
		query,
		input.ProductID,
		input.Size,
		input.Color,
		input.ColorCode,
		input.MRP,
		input.SellingPrice,
		discount,
		input.StockQty,
	)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create product variant"})
		c.AbortWithStatus(500)
		return
	}

	// get the last inserted ID
	id, err := result.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get variant ID",
		})
		return
	}

	// variant := ProductVariant{}

	variantExixts, err := repository.CheckVariantExists(int(id))
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{
			"error": "Failed to validate variant",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !variantExixts {
		c.JSON(400, gin.H{"error": "Variant does not exist or is inactive"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(201, gin.H{
		"message":  "Product variant created successfully",
		"variant ": input,
	})
}

// Check if a variant exists by its ID
func (repository *Repository) CheckVariantExists(variantID int) (bool, error) {
	exists := false

	err := repository.Db.Get(&exists,
		`SELECT EXISTS(SELECT 1 FROM product_variants WHERE id = ?)`,
		variantID,
	)

	return exists, err
}
