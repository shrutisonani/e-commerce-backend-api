package brand

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

func NewBrand(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

// Default brands
func (repository *Repository) Brands(c *gin.Context) {

	brandList, err := repository.GetBrands()

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, brandList)
}

// create the brand (Insert the data query)
func (repository *Repository) CreateBrand(c *gin.Context) {
	req := CreateBrandRequest{}

	if err := c.BindJSON(&req); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query := `
	   INSERT INTO product_brands						
	   (name, slug, logo_url)
	   VALUES (?, ?, ?)`

	_, err := repository.Db.Exec(query, req.Name, req.Slug, req.LogoURL)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create brand"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Brand created successfully",
		"brand":   req,
	})
}

// Get all brands
func (repository *Repository) GetBrands() ([]Brands, error) {

	brands := []Brands{}

	query := "SELECT id, name, slug, logo_url, is_active, created_at, updated_at FROM product_brands WHERE is_active = true"

	err := repository.Db.Select(&brands, query)

	return brands, err
}

// Check brand is exists and active
func (repository *Repository) CheckBrandExists(id int) (bool, error) {
	exists := false

	err := repository.Db.Get(&exists,
		`SELECT EXISTS(SELECT 1 FROM product_brands WHERE id = ? AND is_active = true)`,
		id,
	)

	return exists, err
}

// Get brand by id
func (repository *Repository) BrandByID(c *gin.Context) {

	brand := Brands{}

	id := c.Param("id")

	err := repository.Db.Get(&brand, "SELECT id, name, slug, logo_url, is_active, created_at, updated_at FROM product_brands WHERE id= ? and is_active = true", id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Brand not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, brand)
}

// Update brand by id
func (repository *Repository) UpdateBrand(c *gin.Context) {

	brand := Brands{}

	id := c.Param("id")

	// Check if the brand exists and is active before updating
	err := repository.Db.Get(
		&brand,
		"SELECT id, name, slug, logo_url, is_active, created_at, updated_at FROM product_brands WHERE id= ? and is_active = true", id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Brand not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Bind the JSON input to the UpdateBrandRequest struct
	input := UpdateBrandRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// merge the existing brand data with the input data
	if input.Name != nil {
		brand.Name = *input.Name
	}
	if input.Slug != nil {
		brand.Slug = *input.Slug
	}
	if input.LogoURL != nil {
		brand.LogoURL = *input.LogoURL
	}
	if input.IsActive != nil {
		brand.IsActive = *input.IsActive
	}

	// Update the brand in the database
	query := `
		UPDATE product_brands
		SET name = ?, slug = ?, logo_url = ?, is_active = ?
		WHERE id = ? and is_active = true`

	_, err = repository.Db.Exec(query, brand.Name, brand.Slug, brand.LogoURL, brand.IsActive, id)

	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update brand"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Brand updated successfully",
		"brand":   brand,
	})
}

// Delete brand by id
func (repository *Repository) DeleteBrand(c *gin.Context) {

	id := c.Param("id")

	query := "DELETE FROM product_brands WHERE id = ? and is_active = true"

	_, err := repository.Db.Exec(query, id)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Brand not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete brand"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Brand deleted successfully",
	})
}
