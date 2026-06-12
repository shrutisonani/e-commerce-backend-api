package product

import (
	"database/sql"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewProduct(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

// Default products
func (repository *Repository) Products(c *gin.Context) {

	productList, err := repository.GetProducts()

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, productList)
}

// create the product (Insert the data query)
func (repository *Repository) CreateProduct(c *gin.Context) {
	req := CreateProductRequest{}

	if err := c.BindJSON(&req); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO products
		(brand_id, category_id, title, description, gender_tag, material, care_instructions)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := repository.Db.Exec(query, req.BrandId, req.CategoryId, req.Title, req.Description, req.GenderTag, req.Material, req.CareInstructions)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create product"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"product": req,
	})
}

// Get all products
func (repository *Repository) GetProducts() ([]Products, error) {

	products := []Products{}

	query := `
		SELECT id, brand_id, category_id, title, description, gender_tag, material, care_instructions, is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
	`

	err := repository.Db.Select(&products, query)

	if err != nil {
		return nil, err

	}

	return products, nil
}

// Get product by ID
func (repository *Repository) ProductByID(c *gin.Context) {

	product := Products{}

	id := c.Param("id")

	query := `
		SELECT id, brand_id, category_id, title, description, gender_tag, material, care_instructions, is_active, created_at, updated_at
		FROM products
		WHERE id = ? AND is_active = true
	`
	err := repository.Db.Get(&product, query, id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update product by ID
func (repository *Repository) UpdateProduct(c *gin.Context) {
	product := Products{}

	id := c.Param("id")

	// Check if the product exists and is active before updating
	query := `
		SELECT id, brand_id, category_id, title, description, gender_tag, material, care_instructions, is_active, created_at, updated_at
		FROM products
		WHERE id = ? AND is_active = true
	`
	err := repository.Db.Get(&product, query, id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Bind the JSON input to the UpdateProductRequest struct
	input := UpdateProductRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// merge the existing product data with the input data
	if input.BrandId != nil {
		product.BrandId = *input.BrandId
	}
	if input.CategoryId != nil {
		product.CategoryId = *input.CategoryId
	}
	if input.Title != nil {
		product.Title = *input.Title
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.GenderTag != nil {
		product.GenderTag = *input.GenderTag
	}
	if input.Material != nil {
		product.Material = *input.Material
	}
	if input.CareInstructions != nil {
		product.CareInstructions = *input.CareInstructions
	}
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}

	// update the product in database
	updateQuery := `
		UPDATE products SET
			brand_id = ?, category_id = ?, title = ?, description = ?, gender_tag = ?, material = ?, care_instructions = ?, is_active = ?
		WHERE id = ? and is_active = true
	`
	_, err = repository.Db.Exec(updateQuery, product.BrandId, product.CategoryId, product.Title, product.Description, product.GenderTag, product.Material, product.CareInstructions, product.IsActive, id)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update brand"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

// Delete product by id
func (repository *Repository) DeleteProduct(c *gin.Context) {

	id := c.Param("id")

	query := "DELETE FROM product WHERE id = ? and is_active = true"

	_, err := repository.Db.Exec(query, id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
