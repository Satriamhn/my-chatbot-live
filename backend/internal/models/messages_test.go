package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessageModel(t *testing.T) {
	t.Run("CreationAndFields", func(t *testing.T) {
		orgID := uuid.New()
		sessionID := uuid.New()
		content := "Hello, how can I help you?"

		msg := Message{
			OrganizationID: orgID,
			ChatSessionID:  sessionID,
			Sender:         SenderBot,
			Content:        content,
		}

		assert.Equal(t, orgID, msg.OrganizationID)
		assert.Equal(t, sessionID, msg.ChatSessionID)
		assert.Equal(t, SenderBot, msg.Sender)
		assert.Equal(t, content, msg.Content)
	})

	t.Run("SenderEnums", func(t *testing.T) {
		assert.Equal(t, "user", SenderUser)
		assert.Equal(t, "bot", SenderBot)
		assert.Equal(t, "human", SenderHuman)
	})
}
