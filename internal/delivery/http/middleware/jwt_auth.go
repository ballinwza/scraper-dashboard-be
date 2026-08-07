package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ballinwza/scraper-dashboard-be/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// บันทึก Claims ไว้ใน Context สำหรับ Handler ถัดไปนำไปใช้งาน
		if claims, ok := token.Claims.(handler.AccessClaimsRequest); ok {
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		}

		// บันทึกใน context
		// type contextKey string
		// const UserIDKey contextKey = "username"

		// if customClaims, ok := token.Claims.(*handler.AccessClaimsRequest); ok {
		// 	// แปะค่าลง stdlib context.Context
		// 	ctx := context.WithValue(c.Request.Context(), UserIDKey, customClaims.Username)
		// 	c.Request = c.Request.WithContext(ctx)
		// }

		// วิธีเรียกใช้ใน usecase โดยตรง
		// userID, ok := ctx.Value(UserIDKey).(string)
		//     if !ok {
		//         return errors.New("unauthorized")
		//     }

		c.Next()
	}
}
