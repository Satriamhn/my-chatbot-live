package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AI Provider constants
const (
	AIProviderGemini = "gemini"
	AIProviderOpenAI = "openai"

	// Default daily limit when using platform key
	DefaultDailyLimit = 100
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
	Users              []User              `json:"users,omitempty"`
	BotSettings        []BotSetting        `json:"bot_settings,omitempty"`
	KnowledgeBases     []KnowledgeBase     `json:"knowledge_bases,omitempty"`
	KnowledgeBaseItems []KnowledgeBaseItem `json:"knowledge_base_items,omitempty"`
	ChatSessions       []ChatSession       `json:"chat_sessions,omitempty"`
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
	BotName        string    `gorm:"type:varchar(255);not null" json:"bot_name"`
	WelcomeMessage string    `gorm:"type:text;not null" json:"welcome_message"`
	SystemPrompt   string    `gorm:"type:text;not null" json:"system_prompt"`

	// AI Provider config (Hybrid model)
	AIProvider string `gorm:"column:ai_provider;type:varchar(50);not null;default:'gemini'" json:"ai_provider"` // "gemini" | "openai"
	ModelName  string `gorm:"column:model_name;type:varchar(100);not null;default:'gemini-2.0-flash'" json:"model_name"`
	APIKey     string `gorm:"column:api_key;type:text" json:"-"`                              // BYOK — never exposed in API response
	HasByokKey bool   `gorm:"column:has_byok_key;not null;default:false" json:"has_byok_key"` // frontend can check if key is set

	// Rate limiting (used when tenant uses platform key)
	DailyMessageLimit int       `gorm:"column:daily_message_limit;not null;default:100" json:"daily_message_limit"`
	DailyMessageCount int       `gorm:"column:daily_message_count;not null;default:0" json:"daily_message_count"`
	LastResetDate     time.Time `gorm:"column:last_reset_date;not null;default:now()" json:"last_reset_date"`

	// Relationships
	Organization Organization `json:"organization"`
}

// KnowledgeBase represents a collection of documents/data for RAG context.
type KnowledgeBase struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_org_kb" json:"organization_id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	SourceType     string    `gorm:"type:varchar(50);not null" json:"source_type"` // "document", "url", "manual"
	Source         string    `gorm:"type:text" json:"source"`                      // File path, URL, or raw text
	Status         string    `gorm:"type:varchar(50);not null;default:'queued'" json:"status"`

	// Relationships
	Organization       Organization        `json:"organization"`
	KnowledgeBaseItems []KnowledgeBaseItem `json:"knowledge_base_items,omitempty"`
}

const (
	KBSourceTypeDocument = "document"
	KBSourceTypeURL      = "url"
	KBSourceTypeManual   = "manual"

	KBStatusQueued     = "queued"
	KBStatusProcessing = "processing"
	KBStatusReady      = "ready"
	KBStatusFailed     = "failed"
)

// Message represents a single message in a ChatSession.
type Message struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_org_session;index" json:"organization_id"`
	ChatSessionID  uuid.UUID `gorm:"type:uuid;not null;index:idx_org_session;index" json:"chat_session_id"`
	Sender         string    `gorm:"type:varchar(50)" json:"sender"` // "user", "bot", "human"
	Content        string    `gorm:"type:text;not null" json:"content"`

	// Relationships
	Organization Organization `json:"organization"`
	ChatSession  ChatSession  `json:"chat_session"`
}

const (
	SenderUser  = "user"
	SenderBot   = "bot"
	SenderHuman = "human"
)

// WSMessage represents a message sent over the WebSocket pipeline.
type WSMessage struct {
	Type      string      `json:"type"`       // "chat_message", "token", "error", "pong"
	Payload   interface{} `json:"payload"`    // The actual data payload
	SessionID uuid.UUID   `json:"session_id"` // Associated ChatSession
}

// WSResponseChunk represents a piece of a streaming AI response.
type WSResponseChunk struct {
	Content string `json:"content"`
	IsLast  bool   `json:"is_last"`
}
