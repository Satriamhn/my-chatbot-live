package models

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestKnowledgeBaseModel(t *testing.T) {
	// Setup in-memory sqlite for testing GORM model
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Minimal mock struct to avoid postgres-specific syntax in sqlite
	type TestBaseModel struct {
		ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type TestOrg struct {
		TestBaseModel
		Name string
	}
	type TestKB struct {
		TestBaseModel
		OrganizationID uuid.UUID `gorm:"index"`
		Name           string
		SourceType     string
		Source         string
		Status         string
	}

	// Migrate the schema
	err = db.AutoMigrate(&TestOrg{}, &TestKB{})
	require.NoError(t, err)

	orgID := uuid.New()
	org := TestOrg{
		TestBaseModel: TestBaseModel{ID: orgID},
		Name:          "Test Org",
	}
	err = db.Create(&org).Error
	require.NoError(t, err)

	t.Run("Create KnowledgeBase with all fields", func(t *testing.T) {
		kb := TestKB{
			TestBaseModel:  TestBaseModel{ID: uuid.New()},
			OrganizationID: orgID,
			Name:           "Company Docs",
			SourceType:     KBSourceTypeDocument,
			Source:         "/path/to/docs",
			Status:         KBStatusQueued,
		}

		err := db.Create(&kb).Error
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, kb.ID)
	})

	t.Run("Tenant Scoping Verification", func(t *testing.T) {
		orgID2 := uuid.New()
		kb := TestKB{
			TestBaseModel:  TestBaseModel{ID: uuid.New()},
			OrganizationID: orgID2,
			Name:           "Other Org KB",
			SourceType:     KBSourceTypeURL,
			Source:         "https://example.com",
			Status:         KBStatusReady,
		}
		err := db.Create(&kb).Error
		assert.NoError(t, err)

		// Query for first org's KBs
		var kbs []TestKB
		db.Where("organization_id = ?", orgID).Find(&kbs)
		assert.Len(t, kbs, 1)
		assert.Equal(t, "Company Docs", kbs[0].Name)

		// Query for second org's KBs
		db.Where("organization_id = ?", orgID2).Find(&kbs)
		assert.Len(t, kbs, 1)
		assert.Equal(t, "Other Org KB", kbs[0].Name)
	})

	t.Run("Enum Constants Verification", func(t *testing.T) {
		assert.Equal(t, "document", KBSourceTypeDocument)
		assert.Equal(t, "url", KBSourceTypeURL)
		assert.Equal(t, "manual", KBSourceTypeManual)

		assert.Equal(t, "queued", KBStatusQueued)
		assert.Equal(t, "processing", KBStatusProcessing)
		assert.Equal(t, "ready", KBStatusReady)
		assert.Equal(t, "failed", KBStatusFailed)
	})
}
