package routes

import (
	"github.com/gin-gonic/gin"
	"my-chatbot-backend/internal/api/handlers"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	h := &handlers.AuthHandler{}

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}
