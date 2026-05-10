package models

import (
	"github.com/google/uuid"
)

// KnowledgeBaseItem represents a single source of truth in a KnowledgeBase.
type KnowledgeBaseItem struct {
	BaseModel
	OrganizationID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"organization_id"`
	KnowledgeBaseID uuid.UUID       `gorm:"type:uuid;not null;index" json:"knowledge_base_id"`
	Type            string          `gorm:"type:varchar(50);not null" json:"type"`                   // "file", "url", "manual"
	Status          string          `gorm:"type:varchar(50);not null;default:'ready'" json:"status"` // "queued", "processing", "ready", "failed"
	Content         string          `gorm:"type:text" json:"content"`
	Embedding       string          `gorm:"type:text" json:"embedding"`
	Metadata        string          `gorm:"type:text" json:"metadata"` // JSON metadata

	// Relationships
	Organization  Organization  `json:"organization"`
	KnowledgeBase KnowledgeBase `json:"knowledge_base"`
}
