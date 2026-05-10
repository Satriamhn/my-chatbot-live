package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("dev-secret-change-in-production")

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return JWTSecret
	}
	return []byte(secret)
}

func TenantContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 1. Extract token ────────────────────────────────────────────
		// Support both Authorization header (REST) and ?token= query param (WebSocket)
		tokenString := ""
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}
		if tokenString == "" {
			// Fallback: query param for WebSocket connections
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
			return
		}

		// ── 2. Validate JWT ─────────────────────────────────────────────
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		jwtOrgID, ok := claims["org_id"].(string)
		if !ok || jwtOrgID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing org_id claim"})
			return
		}

		// ── 3. Validate X-Org-ID (header or query param for WS) ────────
		xOrgID := c.GetHeader("X-Org-ID")
		if xOrgID == "" {
			// Fallback: query param for WebSocket
			xOrgID = c.Query("org_id")
		}

		if xOrgID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing org_id"})
			return
		}

		if jwtOrgID != xOrgID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "org_mismatch"})
			return
		}

		// ── 4. Set context values ───────────────────────────────────────
		c.Set("tenant_id", xOrgID)

		if sub, ok := claims["sub"].(string); ok && sub != "" {
			c.Set("user_id", sub)
		}

		c.Next()
	}
}
