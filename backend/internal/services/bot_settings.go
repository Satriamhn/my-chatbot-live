package services

import (
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/google/uuid"
)

type BotSettingsUpdate struct {
	BotName        string `json:"bot_name" validate:"required,max=255"`
	WelcomeMessage string `json:"welcome_message" validate:"required"`
	SystemPrompt   string `json:"system_prompt" validate:"required"`
	AIProvider     string `json:"ai_provider"`
	ModelName      string `json:"model_name"`
	APIKey         string `json:"api_key"`
}

func UpdateBotSettings(orgID uuid.UUID, payload BotSettingsUpdate) (*models.BotSetting, error) {
	var settings models.BotSetting
	tenantID := orgID.String()

	// Default provider/model if not provided
	if payload.AIProvider == "" {
		payload.AIProvider = models.AIProviderGemini
	}
	if payload.ModelName == "" {
		payload.ModelName = DefaultModelFor(payload.AIProvider)
	}

	result := db.DB.Scopes(db.TenantScope(tenantID)).Where("organization_id = ?", tenantID).First(&settings)

	if result.Error != nil {
		// Create new
		settings = models.BotSetting{
			OrganizationID:    orgID,
			BotName:           payload.BotName,
			WelcomeMessage:    payload.WelcomeMessage,
			SystemPrompt:      payload.SystemPrompt,
			AIProvider:        payload.AIProvider,
			ModelName:         payload.ModelName,
			DailyMessageLimit: models.DefaultDailyLimit,
		}
		if payload.APIKey != "" {
			settings.APIKey = payload.APIKey
			settings.HasByokKey = true
		}
		if err := db.DB.Scopes(db.TenantScope(tenantID)).Create(&settings).Error; err != nil {
			return nil, err
		}
	} else {
		// Update existing
		settings.BotName = payload.BotName
		settings.WelcomeMessage = payload.WelcomeMessage
		settings.SystemPrompt = payload.SystemPrompt
		settings.AIProvider = payload.AIProvider
		settings.ModelName = payload.ModelName
		if payload.APIKey != "" {
			settings.APIKey = payload.APIKey
			settings.HasByokKey = true
		}
		if err := db.DB.Scopes(db.TenantScope(tenantID)).Save(&settings).Error; err != nil {
			return nil, err
		}
	}

	return &settings, nil
}
