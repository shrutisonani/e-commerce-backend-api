package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/product"
)

func ImageRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	// image := rg.Group("/image")
}

func AdminImageRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {
	image := rg.Group("/image")

	image.POST("", productRepo.CreateImage)

}
