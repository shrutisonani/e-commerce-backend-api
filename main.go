package main

import (
	"auth"
	"db"
	"middleware"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDB := db.NewSql()

	router := gin.Default()

	userRepo := auth.NewUser(sqlDB)
	userAuth := auth.NewAuth(sqlDB)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      sqlDB != nil,
		})
	})

	// Auth routes
	router.POST("/register", userAuth.Register)
	router.POST("login", userAuth.Login)
	router.GET("/verify-email", userAuth.VerifyEmail)
	router.POST("/refresh", userAuth.RefreshToken)
	router.POST("/logout", userAuth.Logout)

	// Protected routes
	protected := router.Group("/api")

	// Apply auth middleware to all protected routes
	protected.Use(middleware.AuthMiddleware())

	// Second way to used role middleware with allowed roles
	adminRoutes := protected.Group("/")
	adminRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
	{
		adminRoutes.GET("/users", userRepo.Users)
		protected.POST("/logout-all", userAuth.LogoutAll)
	}

	// Second way to used role middleware with allowed roles
	userRoutes := protected.Group("/")
	userRoutes.Use(middleware.RoleMiddleware("USER"))
	{
		userRoutes.GET("/user/:id", userRepo.UserById)
		userRoutes.PATCH("/user/:id", userRepo.UpdateUser)
	}

	return router
}

func main() {

	router = setupRouter()

	router.Run(":8080")
}
