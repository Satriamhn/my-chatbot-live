package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WidgetTenantContract ensures that a valid organization ID is provided for public widget requests.
// Unlike TenantContract, it DOES NOT require a JWT token, because widget users are anonymous.
func WidgetTenantContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get org_id from header first
		orgID := c.GetHeader("X-Org-ID")
		if orgID == "" {
			// Fallback to query parameter (often easier for simple iframe embeds)
			orgID = c.Query("org_id")
		}

		if orgID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing org_id"})
			return
		}

		// Set context value for TenantScope and other handlers to use
		c.Set("tenant_id", orgID)

		// We intentionally DO NOT set "user_id" because this is an anonymous widget session
		c.Next()
	}
}
