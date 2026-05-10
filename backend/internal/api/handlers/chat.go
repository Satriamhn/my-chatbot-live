package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"my-chatbot-backend/internal/services"
)

// ChatHandler handles API requests for chat sessions and messages.
type ChatHandler struct{}

// CreateSession creates a new chat session.
// POST /api/v1/sessions
func (h *ChatHandler) CreateSession(c *gin.Context) {
	orgID, err := getOrgID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var payload services.CreateSessionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if payload.Title == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Title is required"})
		return
	}

	// Extract user ID from context if available
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uID, err := uuid.Parse(userIDStr); err == nil {
			payload.UserID = &uID
		}
	}

	session, err := services.CreateSession(orgID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetSession fetches a single chat session by ID.
// GET /api/v1/sessions/:id
func (h *ChatHandler) GetSession(c *gin.Context) {
	orgID, err := getOrgID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := services.GetSession(orgID, sessionID)
	if err != nil {
		if err.Error() == "org_mismatch" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization mismatch"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// AddMessage adds a new message to a chat session.
// POST /api/v1/sessions/:id/messages
func (h *ChatHandler) AddMessage(c *gin.Context) {
	orgID, err := getOrgID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var payload services.CreateMessagePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if payload.Sender == "" || payload.Content == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Sender and Content are required"})
		return
	}

	message, err := services.AddMessage(orgID, sessionID, payload)
	if err != nil {
		if err.Error() == "org_mismatch" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization mismatch"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add message"})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// TakeoverSession transitions a session to human mode.
// POST /api/v1/sessions/:id/takeover
func (h *ChatHandler) TakeoverSession(c *gin.Context) {
	orgID, err := getOrgID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = services.TakeoverSession(orgID, sessionID)
	if err != nil {
		switch err.Error() {
		case "org_mismatch":
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization mismatch"})
		case "invalid_transition":
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid session transition"})
		case "record not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to takeover session"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Human takeover successful"})
}

// ReturnToBotMode transitions a session back to bot mode.
// POST /api/v1/sessions/:id/return
func (h *ChatHandler) ReturnToBotMode(c *gin.Context) {
	orgID, err := getOrgID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = services.ReturnToBotMode(orgID, sessionID)
	if err != nil {
		switch err.Error() {
		case "org_mismatch":
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization mismatch"})
		case "invalid_transition":
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid session transition"})
		case "record not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to return to bot mode"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Returned to bot mode successfully"})
}

// getOrgID extracts the organization ID from the context.
func getOrgID(c *gin.Context) (uuid.UUID, error) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		return uuid.Nil, http.ErrNoLocation // Just a placeholder error
	}
	return uuid.Parse(tenantID)
}
