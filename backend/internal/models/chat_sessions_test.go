package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChatSessionModel(t *testing.T) {
	t.Run("CreationAndDefaults", func(t *testing.T) {
		orgID := uuid.New()
		userID := uuid.New()

		cs := ChatSession{
			OrganizationID: orgID,
			UserID:         &userID,
			Title:          "Test Session",
			Status:         SessionStatusBotHandling,
		}

		assert.Equal(t, orgID, cs.OrganizationID)
		assert.Equal(t, &userID, cs.UserID)
		assert.Equal(t, "Test Session", cs.Title)
		assert.Equal(t, SessionStatusBotHandling, cs.Status)
	})

	t.Run("Transitions", func(t *testing.T) {
		assert.True(t, IsValidSessionTransition(SessionStatusBotHandling, SessionStatusHumanAssigned))
		assert.True(t, IsValidSessionTransition(SessionStatusBotHandling, SessionStatusClosed))
		assert.True(t, IsValidSessionTransition(SessionStatusHumanAssigned, SessionStatusBotHandling))
		assert.True(t, IsValidSessionTransition(SessionStatusHumanAssigned, SessionStatusClosed))
		assert.False(t, IsValidSessionTransition(SessionStatusClosed, SessionStatusBotHandling))
	})
}
