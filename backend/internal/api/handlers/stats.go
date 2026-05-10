package handlers

import (
	"net/http"
	"time"

	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct{}

type DashboardStats struct {
	TotalMessagesToday int64 `json:"total_messages_today"`
	ActiveSessions     int64 `json:"active_sessions"`
	TotalUsers         int64 `json:"total_users"`
	TotalKnowledge     int64 `json:"total_knowledge"`
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	var stats DashboardStats

	// 1. Total Messages Today
	today := time.Now().Truncate(24 * time.Hour)
	db.DB.Scopes(db.TenantScope(tenantID)).
		Model(&models.Message{}).
		Where("created_at >= ?", today).
		Count(&stats.TotalMessagesToday)

	// 2. Active Sessions (Status not closed)
	db.DB.Scopes(db.TenantScope(tenantID)).
		Model(&models.ChatSession{}).
		Where("status != ?", models.SessionStatusClosed).
		Count(&stats.ActiveSessions)

	// 3. Total Users (in this organization)
	db.DB.Scopes(db.TenantScope(tenantID)).
		Model(&models.User{}).
		Count(&stats.TotalUsers)

	// 4. Total Knowledge Base Items
	db.DB.Scopes(db.TenantScope(tenantID)).
		Model(&models.KnowledgeBaseItem{}).
		Count(&stats.TotalKnowledge)

	c.JSON(http.StatusOK, stats)
}
