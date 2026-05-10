package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createTestToken(orgID, sub string) string {
	claims := jwt.MapClaims{
		"org_id": orgID,
		"sub":    sub,
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JWTSecret)
	return tokenString
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TenantContract())
	r.GET("/protected", func(c *gin.Context) {
		tenantID, _ := c.Get("tenant_id")
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"tenant_id": tenantID,
			"user_id":   userID,
		})
	})
	return r
}

func TestTenantContractMiddleware(t *testing.T) {
	r := setupTestRouter()

	t.Run("Valid token and matching X-Org-ID", func(t *testing.T) {
		token := createTestToken("org-123", "user-456")
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", "org-123")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"tenant_id":"org-123"`)
		assert.Contains(t, w.Body.String(), `"user_id":"user-456"`)
	})

	t.Run("Valid token but missing X-Org-ID", func(t *testing.T) {
		token := createTestToken("org-123", "user-456")
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Valid token but mismatched X-Org-ID", func(t *testing.T) {
		token := createTestToken("org-123", "user-456")
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Org-ID", "org-999")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "org_mismatch")
	})

	t.Run("Missing authorization header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("X-Org-ID", "org-123")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		req.Header.Set("X-Org-ID", "org-123")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
