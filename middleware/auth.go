package middleware

import (
	"net/http"
	"strings"

	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authConfig *config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authConfig.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

func GetCurrentClaims(c *gin.Context) *config.JWTClaims {
	claims, exists := c.Get("claims")
	if !exists {
		return nil
	}
	return claims.(*config.JWTClaims)
}