package product

import (
	"database/sql"
	"net/http"
	"strconv"

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

	// check if the variant exists
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

// Get all product variants, optionally filtered by product ID
func (repository *Repository) GetAllVariants(c *gin.Context) {

	productID := c.Query("product_id")

	variants := []ProductVariant{}

	var err error

	if productID != "" {

		id, parseErr := strconv.Atoi(productID)

		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid product_id",
			})
			return
		}

		// check if the product exists and is active
		err = repository.Db.Select(
			&variants,
			`SELECT pv.*
			FROM product_variants pv
			INNER JOIN products p
				ON p.id = pv.product_id
			WHERE pv.product_id = ?
			AND p.is_active = true
			ORDER BY pv.id ASC`,
			id,
		)

	} else {
		// fetch all variants for active products
		err = repository.Db.Select(
			&variants,
			`SELECT pv.*
			FROM product_variants pv
			INNER JOIN products p
				ON p.id = pv.product_id
			WHERE p.is_active = true
			ORDER BY pv.id ASC`,
		)
	}

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": variants,
	})
}

// Get a product variant by its ID
func (repository *Repository) GetVariantByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid variant ID",
		})
		return
	}

	variant := ProductVariant{}

	// check if the variant exists and is active
	err = repository.Db.Get(
		&variant,
		`SELECT pv.*
		FROM product_variants pv
		INNER JOIN products p
			ON p.id = pv.product_id
		WHERE pv.id = ?
		AND p.is_active = true`,
		id,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Variant not found",
			})
			return
		}

		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"variant": variant,
	})
}

// Update a product variant by its ID
func (repository *Repository) UpdateVariant(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid variant ID",
		})
		return
	}

	variant := ProductVariant{}

	// check if the variant exists
	err = repository.Db.Get(
		&variant,
		`
		SELECT *
		FROM product_variants
		WHERE id = ?
		`,
		id,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Variant not found",
			})
			return
		}

		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch variant",
		})
		return
	}

	// bind the request body to the UpdateVariantRequest struct
	input := UpdateVariantRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if input.ProductID != nil {
		variant.ProductID = *input.ProductID
	}

	if input.Size != nil {
		variant.Size = *input.Size
	}

	if input.Color != nil {
		variant.Color = *input.Color
	}

	if input.ColorCode != nil {
		variant.ColorCode = *input.ColorCode
	}

	if input.MRP != nil {
		variant.MRP = *input.MRP
	}

	if input.SellingPrice != nil {
		variant.SellingPrice = *input.SellingPrice
	}

	if input.StockQty != nil {
		variant.StockQty = *input.StockQty
	}

	if variant.MRP <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "MRP must be greater than zero",
		})
		return
	}

	if variant.SellingPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Selling price must be greater than zero",
		})
		return
	}

	if variant.SellingPrice > variant.MRP {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Selling price cannot be greater than MRP",
		})
		return
	}

	discount := ((variant.MRP - variant.SellingPrice) / variant.MRP) * 100

	_, err = repository.Db.Exec(
		`UPDATE product_variants
		SET
			product_id = ?,
			size = ?,
			color = ?,
			color_code = ?,
			mrp = ?,
			selling_price = ?,
			discount_pct = ?,
			stock_qty = ?
		WHERE id = ?`,
		variant.ProductID,
		variant.Size,
		variant.Color,
		variant.ColorCode,
		variant.MRP,
		variant.SellingPrice,
		discount,
		variant.StockQty,
		id,
	)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update variant",
		})
		return
	}

	variant.DiscountPct = discount

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant updated successfully",
		"variant": variant,
	})
}
