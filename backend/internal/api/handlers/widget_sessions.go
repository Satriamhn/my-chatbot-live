package handlers

import (
	"net/http"

	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WidgetSessionHandler struct{}

// CreateWidgetSession creates a new chat session for an anonymous customer via the public widget.
// It uses the tenant_id set by the WidgetTenantContract middleware.
func (h *WidgetSessionHandler) CreateWidgetSession(c *gin.Context) {
	tenantIDStr, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing"})
		return
	}

	orgUUID, err := uuid.Parse(tenantIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org_id"})
		return
	}

	session := models.ChatSession{
		OrganizationID: orgUUID,
		UserID:         nil, // Anonymous session
		Title:          "Customer Inquiry",
		Status:         models.SessionStatusBotHandling,
	}

	// Create session within tenant scope
	if err := db.DB.Scopes(db.TenantScope(orgUUID.String())).Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create widget session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}
