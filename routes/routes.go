package routes

import (
	"auth"
	"db"

	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/middleware"
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

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	// USER + SUPER_ADMIN
	userRoutes := protected.Group("/v2")
	userRoutes.Use(middleware.RoleMiddleware("USER"))

	UserRoutes(userRoutes, authRepo) // users

	// SUPER_ADMIN
	adminRoutes := protected.Group("/v1")
	adminRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))

	AdminUserRoutes(adminRoutes, authRepo) // users

}
