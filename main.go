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

	router.Group("/user")
	{
		router.POST("/register", userRepo.UserRegister)
	}

	return router
}

func main() {

	router = setupRouter()

	router.Run(":8080")
}
