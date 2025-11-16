package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "MISSING_TOKEN"}})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_TOKEN_FORMAT"}})
			c.Abort()
			return
		}

		// TODO: Validate JWT token and extract claims
		// token := parts[1]
		// claims, err := auth.ValidateAccessToken(jwtConfig, token)
		// if err != nil {
		//     c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_TOKEN"}})
		//     c.Abort()
		//     return
		// }

		// TODO: Set user context
		// c.Set("user_id", claims.UserID)
		// c.Set("email", claims.Email)

		c.Next()
	}
}
