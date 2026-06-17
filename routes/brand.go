package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/brand"
)

func BrandRoutes(rg *gin.RouterGroup, brandRepo *brand.Repository) {

	brand := rg.Group("/brand")

	brand.GET("", brandRepo.Brands)
	brand.GET("/:id", brandRepo.BrandByID)
}

func AdminBrandRoutes(rg *gin.RouterGroup, brandRepo *brand.Repository) {

	brand := rg.Group("/brand")

	brand.POST("", brandRepo.CreateBrand)
	brand.PATCH("/:id", brandRepo.UpdateBrand)
	brand.DELETE("/:id", brandRepo.DeleteBrand)
}
