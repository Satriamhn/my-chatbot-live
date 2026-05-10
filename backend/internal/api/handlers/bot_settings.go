package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
	"my-chatbot-backend/internal/services"
)

type BotSettingsHandler struct{}

type BotSettingsResponse struct {
	BotName        string `json:"bot_name"`
	WelcomeMessage string `json:"welcome_message"`
	SystemPrompt   string `json:"system_prompt"`
}

type FullBotSettingsResponse struct {
	BotName           string `json:"bot_name"`
	WelcomeMessage    string `json:"welcome_message"`
	SystemPrompt      string `json:"system_prompt"`
	AIProvider        string `json:"ai_provider"`
	ModelName         string `json:"model_name"`
	HasByokKey        bool   `json:"has_byok_key"`
	DailyMessageLimit int    `json:"daily_message_limit"`
	DailyMessageCount int    `json:"daily_message_count"`
}

// PublicWidgetSettingsResponse is the public widget contract: only safe display fields.
type PublicWidgetSettingsResponse struct {
	BotName        string `json:"bot_name"`
	WelcomeMessage string `json:"welcome_message"`
}

func (h *BotSettingsHandler) GetBotSettings(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant context missing"})
		return
	}

	var settings models.BotSetting
	result := db.DB.Scopes(db.TenantScope(tenantID)).First(&settings)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bot settings not found"})
		return
	}

	c.JSON(http.StatusOK, FullBotSettingsResponse{
		BotName:           settings.BotName,
		WelcomeMessage:    settings.WelcomeMessage,
		SystemPrompt:      settings.SystemPrompt,
		AIProvider:        settings.AIProvider,
		ModelName:         settings.ModelName,
		HasByokKey:        settings.HasByokKey,
		DailyMessageLimit: settings.DailyMessageLimit,
		DailyMessageCount: settings.DailyMessageCount,
	})
}

func (h *BotSettingsHandler) GetWidgetSettings(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant context missing"})
		return
	}

	var settings models.BotSetting
	result := db.DB.Scopes(db.TenantScope(tenantID)).First(&settings)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bot settings not found"})
		return
	}

	c.JSON(http.StatusOK, PublicWidgetSettingsResponse{
		BotName:        settings.BotName,
		WelcomeMessage: settings.WelcomeMessage,
	})
}

func (h *BotSettingsHandler) UpdateBotSettings(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant context missing"})
		return
	}

	var req services.BotSettingsUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return
	}

	settings, err := services.UpdateBotSettings(orgID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update bot settings"})
		return
	}

	c.JSON(http.StatusOK, FullBotSettingsResponse{
		BotName:           settings.BotName,
		WelcomeMessage:    settings.WelcomeMessage,
		SystemPrompt:      settings.SystemPrompt,
		AIProvider:        settings.AIProvider,
		ModelName:         settings.ModelName,
		HasByokKey:        settings.HasByokKey,
		DailyMessageLimit: settings.DailyMessageLimit,
		DailyMessageCount: settings.DailyMessageCount,
	})
}
