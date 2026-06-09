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