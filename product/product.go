package product

import (
	"category"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/shrutisonani/e-commerce-backend-api/brand"
	log "github.com/sirupsen/logrus"
)

func NewProduct(
	db *sqlx.DB,
	categoryRepo *category.Repository,
	brandRepo *brand.Repository,
) *Repository {

	return &Repository{
		Db:       db,
		Category: categoryRepo,
		Brand:    brandRepo,
	}
}

// Default products
// func (repository *Repository) Products(c *gin.Context) {

// 	productList, err := repository.GetProducts()

// 	if err != nil {
// 		log.Error(err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

// 	c.JSON(http.StatusOK, productList)
// }

// create the product (Insert the data query)
func (repository *Repository) CreateProduct(c *gin.Context) {
	req := CreateProductRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check if the brand exists
	brandExists, err := repository.Brand.CheckBrandExists(req.BrandId)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{
			"error": "Failed to validate brand",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !brandExists {
		c.JSON(400, gin.H{"error": "Brand does not exist or is inactive"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check if the category exists
	categoryExists, err := repository.Category.CheckCategoryExists(req.CategoryId)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{
			"error": "Failed to validate category",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !categoryExists {
		c.JSON(400, gin.H{"error": "Category does not exist or is inactive"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO products
		(brand_id, category_id, title, description, gender_tag, material, care_instructions)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := repository.Db.Exec(
		query,
		req.BrandId,
		req.CategoryId,
		req.Title,
		req.Description,
		req.GenderTag,
		req.Material,
		req.CareInstructions,
	)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create product"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// productId, err = result.LastInsertId()
	_, err = result.LastInsertId()

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get product ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"product": req,
	})
}

// Get all products
func (repository *Repository) GetProducts(c *gin.Context) {

	products := []Products{}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	search := strings.TrimSpace(
		c.DefaultQuery("search", ""),
	)

	query := ""
	args := []interface{}{}

	if search != "" {
		query = `
		SELECT *
			FROM products
			WHERE is_active = true
			AND (
				title LIKE ?
				OR description LIKE ?
			)
			ORDER BY id DESC
			LIMIT ? OFFSET ?
	`
		searchValue := "%" + search + "%"

		args = []interface{}{
			searchValue,
			searchValue,
			limit,
			offset,
		}
	} else {
		query = `
			SELECT *
			FROM products
			WHERE is_active = true
			ORDER BY id DESC
			LIMIT ? OFFSET ?
		`

		args = []interface{}{
			limit,
			offset,
		}
	}
	err := repository.Db.Select(&products, query, args...)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch products",
		})
		return
	}

	total := 0
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE is_active = true
	`

	if search != "" {

		searchValue := "%" + search + "%"

		err = repository.Db.Get(
			&total,
			`
			SELECT COUNT(*)
			FROM products
			WHERE is_active = true
			AND (
				title LIKE ?
				OR description LIKE ?
			)
			`,
			searchValue,
			searchValue,
		)

	} else {

		err = repository.Db.Get(
			&total,
			countQuery,
		)
	}

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// Get product by ID
func (repository *Repository) ProductByID(c *gin.Context) {

	product := Products{}

	id := c.Param("id")

	query := `
		SELECT *
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

	variants := []ProductVariant{}

	err = repository.Db.Select(
		&variants,
		`SELECT *
		 FROM product_variants
		 WHERE product_id = ?
		 ORDER BY id ASC`,
		id,
	)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product variants",
		})
		return
	}

	images := []ProductImages{}

	err = repository.Db.Select(
		&images,
		`SELECT *
		 FROM product_images
		 WHERE product_id = ?
		 ORDER BY display_order ASC, id ASC`,
		id,
	)

	if err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product images",
		})
		return
	}

	response := ProductDetailResponse{
		Product:  product,
		Variants: variants,
		Images:   images,
	}

	c.JSON(http.StatusOK, response)
}

// Update product by ID
func (repository *Repository) UpdateProduct(c *gin.Context) {
	product := Products{}

	// id := c.Param("id")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Check if the product exists and is active before updating
	err = repository.Db.Get(&product,
		`
		SELECT *
		FROM products
		WHERE id = ? AND is_active = true
	`, id,
	)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found or inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	// Bind the JSON input to the UpdateProductRequest struct
	input := UpdateProductRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
		UPDATE products 
		SET
			brand_id = ?, 
			category_id = ?, 
			title = ?, 
			description = ?, 
			gender_tag = ?, 
			material = ?, 
			care_instructions = ?, 
			is_active = ?
		WHERE id = ?
	`
	_, err = repository.Db.Exec(
		updateQuery,
		product.BrandId,
		product.CategoryId,
		product.Title,
		product.Description,
		product.GenderTag,
		product.Material,
		product.CareInstructions,
		product.IsActive,
		id,
	)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

// Delete product by id
func (repository *Repository) DeleteProduct(c *gin.Context) {

	// id := c.Param("id")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	result, err := repository.Db.Exec(
		`
		UPDATE products
		SET is_active = false
		WHERE id = ? AND is_active = true
		`,
		id,
	)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	rows, err := result.RowsAffected()

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify deletion"})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or already inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
