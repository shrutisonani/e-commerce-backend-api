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
		adminRoutes.GET("/users", userRepo.Users)
	}

	// Second way to used role middleware with allowed roles
	// If we set required role to USER then both USER and SUPER_ADMIN can access the api because of the role hierarchy
	userRoutes := protected.Group("/v2")
	userRoutes.Use(middleware.RoleMiddleware("USER"))
	{
		userRoutes.GET("/user/:id", userRepo.UserById)
		userRoutes.PATCH("/user/:id", userRepo.UpdateUser)
		userRoutes.POST("/user/logout", userAuth.Logout)
		userRoutes.POST("/user/logout-all", userAuth.LogoutAll)

		categoryRoutes := userRoutes.Group("/categories")
		{
			categoryRoutes.POST("", categoryRepo.CreateCategory)
			categoryRoutes.GET("", categoryRepo.Categories)
			categoryRoutes.GET("/:id", categoryRepo.CategoryByID)
			categoryRoutes.PATCH("/:id", categoryRepo.UpdateCategory)
			categoryRoutes.DELETE("/:id", categoryRepo.DeleteCategory)
		}
	}

	return router
}

func main() {

	godotenv.Load()

	router = setupRouter()

	router.Run(":8080")
}
