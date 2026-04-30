package middleware

import (
	"net/http"
	"strings"
	"uniS/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "")
		claims, err := jwt.Parse(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("open_id", claims.OpenID)
		c.Next()
	}
}
