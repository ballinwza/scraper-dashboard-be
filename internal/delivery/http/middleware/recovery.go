package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🔥 [PANIC RECOVERED]: %v", err)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "internal server error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}