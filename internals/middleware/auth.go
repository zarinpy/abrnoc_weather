package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zarinpy/abrnoc_weather/pkg/respond"
)

// Claims represents the expected JWT payload.
type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.RegisteredClaims
}

const ContextUserKey = "user"

// AuthRequired validates JWT tokens and injects claims into the context.
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respond.Error(c, http.StatusUnauthorized, "missing_auth_header", "Authorization header required")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			respond.Error(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			respond.Error(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
			return
		}

		c.Set(ContextUserKey, claims)
		c.Next()
	}
}
