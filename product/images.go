package product

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewProductImage(db *sqlx.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

// CreateImage handles the creation of a new product image. It validates the input, checks if the product and variant (if provided) exist, and inserts the new image into the database.
func (repository *Repository) CreateImage(c *gin.Context) {

	input := CreateImageRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Check product
	var productID int

	err := repository.Db.Get(
		&productID,
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

		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate product",
		})
		return
	}

	// If variant is supplied, verify it belongs to the product
	if input.VariantID != 0 {

		var variantID int

		err = repository.Db.Get(
			&variantID,
			`
			SELECT id
			FROM product_variants
			WHERE id = ?
			AND product_id = ?
			`,
			input.VariantID,
			input.ProductID,
		)

		if err != nil {

			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Variant does not belong to this product",
				})
				return
			}

			log.Error(err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to validate variant",
			})
			return
		}
	}

	result, err := repository.Db.Exec(
		`
		INSERT INTO product_images
		(
			product_id,
			image_url,
			variant_id,
			all_text,
			display_order
		)
		VALUES (?, ?, ?, ?, ?)
		`,
		input.ProductID,
		input.ImageURL,
		input.VariantID,
		input.AllText,
		input.DisplayOrder,
	)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product image",
		})
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get image ID",
		})
		return
	}

	var image ProductImages

	err = repository.Db.Get(
		&image,
		`
		SELECT *
		FROM product_images
		WHERE id = ?
		`,
		id,
	)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch created image",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product image created successfully",
		"image":   image,
	})
}
