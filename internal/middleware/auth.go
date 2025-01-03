package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 权限检查
func AuthCheck(c *gin.Context) {
	// Check if the user is authenticated
	// If the user is authenticated, call c.Next()
	// If the user is not authenticated, return an error response
	c.Next()
}

// Token检查
func TokenCheck(c *gin.Context) {
	// Check if the token is valid
	// If the token is valid, call c.Next()
	// If the token is invalid, return an error response
	accessToken := c.GetHeader("Authorization")
	if accessToken != "Bearer token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
