package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthRoute(t *testing.T) {
	router := setupRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "up", response["status"])
}

func TestCORSAllowsWidgetRuntimeOrigins(t *testing.T) {
	t.Setenv("WIDGET_RUNTIME_ORIGINS", "https://widget.example.com, https://widget-staging.example.com")
	router := setupRouter(nil)

	tests := []struct {
		name       string
		origin     string
		allowValue string
	}{
		{name: "localhost 5173", origin: "http://localhost:5173", allowValue: "http://localhost:5173"},
		{name: "localhost 3000", origin: "http://localhost:3000", allowValue: "http://localhost:3000"},
		{name: "localhost 4173", origin: "http://localhost:4173", allowValue: "http://localhost:4173"},
		{name: "configured widget origin", origin: "https://widget.example.com", allowValue: "https://widget.example.com"},
		{name: "customer site not allowed", origin: "https://customer.example.com", allowValue: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)
			req.Header.Set("Origin", tc.origin)
			router.ServeHTTP(w, req)

			if tc.allowValue == "" {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
			assert.Equal(t, tc.allowValue, w.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}
