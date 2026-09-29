package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/product"
)

func VariantRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	variant := rg.Group("/variant")

	variant.GET("", productRepo.GetAllVariants)
	variant.GET("/:id", productRepo.GetVariantByID)
}

func AdminVariantRoutes(rg *gin.RouterGroup, productRepo *product.Repository) {

	variant := rg.Group("/variant")

	variant.POST("", productRepo.CreateProductVariant)
	variant.PUT("/:id", productRepo.UpdateVariant)
	variant.DELETE("/:id", productRepo.DeleteVariant)
}
