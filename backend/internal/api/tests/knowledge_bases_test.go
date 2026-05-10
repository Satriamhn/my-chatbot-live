package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	"my-chatbot-backend/internal/api"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/middleware"
	"my-chatbot-backend/internal/models"
)

func setupTestDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}
	testDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	testDB.Exec("DROP TABLE IF EXISTS knowledge_base_items")
	testDB.Exec("DROP TABLE IF EXISTS knowledge_bases")
	testDB.Exec("DROP TABLE IF EXISTS organizations")

	testDB.AutoMigrate(
		&models.Organization{},
		&models.KnowledgeBase{},
		&models.KnowledgeBaseItem{},
	)
	db.DB = testDB
}

func createTestToken(orgID string) string {
	claims := jwt.MapClaims{
		"org_id": orgID,
		"sub":    uuid.New().String(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(middleware.JWTSecret)
	return tokenString
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	kbHandler := &api.KnowledgeBaseHandler{}

	v1 := r.Group("/api/v1")
	v1.Use(middleware.TenantContract())
	{
		kb := v1.Group("/knowledge-base")
		{
			kb.POST("/items", kbHandler.CreateItem)
			kb.GET("/items", kbHandler.ListItems)
			kb.PATCH("/items/:id/status", kbHandler.UpdateItemStatus)
		}
	}
	return r
}

func TestKnowledgeBaseAPI(t *testing.T) {
	setupTestDB(t)
	r := setupTestRouter()

	orgID := uuid.New()
	token := createTestToken(orgID.String())
	kbID := uuid.New()

	// Create a dummy KB in DB first (optional for this test but good practice)
	db.DB.Create(&models.KnowledgeBase{
		BaseModel:      models.BaseModel{ID: kbID},
		OrganizationID: orgID,
		Name:           "Test KB",
	})

	t.Run("Create Item - Success", func(t *testing.T) {
		payload := map[string]interface{}{
			"knowledge_base_id": kbID,
			"type":              "url",
			"content":           "https://example.com",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/knowledge-base/items", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.KnowledgeBaseItem
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "url", response.Type)
		assert.Equal(t, "queued", response.Status)
	})

	t.Run("Create Item - Invalid Type (422)", func(t *testing.T) {
		payload := map[string]interface{}{
			"knowledge_base_id": kbID,
			"type":              "invalid_type",
			"content":           "some content",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/knowledge-base/items", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Update Status - Success", func(t *testing.T) {
		item := models.KnowledgeBaseItem{
			BaseModel:       models.BaseModel{ID: uuid.New()},
			OrganizationID:  orgID,
			KnowledgeBaseID: kbID,
			Type:            "manual",
			Status:          "queued",
		}
		db.DB.Create(&item)

		payload := map[string]interface{}{"status": "ready"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/knowledge-base/items/%s/status", item.ID), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.KnowledgeBaseItem
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "ready", response.Status)
	})

	t.Run("Update Status - Invalid Lifecycle (422)", func(t *testing.T) {
		item := models.KnowledgeBaseItem{
			BaseModel:       models.BaseModel{ID: uuid.New()},
			OrganizationID:  orgID,
			KnowledgeBaseID: kbID,
			Type:            "manual",
			Status:          "queued",
		}
		db.DB.Create(&item)

		payload := map[string]interface{}{"status": "completed"} // "ready" is valid, "completed" is not
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/knowledge-base/items/%s/status", item.ID), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", orgID.String())

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
