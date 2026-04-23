package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel includes standard GORM fields with UUID primary keys.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

// Organization represents a tenant in the SaaS platform.
type Organization struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	// Relationships
	Users           []User          `json:"users,omitempty"`
	BotSettings     []BotSetting    `json:"bot_settings,omitempty"`
	KnowledgeBases  []KnowledgeBase `json:"knowledge_bases,omitempty"`
	ChatSessions    []ChatSession   `json:"chat_sessions,omitempty"`
}

// User represents an authentication account within an organization.
type User struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	Email          string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash   string    `gorm:"type:varchar(255);not null" json:"-"`
	Role           string    `gorm:"type:varchar(50);not null;default:'user'" json:"role"`

	// Relationships
	Organization Organization  `json:"organization"`
	ChatSessions []ChatSession `json:"chat_sessions,omitempty"`
}

// BotSetting stores configuration for the AI chatbot per organization.
type BotSetting struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	PromptTemplate string    `gorm:"type:text;not null" json:"prompt_template"`

	// Relationships
	Organization Organization `json:"organization"`
}

// KnowledgeBase represents a collection of documents/data for RAG context.
type KnowledgeBase struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`

	// Relationships
	Organization Organization `json:"organization"`
}

// ChatSession represents an active or historical conversation.
type ChatSession struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`

	// Relationships
	Organization Organization `json:"organization"`
	User         User         `json:"user"`
	Messages     []Message    `json:"messages,omitempty"`
}

// Message represents a single message in a ChatSession.
type Message struct {
	BaseModel
	ChatSessionID uuid.UUID `gorm:"type:uuid;not null;index" json:"chat_session_id"`
	SenderType    string    `gorm:"type:varchar(50);not null" json:"sender_type"` // e.g., "user", "bot", "system"
	Content       string    `gorm:"type:text;not null" json:"content"`

	// Relationships
	ChatSession ChatSession `json:"chat_session"`
}
