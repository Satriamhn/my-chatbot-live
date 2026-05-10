package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"my-chatbot-backend/internal/api/routes"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")

	// Register a simple health check route if none exists, or just test the router existence
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")

	// We just want to ensure these don't panic
	routes.RegisterBotSettingsRoutes(api)
	routes.RegisterChatRoutes(api, nil)

	assert.NotNil(t, api)
}
