package category

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

func NewCategory(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}

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