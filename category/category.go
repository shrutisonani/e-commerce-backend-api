package category

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

func NewCategory(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

// Default categories
func (repository *Repository) Categories(c *gin.Context) {

	categoryList, err := repository.GetCategories()

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, categoryList)
}

// create the category (Insert the data query)
func (repository *Repository) CreateCategory(c *gin.Context) {

	 req := CreateCategoryRequest{} 

	 if err := c.BindJSON(&req); err != nil {
		log.Error(err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		c.AbortWithStatus(http.StatusBadRequest)
		return
	 }

	 query := `
		INSERT INTO categories
		(name, description, slug, image_url, parent_id)
		VALUES (?, ?, ?, ?, ?)`

	_, err := repository.Db.Exec(query, req.Name, req.Description, req.Slug, req.ImageURL, req.ParentID)
	 if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create category"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	 }

	 c.JSON(http.StatusOK, gin.H{
		"message": "Category created successfully", 
		"category": req,
	})
}

// Get all categories
func (repository *Repository) GetCategories() ([]Categories, error) {

	categories := []Categories{}
	
	err := repository.Db.Select(&categories, "SELECT * FROM categories")
	
	return categories, err
}

// Get category by id
func (repository *Repository) CategoryByID(c *gin.Context) {

	category := Categories{}

	id := c.Param("id")

	err := repository.Db.Get(&category, "SELECT * FROM categories WHERE id= '"+id+"' and is_active = true")

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, category)	
}

// Update category by id
func (repository *Repository) UpdateCategory(c *gin.Context) {

	category := Categories{}

	id := c.Param("id")

	// Check if the category exists and is active before updating
	err := repository.Db.Get(
		&category, 
		"SELECT * FROM categories WHERE id= ? and is_active = true",
		id,
	)

	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found or inactive"})
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}	

	// Bind the JSON body to the input struct
	input := UpdateCategoryRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// merge the existing category data with the new input data
	if input.Name != nil  {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}
	if input.Slug != nil {
		category.Slug = *input.Slug
	}
	if input.ImageURL != nil {
		category.ImageURL = *input.ImageURL
	}
	if input.ParentID != nil {
		category.ParentID = input.ParentID
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	// Update the category in the database
	query := `
		UPDATE categories
		SET name = ?, description = ?, slug = ?, image_url = ?, parent_id = ?, is_active = ?
		WHERE id = ? and is_active = true`

	_, err = repository.Db.Exec(query, category.Name, category.Description, category.Slug, category.ImageURL, category.ParentID, category.IsActive, id)
	
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update category"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"category": category,
	})
}

func (repository *Repository) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM categories WHERE id = ?"	
	_, err := repository.Db.Exec(query, id)
	if err != nil {
		log.Error(err)
		c.JSON(500, gin.H{"error": "Failed to delete category"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}	