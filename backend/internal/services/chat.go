package services

import (
	"fmt"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"

	"github.com/google/uuid"
)

type CreateSessionPayload struct {
	Title  string
	UserID *uuid.UUID
}

type CreateMessagePayload struct {
	Sender  string
	Content string
}

func CreateSession(orgID uuid.UUID, payload CreateSessionPayload) (models.ChatSession, error) {
	session := models.ChatSession{
		OrganizationID: orgID,
		UserID:         payload.UserID,
		Title:          payload.Title,
		Status:         models.SessionStatusBotHandling,
	}

	if err := db.DB.Create(&session).Error; err != nil {
		return models.ChatSession{}, err
	}

	return session, nil
}

func GetSession(orgID, sessionID uuid.UUID) (models.ChatSession, error) {
	var session models.ChatSession
	if err := db.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		return models.ChatSession{}, err
	}

	if session.OrganizationID != orgID {
		return models.ChatSession{}, fmt.Errorf("org_mismatch")
	}

	return session, nil
}

func AddMessage(orgID, sessionID uuid.UUID, msg CreateMessagePayload) (models.Message, error) {
	// Validate session exists and belongs to org
	session, err := GetSession(orgID, sessionID)
	if err != nil {
		return models.Message{}, err
	}

	message := models.Message{
		OrganizationID: orgID,
		ChatSessionID:  session.ID,
		Sender:         msg.Sender,
		Content:        msg.Content,
	}

	if err := db.DB.Create(&message).Error; err != nil {
		return models.Message{}, err
	}

	return message, nil
}

func TakeoverSession(orgID, sessionID uuid.UUID) error {
	session, err := GetSession(orgID, sessionID)
	if err != nil {
		return err
	}

	if !models.IsValidSessionTransition(session.Status, models.SessionStatusHumanAssigned) {
		return fmt.Errorf("invalid_transition")
	}

	session.Status = models.SessionStatusHumanAssigned
	if err := db.DB.Save(&session).Error; err != nil {
		return err
	}

	return nil
}

func ReturnToBotMode(orgID, sessionID uuid.UUID) error {
	session, err := GetSession(orgID, sessionID)
	if err != nil {
		return err
	}

	if !models.IsValidSessionTransition(session.Status, models.SessionStatusBotHandling) {
		return fmt.Errorf("invalid_transition")
	}

	session.Status = models.SessionStatusBotHandling
	if err := db.DB.Save(&session).Error; err != nil {
		return err
	}

	return nil
}
