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

		// this conditon allows USER and SUPER_ADMIN both api access if required role is USER
		if userLevel >= requiredLevel {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
	}
}
