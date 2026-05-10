package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/middleware"
	"my-chatbot-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}
	newDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	// Clean up and migrate
	newDB.Exec("DROP TABLE IF EXISTS messages")
	newDB.Exec("DROP TABLE IF EXISTS chat_sessions")
	newDB.Exec("DROP TABLE IF EXISTS organizations")
	newDB.Exec("DROP TABLE IF EXISTS users")

	newDB.AutoMigrate(&models.Organization{}, &models.User{}, &models.ChatSession{}, &models.Message{})
	db.DB = newDB
	return newDB
}

func createTestToken(orgID, sub string) string {
	claims := jwt.MapClaims{
		"org_id": orgID,
		"sub":    sub,
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middleware.JWTSecret)
	return tokenString
}

func TestChatSessionAPI(t *testing.T) {
	testDB := setupTestDB(t)
	gin.SetMode(gin.TestMode)

	orgID1 := uuid.New()
	orgID2 := uuid.New()
	userID1 := uuid.New()

	// Create test orgs
	testDB.Create(&models.Organization{BaseModel: models.BaseModel{ID: orgID1}, Name: "Org 1"})
	testDB.Create(&models.Organization{BaseModel: models.BaseModel{ID: orgID2}, Name: "Org 2"})

	h := &ChatSessionHandler{}
	r := gin.New()
	r.Use(middleware.TenantContract())
	r.POST("/sessions", h.CreateSession)
	r.POST("/sessions/:id/takeover", h.TakeoverSession)
	r.POST("/sessions/:id/messages", h.AddMessage)

	t.Run("Create Session", func(t *testing.T) {
		token := createTestToken(orgID1.String(), userID1.String())
		body := bytes.NewBufferString(`{"title": "Test Session"}`)
		req, _ := http.NewRequest("POST", "/sessions", body)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID1.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var session models.ChatSession
		json.Unmarshal(w.Body.Bytes(), &session)
		assert.Equal(t, "Test Session", session.Title)
		assert.Equal(t, orgID1, session.OrganizationID)
	})

	t.Run("Takeover and Block Bot", func(t *testing.T) {
		// 1. Create a session
		session := models.ChatSession{
			OrganizationID: orgID1,
			UserID:         &userID1,
			Title:          "Takeover Test",
			Status:         models.SessionStatusBotHandling,
		}
		testDB.Create(&session)

		token := createTestToken(orgID1.String(), userID1.String())

		// 2. Takeover
		req, _ := http.NewRequest("POST", fmt.Sprintf("/sessions/%s/takeover", session.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID1.String())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 3. Try to add bot message (should be blocked)
		msgBody := bytes.NewBufferString(`{"content": "Hello from bot", "sender_type": "bot"}`)
		req, _ = http.NewRequest("POST", fmt.Sprintf("/sessions/%s/messages", session.ID), msgBody)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID1.String())
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "human takeover active")
	})

	t.Run("Unauthorized Tenant Access", func(t *testing.T) {
		// Session belongs to Org 1
		session := models.ChatSession{
			OrganizationID: orgID1,
			UserID:         &userID1,
			Title:          "Org 1 Session",
		}
		testDB.Create(&session)

		// Try to access with Org 2 token
		token2 := createTestToken(orgID2.String(), userID1.String())
		req, _ := http.NewRequest("POST", fmt.Sprintf("/sessions/%s/takeover", session.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token2)
		req.Header.Set("X-Org-ID", orgID2.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Should be 404 because TenantScope filters it out
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid Transition (Closed)", func(t *testing.T) {
		session := models.ChatSession{
			OrganizationID: orgID1,
			UserID:         &userID1,
			Title:          "Closed Session",
			Status:         models.SessionStatusClosed,
		}
		testDB.Create(&session)

		token := createTestToken(orgID1.String(), userID1.String())
		req, _ := http.NewRequest("POST", fmt.Sprintf("/sessions/%s/takeover", session.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID1.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "invalid state transition")
	})
}
