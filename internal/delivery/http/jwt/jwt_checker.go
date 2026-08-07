package http_jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProtectedData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := c.GetString("role")

	_ = userID
	_ = role
}
