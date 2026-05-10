package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
)

type KnowledgeBaseHandler struct{}

func (h *KnowledgeBaseHandler) CreateItem(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	orgID, err := uuid.Parse(tenantID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return
	}

	var req struct {
		KnowledgeBaseID uuid.UUID `json:"knowledge_base_id" binding:"required"`
		Type            string    `json:"type" binding:"required"`
		Content         string    `json:"content"`
		Metadata        string    `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate type
	switch req.Type {
	case "file", "url", "manual":
		// OK
	default:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid payload type"})
		return
	}

	item := models.KnowledgeBaseItem{
		OrganizationID:  orgID,
		KnowledgeBaseID: req.KnowledgeBaseID,
		Type:            req.Type,
		Status:          "queued",
		Content:         req.Content,
		Metadata:        req.Metadata,
	}

	if err := db.DB.Scopes(db.TenantScope(tenantID.(string))).Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *KnowledgeBaseHandler) ListItems(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	kbIDStr := c.Query("knowledge_base_id")

	var items []models.KnowledgeBaseItem
	query := db.DB.Scopes(db.TenantScope(tenantID.(string)))

	if kbIDStr != "" {
		kbID, err := uuid.Parse(kbIDStr)
		if err == nil {
			query = query.Where("knowledge_base_id = ?", kbID)
		}
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *KnowledgeBaseHandler) UpdateItemStatus(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status lifecycle
	switch req.Status {
	case "queued", "processing", "ready", "failed":
		// OK
	default:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid status"})
		return
	}

	var item models.KnowledgeBaseItem
	if err := db.DB.Scopes(db.TenantScope(tenantID.(string))).First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	item.Status = req.Status
	if err := db.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, item)
}
