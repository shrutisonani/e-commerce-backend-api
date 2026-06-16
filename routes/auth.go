package routes

import (
	"auth"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authRepo *auth.Repository) {

	authGroup := rg.Group("/auth")
	
	authGroup.POST("/register", authRepo.Register)
	authGroup.POST("/login", authRepo.Login)
	authGroup.GET("/verify-email", authRepo.VerifyEmail)

	// Public authenticated auth api
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/refresh", authRepo.RefreshToken)
	protected.POST("/logout", authRepo.Logout)
	protected.POST("/logout-all", authRepo.LogoutAll)
	
}