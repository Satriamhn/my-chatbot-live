package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"my-chatbot-backend/internal/api"
	"my-chatbot-backend/internal/api/routes"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/middleware"
	"my-chatbot-backend/internal/services"
)

var defaultWidgetDevOrigins = []string{
	"http://localhost:5173",
	"http://localhost:3000",
	"http://localhost:4173",
}

// widgetRuntimeOrigins returns the browser origins allowed to call /api/v1/widget/* directly.
// Keep localhost dev origins enabled, and add production widget app origins via WIDGET_RUNTIME_ORIGINS.
func widgetRuntimeOrigins() []string {
	origins := make([]string, 0, len(defaultWidgetDevOrigins))
	seen := make(map[string]struct{}, len(defaultWidgetDevOrigins))
	add := func(origin string) {
		if origin == "" {
			return
		}
		if _, ok := seen[origin]; ok {
			return
		}
		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	for _, origin := range defaultWidgetDevOrigins {
		add(origin)
	}

	for origin := range strings.SplitSeq(os.Getenv("WIDGET_RUNTIME_ORIGINS"), ",") {
		add(strings.TrimSpace(origin))
	}

	return origins
}

func setupRouter(aiSvc services.AIService) *gin.Engine {
	r := gin.Default()

	// CORS is intentionally limited to the browser origins that directly call /api/v1/widget/*.
	// Customer sites embedding the script are not added here; the widget app origin is.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     widgetRuntimeOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Org-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	kbHandler := &api.KnowledgeBaseHandler{}
	v1 := r.Group("/api/v1")

	// Public routes (no auth)
	routes.RegisterAuthRoutes(v1)
	routes.RegisterWidgetRoutes(v1, aiSvc)

	// Protected routes (JWT required)
	protected := v1.Group("")
	protected.Use(middleware.TenantContract())
	{
		routes.RegisterBotSettingsRoutes(protected)
		routes.RegisterChatRoutes(protected, aiSvc)
		routes.RegisterMeRoute(protected)
		routes.RegisterStatsRoutes(protected)

		kb := protected.Group("/knowledge-base")
		{
			kb.POST("/items", kbHandler.CreateItem)
			kb.GET("/items", kbHandler.ListItems)
			kb.PATCH("/items/:id/status", kbHandler.UpdateItemStatus)
		}
	}

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if _, err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected and migrated")

	// Init AI service (optional — graceful if key missing)
	ctx := context.Background()
	aiSvc, err := services.NewAIService(ctx)
	if err != nil {
		log.Printf("AI service unavailable (set GEMINI_API_KEY to enable): %v", err)
		aiSvc = nil
	} else {
		log.Println("Gemini AI service initialized")
	}

	r := setupRouter(aiSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
