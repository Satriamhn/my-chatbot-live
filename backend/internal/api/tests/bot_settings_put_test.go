package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"my-chatbot-backend/internal/api/handlers"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
	"my-chatbot-backend/internal/services"
)

func TestUpdateBotSettings(t *testing.T) {
	setupBotSettingsTestDB(t)
	r := setupBotSettingsRouter()

	orgID := uuid.New()
	token := createBotSettingsToken(orgID.String())

	t.Run("Update Bot Settings - Success (Create New)", func(t *testing.T) {
		payload := services.BotSettingsUpdate{
			BotName:        "New Bot",
			WelcomeMessage: "Hello!",
			SystemPrompt:   "Be helpful.",
			AIProvider:     "openai",
			ModelName:      "gpt-4o-mini",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/settings/bot", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response handlers.FullBotSettingsResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, payload.BotName, response.BotName)
		assert.Equal(t, payload.AIProvider, response.AIProvider)
		assert.Equal(t, payload.ModelName, response.ModelName)

		// Verify in DB
		var settings models.BotSetting
		db.DB.Where("organization_id = ?", orgID).First(&settings)
		assert.Equal(t, payload.BotName, settings.BotName)
		assert.Equal(t, payload.AIProvider, settings.AIProvider)
		assert.Equal(t, payload.ModelName, settings.ModelName)
	})

	t.Run("Update Bot Settings - Success (Update Existing)", func(t *testing.T) {
		payload := services.BotSettingsUpdate{
			BotName:        "Updated Bot",
			WelcomeMessage: "Welcome back!",
			SystemPrompt:   "Be more helpful.",
			AIProvider:     "gemini",
			ModelName:      "gemini-2.0-pro",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/settings/bot", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response handlers.FullBotSettingsResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, payload.BotName, response.BotName)
		assert.Equal(t, payload.AIProvider, response.AIProvider)
		assert.Equal(t, payload.ModelName, response.ModelName)

		// Verify in DB
		var settings models.BotSetting
		db.DB.Where("organization_id = ?", orgID).First(&settings)
		assert.Equal(t, payload.BotName, settings.BotName)
		assert.Equal(t, payload.AIProvider, settings.AIProvider)
		assert.Equal(t, payload.ModelName, settings.ModelName)

		getReq, _ := http.NewRequest(http.MethodGet, "/api/v1/settings/bot", nil)
		getReq.Header.Set("Authorization", "Bearer "+token)
		getReq.Header.Set("X-Org-ID", orgID.String())

		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusOK, getW.Code)
		var getResponse handlers.FullBotSettingsResponse
		json.Unmarshal(getW.Body.Bytes(), &getResponse)
		assert.Equal(t, payload.AIProvider, getResponse.AIProvider)
		assert.Equal(t, payload.ModelName, getResponse.ModelName)
	})

	t.Run("Update Bot Settings - Validation Failure (400)", func(t *testing.T) {
		payload := services.BotSettingsUpdate{
			BotName: "", // Required
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/settings/bot", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("Update Bot Settings - Org Mismatch (403)", func(t *testing.T) {
		otherOrgID := uuid.New().String()
		payload := services.BotSettingsUpdate{
			BotName:        "Bot",
			WelcomeMessage: "Hi",
			SystemPrompt:   "Prompt",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/settings/bot", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", otherOrgID)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "org_mismatch")
	})

	t.Run("Update Bot Settings - Unauthorized (401)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/settings/bot", nil)
		// No token

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
