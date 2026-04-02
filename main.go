package main

import (
	"auth"
	"db"
	"middleware"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDb := db.NewSql()

	router := gin.Default()

	userRepo := auth.NewUser(sqlDb)
	userAuth := auth.NewAuth(sqlDb)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      sqlDb != nil,
		})
	})

	router.POST("/register", userAuth.Register)
	router.POST("login", userAuth.Login)
	router.GET("/verify-email", userAuth.VerifyEmail)
	router.POST("/refresh", userAuth.RefreshToken)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", userRepo.Users)
	}

	user := router.Group("/api/user")
	{
		user.GET("/", userRepo.Users)
		user.GET("/:id", userRepo.UserById)
		user.PATCH(":id", userRepo.UpdateUser)
	}

	return router
}

func main() {

	router = setupRouter()

	router.Run(":8080")
}
