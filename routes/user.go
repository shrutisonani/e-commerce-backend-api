package routes

import (
	"auth"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {

	user := rg.Group("/user")

	user.GET("/:id", authRepo.UserById)
	user.PATCH("/:id", authRepo.UpdateUser)
	// user.DELETE("/:id", authRepo.DeleteUser)
}

func AdminUserRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {

	user := rg.Group("/users")

	user.GET("", authRepo.Users) //admin
}
