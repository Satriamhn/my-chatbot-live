package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"my-chatbot-backend/internal/api/handlers"
	"my-chatbot-backend/internal/api/routes"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/middleware"
	"my-chatbot-backend/internal/models"
)

func setupBotSettingsTestDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	testDB.Exec("DROP TABLE IF EXISTS bot_settings")
	testDB.Exec("DROP TABLE IF EXISTS organizations")

	testDB.AutoMigrate(
		&models.Organization{},
		&models.BotSetting{},
	)
	db.DB = testDB
}

func createBotSettingsToken(orgID string) string {
	claims := jwt.MapClaims{
		"org_id": orgID,
		"sub":    uuid.New().String(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middleware.JWTSecret)
	return tokenString
}

func setupBotSettingsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	routes.RegisterBotSettingsRoutes(v1)
	return r
}

func TestGetBotSettings(t *testing.T) {
	setupBotSettingsTestDB(t)
	r := setupBotSettingsRouter()

	orgID := uuid.New()
	token := createBotSettingsToken(orgID.String())

	// Seed bot settings
	settings := models.BotSetting{
		OrganizationID: orgID,
		BotName:        "Test Bot",
		WelcomeMessage: "Welcome!",
		SystemPrompt:   "Be helpful.",
		AIProvider:     "gemini",
		ModelName:      "gemini-2.0-flash",
	}
	db.DB.Create(&settings)

	t.Run("Get Bot Settings - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/settings/bot", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fullResponse handlers.FullBotSettingsResponse
		json.Unmarshal(w.Body.Bytes(), &fullResponse)
		assert.Equal(t, "Test Bot", fullResponse.BotName)
		assert.Equal(t, "Welcome!", fullResponse.WelcomeMessage)
		assert.Equal(t, "Be helpful.", fullResponse.SystemPrompt)
		assert.Equal(t, "gemini", fullResponse.AIProvider)
		assert.Equal(t, "gemini-2.0-flash", fullResponse.ModelName)
	})

	t.Run("Get Bot Settings - Org Mismatch (403)", func(t *testing.T) {
		otherOrgID := uuid.New().String()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/settings/bot", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", otherOrgID)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "org_mismatch")
	})

	t.Run("Get Bot Settings - Unauthorized (401)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/settings/bot", nil)
		// No token

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
