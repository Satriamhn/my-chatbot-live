package handlers

import (
	"net/http"

	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatSessionHandler struct{}

func (h *ChatSessionHandler) CreateSession(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userIDStr := c.GetString("user_id")

	var input struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgUUID, _ := uuid.Parse(tenantID)
	userUUID, _ := uuid.Parse(userIDStr)

	session := models.ChatSession{
		OrganizationID: orgUUID,
		UserID:         &userUUID,
		Title:          input.Title,
		Status:         models.SessionStatusBotHandling,
	}

	if err := db.DB.Scopes(db.TenantScope(tenantID)).Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *ChatSessionHandler) ListSessions(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var sessions []models.ChatSession

	if err := db.DB.Scopes(db.TenantScope(tenantID)).Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *ChatSessionHandler) TakeoverSession(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	sessionID := c.Param("id")

	var session models.ChatSession
	if err := db.DB.Scopes(db.TenantScope(tenantID)).First(&session, "id = ?", sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	if !models.IsValidSessionTransition(session.Status, models.SessionStatusHumanAssigned) {
		c.JSON(http.StatusConflict, gin.H{"error": "invalid state transition"})
		return
	}

	session.Status = models.SessionStatusHumanAssigned
	if err := db.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session status"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *ChatSessionHandler) ReturnToBotMode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	sessionID := c.Param("id")

	var session models.ChatSession
	if err := db.DB.Scopes(db.TenantScope(tenantID)).First(&session, "id = ?", sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	session.Status = models.SessionStatusBotHandling
	if err := db.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session status"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *ChatSessionHandler) AddMessage(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	sessionID := c.Param("id")

	var input struct {
		Content    string `json:"content" binding:"required"`
		SenderType string `json:"sender_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Get session and verify tenant (TenantScope handles this)
	var session models.ChatSession
	if err := db.DB.Scopes(db.TenantScope(tenantID)).First(&session, "id = ?", sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found or access denied"})
		return
	}

	// 2. Enforce human takeover rule: if status is human_assigned, bot cannot send messages
	if session.Status == models.SessionStatusHumanAssigned && input.SenderType == models.SenderBot {
		c.JSON(http.StatusForbidden, gin.H{"error": "human takeover active: bot messages blocked"})
		return
	}

	orgUUID, _ := uuid.Parse(tenantID)
	sessionUUID, _ := uuid.Parse(sessionID)

	message := models.Message{
		OrganizationID: orgUUID,
		ChatSessionID:  sessionUUID,
		Sender:         input.SenderType,
		Content:        input.Content,
	}

	if err := db.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *ChatSessionHandler) ListMessages(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	sessionID := c.Param("id")

	// Verify session exists for this tenant
	var session models.ChatSession
	if err := db.DB.Scopes(db.TenantScope(tenantID)).First(&session, "id = ?", sessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	var messages []models.Message
	if err := db.DB.Scopes(db.TenantScope(tenantID)).Where("chat_session_id = ?", sessionID).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
