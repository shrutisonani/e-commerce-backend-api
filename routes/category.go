package routes

func CategoryRoutes(rg *gin.RouterGroup, categoryRepo *category.Repository) {

	category := rg.Group("/category")

	category.GET("", categoryRepo.Categories)
	category.GET("/:id", categoryRepo.CategoryByID)
}

func AdminCategoryRoutes(rg *gin.RouterGroup, categoryRepo *category.Repository) {

	category := rg.Group("/category")

	category.POST("", categoryRepo.CreateCategory)
	category.PATCH("/:id", categoryRepo.UpdateCategory)
	category.DELETE("/:id", categoryRepo.DeleteCategory)
}