package main

import (
	"db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shrutisonani/e-commerce-backend-api/routes"
)

var router *gin.Engine

func setupRouter() *gin.Engine {

	// mysql database connection
	sqlDB := db.NewSql()

	// Create router
	router := gin.Default()

	// Connect with routes handler
	handler := routes.NewHandler(sqlDB)

	// Call the routes register fuction
	handler.RegisterRoutes(router)

	return router
}

func main() {

	godotenv.Load()

	router = setupRouter()

	router.Run(":8080")
}
