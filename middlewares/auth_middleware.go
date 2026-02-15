package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Pastikan os.Getenv("JWT_SECRET") Mas sudah ter-set di environment atau .env
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// --- PERBAIKAN DI SINI ---
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Ekstrak UID dari claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Pastikan key "uid" sesuai dengan saat Mas melakukan GenerateToken
			uid := fmt.Sprintf("%v", claims["uid"])

			// Titipkan UID ke context agar bisa dipanggil di Controller pakai c.GetString("user_id")
			c.Set("user_id", uid)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
			return
		}

		c.Next()
	}
}
