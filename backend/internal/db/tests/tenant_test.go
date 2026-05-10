package tests

import (
	"os"
	"testing"

	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTenantIsolation(t *testing.T) {
	// Skip if no DATABASE_URL is provided (or use a test-specific one)
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("Skipping DB integration test: DATABASE_URL not set")
	}

	database, err := db.InitDB()
	assert.NoError(t, err)

	// Create two organizations
	org1 := models.Organization{Name: "Org 1"}
	org2 := models.Organization{Name: "Org 2"}
	database.Create(&org1)
	database.Create(&org2)

	// Create data for Org 1
	user1 := models.User{
		OrganizationID: org1.ID,
		Email:          "user1@org1.com",
		PasswordHash:   "hash1",
	}
	database.Create(&user1)

	// Create data for Org 2
	user2 := models.User{
		OrganizationID: org2.ID,
		Email:          "user2@org2.com",
		PasswordHash:   "hash2",
	}
	database.Create(&user2)

	// Test cross-tenant isolation with Scope
	t.Run("Org 1 should only see its own users", func(t *testing.T) {
		var users []models.User
		database.Scopes(db.TenantScope(org1.ID.String())).Find(&users)

		assert.Equal(t, 1, len(users))
		assert.Equal(t, user1.ID, users[0].ID)
	})

	t.Run("Org 2 should only see its own users", func(t *testing.T) {
		var users []models.User
		database.Scopes(db.TenantScope(org2.ID.String())).Find(&users)

		assert.Equal(t, 1, len(users))
		assert.Equal(t, user2.ID, users[0].ID)
	})

	t.Run("RLS should block cross-tenant read even if we try to bypass via raw SQL", func(t *testing.T) {
		// Set session to Org 1
		database.Exec("SET LOCAL app.current_tenant = ?", org1.ID.String())

		var count int64
		// Try to count all users in Org 2 while session is Org 1
		database.Model(&models.User{}).Where("organization_id = ?", org2.ID).Count(&count)

		// RLS should filter this out
		assert.Equal(t, int64(0), count)
	})
}
