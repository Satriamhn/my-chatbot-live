package services

import (
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ChatServiceTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *ChatServiceTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to connect database: %v", err)
	}

	db.DB = s.db

	// Custom migration for SQLite to avoid gen_random_uuid() and GORM AutoMigrate issues
	s.db.Exec("CREATE TABLE organizations (id uuid PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime, name text)")
	s.db.Exec("CREATE TABLE users (id uuid PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime, organization_id uuid, email text, password_hash text, role text)")
	s.db.Exec("CREATE TABLE chat_sessions (id uuid PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime, organization_id uuid, user_id uuid, title text, status text)")
	s.db.Exec("CREATE TABLE messages (id uuid PRIMARY KEY, created_at datetime, updated_at datetime, deleted_at datetime, organization_id uuid, chat_session_id uuid, sender text, content text)")
}

func (s *ChatServiceTestSuite) TearDownTest() {
	s.db.Exec("DELETE FROM messages")
	s.db.Exec("DELETE FROM chat_sessions")
}

func TestChatServiceSuite(t *testing.T) {
	suite.Run(t, new(ChatServiceTestSuite))
}

func (s *ChatServiceTestSuite) TestCreateSession() {
	orgID := uuid.New()
	userID := uuid.New()
	payload := CreateSessionPayload{
		Title:  "Test Chat",
		UserID: &userID,
	}

	session, err := CreateSession(orgID, payload)
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, session.ID)
	assert.Equal(s.T(), orgID, session.OrganizationID)
	assert.Equal(s.T(), &userID, session.UserID)
	assert.Equal(s.T(), "Test Chat", session.Title)
	assert.Equal(s.T(), models.SessionStatusBotHandling, session.Status)
}

func (s *ChatServiceTestSuite) TestGetSession() {
	orgID := uuid.New()
	userID := uuid.New()
	session, _ := CreateSession(orgID, CreateSessionPayload{Title: "Test", UserID: &userID})

	t := s.T()
	t.Run("Success", func(t *testing.T) {
		fetched, err := GetSession(orgID, session.ID)
		assert.NoError(t, err)
		assert.Equal(t, session.ID, fetched.ID)
	})

	t.Run("OrgMismatch", func(t *testing.T) {
		wrongOrgID := uuid.New()
		_, err := GetSession(wrongOrgID, session.ID)
		assert.Error(t, err)
		assert.Equal(t, "org_mismatch", err.Error())
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := GetSession(orgID, uuid.New())
		assert.Error(t, err)
		assert.NotEqual(t, "org_mismatch", err.Error())
	})
}

func (s *ChatServiceTestSuite) TestAddMessage() {
	orgID := uuid.New()
	userID := uuid.New()
	session, _ := CreateSession(orgID, CreateSessionPayload{Title: "Test", UserID: &userID})

	payload := CreateMessagePayload{
		Sender:  models.SenderUser,
		Content: "Hello Bot",
	}

	message, err := AddMessage(orgID, session.ID, payload)
	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, message.ID)
	assert.Equal(s.T(), orgID, message.OrganizationID)
	assert.Equal(s.T(), session.ID, message.ChatSessionID)
	assert.Equal(s.T(), models.SenderUser, message.Sender)
	assert.Equal(s.T(), "Hello Bot", message.Content)
}

func (s *ChatServiceTestSuite) TestTakeoverSession() {
	orgID := uuid.New()
	userID := uuid.New()
	session, _ := CreateSession(orgID, CreateSessionPayload{Title: "Test", UserID: &userID})

	t := s.T()
	t.Run("ValidTransition", func(t *testing.T) {
		err := TakeoverSession(orgID, session.ID)
		assert.NoError(t, err)

		fetched, _ := GetSession(orgID, session.ID)
		assert.Equal(t, models.SessionStatusHumanAssigned, fetched.Status)
	})

	t.Run("InvalidTransition", func(t *testing.T) {
		// Already in HumanAssigned, try to takeover again
		_ = TakeoverSession(orgID, session.ID)
		err := TakeoverSession(orgID, session.ID)
		assert.Error(t, err)
		assert.Equal(t, "invalid_transition", err.Error())
	})

	t.Run("OrgMismatch", func(t *testing.T) {
		userID2 := uuid.New()
		session2, _ := CreateSession(orgID, CreateSessionPayload{Title: "Test 2", UserID: &userID2})
		err := TakeoverSession(uuid.New(), session2.ID)
		assert.Error(t, err)
		assert.Equal(t, "org_mismatch", err.Error())
	})
}

func (s *ChatServiceTestSuite) TestReturnToBotMode() {
	orgID := uuid.New()
	userID := uuid.New()
	session, _ := CreateSession(orgID, CreateSessionPayload{Title: "Test", UserID: &userID})

	// First takeover to put it in human_assigned
	_ = TakeoverSession(orgID, session.ID)

	t := s.T()
	t.Run("ValidTransition", func(t *testing.T) {
		err := ReturnToBotMode(orgID, session.ID)
		assert.NoError(t, err)

		fetched, _ := GetSession(orgID, session.ID)
		assert.Equal(t, models.SessionStatusBotHandling, fetched.Status)
	})

	t.Run("InvalidTransition", func(t *testing.T) {
		// Already in BotHandling
		err := ReturnToBotMode(orgID, session.ID)
		assert.Error(t, err)
		assert.Equal(t, "invalid_transition", err.Error())
	})
}
