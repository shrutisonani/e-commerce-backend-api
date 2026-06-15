package main

import (
	"auth"
	"brand"
	"category"
	"db"
	"middleware"
	"product"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shrutisonani/e-commerce-backend-api/routes"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	sqlDB := db.NewSql()

	router := gin.Default()

	userRepo := auth.NewUser(sqlDB)
	userAuth := auth.NewAuth(sqlDB)
	categoryRepo := category.NewCategory(sqlDB)
	brandRepo := brand.NewBrand(sqlDB)
	productRepo := product.NewProduct(sqlDB)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
			"db":      sqlDB != nil,
		})
	})

	api := router.Group("/api")

	// Public auth api
	routes.RegisterAuthRoutes(api, userAuth)

	// Public authenticated auth api
	protect := api.Group("")
	protect.Use(middleware.AuthMiddleware())

	routes.RegisterProtectedAuthRoutes(protect, userAuth)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// Second way to used role middleware with allowed roles
	adminRoutes := protected.Group("/v1")
	adminRoutes.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
	{
		// GEt all users for SUPER_ADMIN only
		adminRoutes.GET("/users", userRepo.Users)

		// Category routes for SUPER_ADMIN only
		categoryRoutes := adminRoutes.Group("/category")
		{
			categoryRoutes.POST("", categoryRepo.CreateCategory)
			categoryRoutes.PATCH("/:id", categoryRepo.UpdateCategory)
			categoryRoutes.DELETE("/:id", categoryRepo.DeleteCategory)
		}

		// Brand routes for SUPER_ADMIN only
		brandRoutes := adminRoutes.Group("/brand")
		{
			brandRoutes.POST("", brandRepo.CreateBrand)
			brandRoutes.PATCH("/:id", brandRepo.UpdateBrand)
			brandRoutes.DELETE("/:id", brandRepo.DeleteBrand)
		}

		// Product routes for SUPER_ADMIN only
		productRoutes := adminRoutes.Group("/product")
		{
			productRoutes.POST("", productRepo.CreateProduct)
			productRoutes.PATCH("/:id", productRepo.UpdateProduct)
			productRoutes.DELETE("/:id", productRepo.DeleteProduct)
		}
	}

	// Second way to used role middleware with allowed roles
	// If we set required role to USER then both USER and SUPER_ADMIN can access the api because of the role hierarchy
	userRoutes := protected.Group("/v2")
	userRoutes.Use(middleware.RoleMiddleware("USER"))
	{
		// User routes for both USER and SUPER_ADMIN
		user := userRoutes.Group("/user")
		{
			user.GET("/:id", userRepo.UserById)
			user.PATCH("/:id", userRepo.UpdateUser)
		}

		// Category routes for both USER and SUPER_ADMIN
		categoryRoutes := userRoutes.Group("/category")
		{
			categoryRoutes.GET("", categoryRepo.Categories)
			categoryRoutes.GET("/:id", categoryRepo.CategoryByID)
		}

		// Brand routes for both USER and SUPER_ADMIN
		brandRoutes := userRoutes.Group("/brand")
		{
			brandRoutes.GET("", brandRepo.Brands)
			brandRoutes.GET("/:id", brandRepo.BrandByID)
		}

		// Product routes for both USER and SUPER_ADMIN
		productRoutes := userRoutes.Group("/product")
		{
			productRoutes.GET("", productRepo.Products)
			productRoutes.GET("/:id", productRepo.ProductByID)
		}
	}

	return router
}

func main() {

	godotenv.Load()

	router = setupRouter()

	router.Run(":8080")
}
