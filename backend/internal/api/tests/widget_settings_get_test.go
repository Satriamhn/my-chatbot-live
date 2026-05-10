package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"my-chatbot-backend/internal/api/routes"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
)

func setupWidgetSettingsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	routes.RegisterWidgetRoutes(v1, nil)
	return r
}

func TestGetWidgetSettings(t *testing.T) {
	setupBotSettingsTestDB(t)
	r := setupWidgetSettingsRouter()

	orgID := uuid.New()
	db.DB.Create(&models.Organization{BaseModel: models.BaseModel{ID: orgID}, Name: "Test Org"})
	db.DB.Create(&models.BotSetting{
		OrganizationID: orgID,
		BotName:        "Test Bot",
		WelcomeMessage: "Welcome!",
		SystemPrompt:   "Be helpful.",
		AIProvider:     "gemini",
		ModelName:      "gemini-2.0-flash",
	})

	t.Run("Get Widget Settings - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/widget/settings?org_id="+orgID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"bot_name":"Test Bot","welcome_message":"Welcome!"}`, w.Body.String())
		assert.NotContains(t, w.Body.String(), "system_prompt")
		assert.NotContains(t, w.Body.String(), "api_key")
		assert.NotContains(t, w.Body.String(), "has_byok_key")
		assert.NotContains(t, w.Body.String(), "daily_message_limit")
		assert.NotContains(t, w.Body.String(), "daily_message_count")
	})

	t.Run("Get Widget Settings - Missing Org ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/widget/settings", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "missing org_id")
	})

	t.Run("Get Widget Settings - Not Found", func(t *testing.T) {
		missingOrgID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/widget/settings?org_id="+missingOrgID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "bot settings not found")
	})
}
