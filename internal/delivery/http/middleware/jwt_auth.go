package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		claims := &handler.AccessClaimsRequest{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(accessSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Access token expired or invalid",
			})
			c.Abort()
			return
		}

		if customClaims, ok := token.Claims.(*handler.AccessClaimsRequest); ok {
			c.Set(config.USERNAME_KEY, customClaims.Username)
			c.Set("role", customClaims.Role)
			c.Set(config.USER_ID_KEY, customClaims.UserId)

			// บันทึกลง c.Request.Context() (สำหรับ Usecase ดึงจาก ctx.Value)
			ctx := context.WithValue(c.Request.Context(), config.USERNAME_KEY, customClaims.Username)
			// ctx = context.WithValue(ctx, RoleKey, customClaims.Role)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}
