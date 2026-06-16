package routes

// import (
// 	"auth"

// 	"github.com/gin-gonic/gin"
// )

// func (h *Handler) RegisterRoutes(router *gin.Engine) {

// 	// repository
// 	authRepo := auth.NewRepository(h.DB)

// 	router.GET("/api/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Pong",
// 			"db":      sqlDB != nil,
// 		})
// 	})

// 	api := router.Group("/api")

// 	// Public Routes
// 	RegisterAuthRoutes(api, authRepo)
// }
