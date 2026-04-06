package middleware

import (
	"fmt"
	"strconv"
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

func OwnerOrAdminMiddleware(param string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetInt("user_id")
		role := c.GetString("role")

		// Get ID from URL param (e.g. /user/:id)
		paramID, err := strconv.Atoi(c.Param(param))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid id"})
			return
		}

		fmt.Printf("UserID: %d, Role: %s, ParamID: %d\n", userID, role, paramID)
		fmt.Printf("Is Owner: %v, Is Admin: %v\n", userID == paramID, role == "SUPER_ADMIN")
		fmt.Printf("Not Admin and not Owner: %v\n", role != "SUPER_ADMIN" && userID != paramID)
		// check if user is SUPER_ADMIN or owner of the resource
		if role != "SUPER_ADMIN" && userID != paramID {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
