package routes

import (
	"auth"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", authRepo.Register)
		authGroup.POST("/login", authRepo.Login)
		authGroup.GET("/verify-email", authRepo.VerifyEmail)
	}
}

func RegisterProtectedAuthRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/refresh", authRepo.RefreshToken)
		authGroup.POST("/logout", authRepo.Logout)
		authGroup.POST("/logout-all", authRepo.LogoutAll)
	}
}
