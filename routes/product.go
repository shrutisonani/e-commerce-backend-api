package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/product"
)

func ProductRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	product := rg.Group("/product")

	product.GET("", productRepo.Products)
	product.GET("/:id", productRepo.ProductByID)
}

func AdminProductRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	product := rg.Group("/product")

	product.POST("", productRepo.CreateProduct)
	product.PATCH("/:id", productRepo.UpdateProduct)
	product.DELETE("/:id", productRepo.DeleteProduct)
}
