package middleware

import (
	"fmt"
	"net/http"

	"github.com/ballinwza/scraper-dashboard-be/config"
	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	"github.com/ballinwza/scraper-dashboard-be/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTRefreshMiddleware(refreshSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. ลองดึง refresh_token จาก Cookie
		cookieToken, err := c.Cookie(config.COOKIE_REFRESH_TOKEN_KEY)
		if err == nil && cookieToken != "" {
			tokenString = cookieToken
		}

		// หากหา Token ไม่เจอเลย
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "refresh_token cookie or authorization header is required",
			})
			c.Abort()
			return
		}

		// 3. Parse และ Validate JWT Signature ด้วย refreshSecret
		claims := &handler.RefreshClaimsRequest{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(refreshSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "refresh token expired or invalid",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*handler.RefreshClaimsRequest); ok {
			c.Set(config.USERNAME_KEY, claims.Username)
			tokenHash := helper.HashTokenSHA256(tokenString)
			c.Set(config.REFRESH_HASH_TOKEN_KEY, tokenHash)
		}

		c.Next()
	}
}
