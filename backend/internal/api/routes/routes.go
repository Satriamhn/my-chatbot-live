package routes

import (
	"github.com/gin-gonic/gin"
	"my-chatbot-backend/internal/api/handlers"
	"my-chatbot-backend/internal/middleware"
	"my-chatbot-backend/internal/services"
)

func RegisterBotSettingsRoutes(r *gin.RouterGroup) {
	h := &handlers.BotSettingsHandler{}

	settings := r.Group("/settings", middleware.TenantContract())
	{
		settings.GET("/bot", h.GetBotSettings)
		settings.PUT("/bot", h.UpdateBotSettings)
	}
}

func RegisterChatRoutes(r *gin.RouterGroup, aiSvc services.AIService) {
	h := &handlers.ChatSessionHandler{}

	sessions := r.Group("/sessions")
	{
		sessions.POST("", h.CreateSession)
		sessions.GET("", h.ListSessions)
		sessions.POST("/:id/takeover", h.TakeoverSession)
		sessions.POST("/:id/return-to-bot", h.ReturnToBotMode)
		sessions.POST("/:id/messages", h.AddMessage)
		sessions.GET("/:id/messages", h.ListMessages)

		wsH := &handlers.WSHandler{AISvc: aiSvc}
		sessions.GET("/ws", wsH.ServeWS)
	}
}

func RegisterMeRoute(r *gin.RouterGroup) {
	h := &handlers.AuthHandler{}
	r.GET("/me", middleware.TenantContract(), h.Me)
}

func RegisterStatsRoutes(r *gin.RouterGroup) {
	h := &handlers.StatsHandler{}
	r.GET("/stats", h.GetDashboardStats)
}

func RegisterWidgetRoutes(r *gin.RouterGroup, aiSvc services.AIService) {
	h := &handlers.WidgetSessionHandler{}
	settingsH := &handlers.BotSettingsHandler{}
	wsH := &handlers.WSHandler{AISvc: aiSvc} // reuse WSHandler

	widget := r.Group("/widget", middleware.WidgetTenantContract())
	{
		widget.GET("/settings", settingsH.GetWidgetSettings)
		widget.POST("/sessions", h.CreateWidgetSession)
		widget.GET("/sessions/ws", wsH.ServeWS)
	}
}
