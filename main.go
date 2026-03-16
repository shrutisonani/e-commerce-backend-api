package main

import (
	"db"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDb := db.NewSql()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API running",
			"db": sqlDb != nil,
		})
	})

	return router
}

func main() {

	router = setupRouter()	

	router.Run(":8080")
}