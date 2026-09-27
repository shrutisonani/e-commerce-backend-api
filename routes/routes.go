package routes

import (
	"auth"
	"category"
	"db"

	"github.com/gin-gonic/gin"
	"github.com/shrutisonani/e-commerce-backend-api/brand"
	"github.com/shrutisonani/e-commerce-backend-api/middleware"
	"github.com/shrutisonani/e-commerce-backend-api/product"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	// repository
	authRepo := auth.NewAuth(h.DB)
	categoryRepo := category.NewCategory(h.DB)
	brandRepo := brand.NewBrand(h.DB)
	productRepo := product.NewProduct(h.DB, categoryRepo, brandRepo)
	variantRepo := product.NewProductVariant(h.DB)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      db.NewSql() != nil,
		})
	})

	api := router.Group("/api")

	// Public Routes
	AuthRoutes(api, authRepo)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	// USER + SUPER_ADMIN
	userRoutes := protected.Group("/v2")
	userRoutes.Use(middleware.RoleMiddleware("USER"))

	UserRoutes(userRoutes, authRepo)         // users
	CategoryRoutes(userRoutes, categoryRepo) // category
	BrandRoutes(userRoutes, brandRepo)       //brand
	ProductRoutes(userRoutes, productRepo)   // product
	VariantRoutes(userRoutes, variantRepo)   // product variants

	// SUPER_ADMIN
	adminRoutes := protected.Group("/v1")
	adminRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))

	AdminUserRoutes(adminRoutes, authRepo)         // users
	AdminCategoryRoutes(adminRoutes, categoryRepo) // category
	AdminBrandRoutes(adminRoutes, brandRepo)       // brand
	AdminProductRoutes(adminRoutes, productRepo)   // product
	AdminVariantRoutes(adminRoutes, variantRepo)   // product variants

}
