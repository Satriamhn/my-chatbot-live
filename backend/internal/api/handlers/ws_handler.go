package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
	"my-chatbot-backend/internal/services"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production: check allowed origins
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	conn     *websocket.Conn
	send     chan models.WSMessage
	tenantID uuid.UUID
	userID   uuid.UUID
}

func (c *Client) readPump(aiSvc services.AIService) {
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Printf("[WS] close connection error: %v", err)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("[WS] set read deadline error: %v", err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("[WS] pong deadline error: %v", err)
			return err
		}
		return nil
	})

	for {
		var msg models.WSMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS unexpected close: %v", err)
			}
			break
		}

		log.Printf("[WS] received type=%s session_id=%s", msg.Type, msg.SessionID)
		if msg.Type == "chat_message" {
			go c.handleChatMessage(msg, aiSvc)
		}
	}
}

func (c *Client) handleChatMessage(msg models.WSMessage, _ services.AIService) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("[WS] invalid payload type: %T", msg.Payload)
		return
	}
	content, _ := payload["content"].(string)
	if content == "" {
		log.Printf("[WS] empty content")
		return
	}

	sessionIDStr, _ := msg.SessionID.MarshalText()
	sessionID := string(sessionIDStr)
	log.Printf("[WS] handleChatMessage: session=%s content=%q", sessionID, content)

	// 1. Load session
	var session models.ChatSession
	if err := db.DB.Scopes(db.TenantScope(c.tenantID.String())).
		First(&session, "id = ?", sessionID).Error; err != nil {
		log.Printf("[WS] session not found: %v", err)
		c.send <- models.WSMessage{Type: "error", Payload: gin.H{"message": "session not found"}}
		return
	}

	// 2. Save user message
	userMsg := models.Message{
		OrganizationID: c.tenantID,
		ChatSessionID:  session.ID,
		Sender:         models.SenderUser,
		Content:        content,
	}
	db.DB.Create(&userMsg)

	// 3. Skip AI if human has taken over
	if session.Status == models.SessionStatusHumanAssigned {
		return
	}

	// 4. Load bot settings (provider config + rate limit)
	var botSetting models.BotSetting
	if err := db.DB.Scopes(db.TenantScope(c.tenantID.String())).First(&botSetting).Error; err != nil {
		log.Printf("[WS] bot_settings not found, using defaults")
		botSetting = models.BotSetting{
			SystemPrompt:      "You are a helpful customer support assistant.",
			AIProvider:        models.AIProviderGemini,
			ModelName:         "gemini-2.0-flash",
			DailyMessageLimit: models.DefaultDailyLimit,
		}
	}
	log.Printf("[WS] using provider=%s model=%s byok=%v", botSetting.AIProvider, botSetting.ModelName, botSetting.HasByokKey)

	// 5. Create AI service via factory (BYOK or platform key)
	ctx := context.Background()
	aiSvc, usingByok, err := services.NewAIServiceFromSetting(ctx, &botSetting)
	if err != nil {
		log.Printf("[WS] AI service error: %v", err)
		c.send <- models.WSMessage{
			Type:    "error",
			Payload: gin.H{"message": "AI tidak tersedia. Hubungi admin untuk mengatur API key."},
		}
		return
	}

	// 6. Rate limit check (only for platform key users)
	if !services.CheckAndIncrementRateLimit(&botSetting, usingByok) {
		c.send <- models.WSMessage{
			Type: "error",
			Payload: gin.H{
				"message": "Batas pesan harian tercapai. Upgrade ke Pro atau masukkan API key sendiri di Settings.",
				"code":    "RATE_LIMITED",
			},
		}
		return
	}

	// Increment counter + reset daily if needed
	today := time.Now().Truncate(24 * time.Hour)
	if botSetting.LastResetDate.Before(today) {
		db.DB.Model(&botSetting).Updates(map[string]interface{}{
			"daily_message_count": 1,
			"last_reset_date":     today,
		})
	} else {
		db.DB.Model(&botSetting).UpdateColumn("daily_message_count", gorm.Expr("daily_message_count + 1"))
	}

	// 7. Load conversation history (last 10 messages)
	var history []models.Message
	db.DB.Scopes(db.TenantScope(c.tenantID.String())).
		Where("chat_session_id = ? AND id != ?", session.ID, userMsg.ID).
		Order("created_at ASC").
		Limit(10).
		Find(&history)

	chatHistory := make([]services.ChatTurn, 0, len(history))
	for _, m := range history {
		role := "user"
		if m.Sender == models.SenderBot {
			role = "model"
		}
		chatHistory = append(chatHistory, services.ChatTurn{Role: role, Content: m.Content})
	}

	// 8. Stream AI response
	tokenCh := make(chan string, 64)
	var fullReply strings.Builder

	go func() {
		if err := aiSvc.StreamReply(ctx, botSetting.SystemPrompt, content, chatHistory, tokenCh); err != nil {
			log.Printf("AI stream error: %v", err)
			c.send <- models.WSMessage{
				Type:    "error",
				Payload: gin.H{"message": "AI Stream Error: " + err.Error()},
			}
		}
	}()

	for token := range tokenCh {
		fullReply.WriteString(token)
		log.Printf("[WS] streaming token len=%d", len(token))
		c.send <- models.WSMessage{
			Type:      "token",
			SessionID: session.ID,
			Payload:   models.WSResponseChunk{Content: token, IsLast: false},
		}
	}

	// Signal stream end
	c.send <- models.WSMessage{
		Type:      "token",
		SessionID: session.ID,
		Payload:   models.WSResponseChunk{Content: "", IsLast: true},
	}

	// 9. Save bot reply
	if fullReply.Len() > 0 {
		db.DB.Create(&models.Message{
			OrganizationID: c.tenantID,
			ChatSessionID:  session.ID,
			Sender:         models.SenderBot,
			Content:        fullReply.String(),
		})
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			log.Printf("[WS] close connection error: %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("[WS] set write deadline error: %v", err)
				return
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("[WS] close message write error: %v", err)
				}
				return
			}
			log.Printf("[WS] writePump sending type=%s", message.Type)
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("[WS] writePump write error: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("[WS] ping deadline error: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WS] ping write error: %v", err)
				return
			}
		}
	}
}

// WSHandler handles websocket requests.
type WSHandler struct {
	AISvc services.AIService
}

func (h *WSHandler) ServeWS(c *gin.Context) {
	tenantIDStr, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing from context"})
		return
	}
	tenantID, err := uuid.Parse(tenantIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant_id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	var userID uuid.UUID
	if userIDStr != nil {
		userID, _ = uuid.Parse(userIDStr.(string))
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WS upgrade failed: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan models.WSMessage, 256),
		tenantID: tenantID,
		userID:   userID,
	}

	go client.writePump()
	go client.readPump(h.AISvc)

	client.send <- models.WSMessage{
		Type:    "system",
		Payload: gin.H{"message": "Connected to AI delivery pipeline"},
	}
}
