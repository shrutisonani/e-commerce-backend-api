package main

import (
	"auth"
	"db"
	"middleware"
	"category"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDB := db.NewSql()

	router := gin.Default()

	userRepo := auth.NewUser(sqlDB)
	userAuth := auth.NewAuth(sqlDB)
	categoryRepo := category.NewCategory(sqlDB)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      sqlDB != nil,
		})
	})

	// Auth routes
	router.POST("/register", userAuth.Register)
	router.POST("/login", userAuth.Login)
	router.GET("/verify-email", userAuth.VerifyEmail)
	router.POST("/refresh", userAuth.RefreshToken)
	router.POST("/logout", userAuth.Logout)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// Second way to used role middleware with allowed roles
	adminRoutes := protected.Group("/v1")
	adminRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
	{
		// GEt all users for SUPER_ADMIN only
		adminRoutes.GET("/users", userRepo.Users)

		// Category routes for SUPER_ADMIN only
		categoryRoutes := adminRoutes.Group("/category")
		{
			categoryRoutes.POST("", categoryRepo.CreateCategory)
			categoryRoutes.PATCH("/:id", categoryRepo.UpdateCategory)
			categoryRoutes.DELETE("/:id", categoryRepo.DeleteCategory)
		}
	}

	// Second way to used role middleware with allowed roles
	// If we set required role to USER then both USER and SUPER_ADMIN can access the api because of the role hierarchy
	userRoutes := protected.Group("/v2")
	userRoutes.Use(middleware.RoleMiddleware("USER"))
	{
		// User routes for both USER and SUPER_ADMIN
		user := userRoutes.Group("/user")
		{
			user.GET("/:id", userRepo.UserById)
			user.PATCH("/:id", userRepo.UpdateUser)
			user.POST("/logout", userAuth.Logout)
			user.POST("/logout-all", userAuth.LogoutAll)
		}
		

		// Category routes for both USER and SUPER_ADMIN
		categoryRoutes := userRoutes.Group("/category")
		{
			categoryRoutes.GET("", categoryRepo.Categories)
			categoryRoutes.GET("/:id", categoryRepo.CategoryByID)
		}
	}

	return router
}

func main() {

	godotenv.Load()

	router = setupRouter()

	router.Run(":8080")
}
