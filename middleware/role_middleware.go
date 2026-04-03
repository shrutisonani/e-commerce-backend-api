package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var RoleHierarchy = map[string]int{
	"USER":        1,
	"SUPER_ADMIN": 2,
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleInterface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "role not found"})
			return
		}

		userRole := strings.ToUpper(strings.TrimSpace(roleInterface.(string)))

		userLevel := RoleHierarchy[userRole]
		requiredLevel := RoleHierarchy[strings.ToUpper(requiredRole)]

		fmt.Printf("UserRole: %s (%d), Required: %s (%d)\n",
			userRole, userLevel, requiredRole, requiredLevel)

		// MAIN LOGIC
		if userLevel >= requiredLevel {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
	}
}

// Only user access
func UserOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "role not found"})
			return
		}

		if role != "user" {
			c.AbortWithStatusJSON(403, gin.H{"error": "user access only"})
			return
		}

		c.Next()
	}
}

// admin access only
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "role not found"})
			return
		}

		if role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "admin access only"})
			return
		}

		c.Next()
	}
}
