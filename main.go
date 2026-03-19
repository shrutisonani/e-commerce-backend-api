package main

import (
	"db"
	"user"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDb := db.NewSql()

	router := gin.Default()

	userRepo := user.NewUser(sqlDb)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      sqlDb != nil,
		})
	})

	user := router.Group("/api/user")
	{
		user.POST("/register", userRepo.UserRegister)
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
