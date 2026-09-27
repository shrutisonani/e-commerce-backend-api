package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/product"
)

func VariantRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	// variant := rg.Group("/variant")

	// variant.GET("", productRepo.GetProductVariants)
}

func AdminVariantRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	variant := rg.Group("/variant")

	variant.POST("", productRepo.CreateProductVariant)
}
