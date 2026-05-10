package models

import (
	"github.com/google/uuid"
)

// ChatSession represents an active or historical conversation.
type ChatSession struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_org_status;index" json:"organization_id"`
	UserID         *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`
	Status         string    `gorm:"type:varchar(50);not null;default:'bot_handling';index:idx_org_status" json:"status"` // "bot_handling", "human_assigned", "closed"

	// Relationships
	Organization Organization `json:"organization"`
	User         User         `json:"user"`
	Messages     []Message    `json:"messages,omitempty"`
}

const (
	SessionStatusBotHandling   = "bot_handling"
	SessionStatusHumanAssigned = "human_assigned"
	SessionStatusClosed        = "closed"
)

func IsValidSessionTransition(oldStatus, newStatus string) bool {
	switch oldStatus {
	case SessionStatusBotHandling:
		return newStatus == SessionStatusHumanAssigned || newStatus == SessionStatusClosed
	case SessionStatusHumanAssigned:
		return newStatus == SessionStatusBotHandling || newStatus == SessionStatusClosed
	case SessionStatusClosed:
		return false // Once closed, it stays closed or needs a new session
	default:
		return false
	}
}
