package routes

import (
	"auth"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {

	user := rg.Group("/user")

	user.GET("", authRepo.Users)            //admin
	user.GET("/:id", authRepo.UserById)     // public
	user.PATCH("/:id", authRepo.UpdateUser) //public
	// user.DELETE("/:id", authRepo.DeleteUser) //public
}
