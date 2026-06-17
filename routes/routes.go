package routes

import (
	"auth"
	"db"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	// repository
	authRepo := auth.NewAuth(h.DB)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      db.NewSql() != nil,
		})
	})

	api := router.Group("/api")

	// Public Routes
	AuthRoutes(api, authRepo)
}
