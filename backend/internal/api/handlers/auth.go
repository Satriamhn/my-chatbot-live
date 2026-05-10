package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"my-chatbot-backend/internal/db"
	"my-chatbot-backend/internal/models"
)

type AuthHandler struct{}

// --- Request/Response Types ---

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	OrgName  string `json:"org_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// Register creates a new Organization and its first admin User.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existing models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create org + user in a transaction
	var org models.Organization
	var user models.User

	txErr := db.DB.Transaction(func(tx *gorm.DB) error {
		org = models.Organization{Name: req.OrgName}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		user = models.User{
			OrganizationID: org.ID,
			Email:          req.Email,
			PasswordHash:   string(hash),
			Role:           "admin",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Default bot settings — auto-detect available AI provider
		defaultProvider := models.AIProviderGemini
		defaultModel := "gemini-2.0-flash"
		if os.Getenv("OPENAI_API_KEY") != "" && os.Getenv("GEMINI_API_KEY") == "" {
			defaultProvider = models.AIProviderOpenAI
			defaultModel = "gpt-4o-mini"
		}

		botSetting := models.BotSetting{
			OrganizationID:    org.ID,
			BotName:           req.OrgName + " Bot",
			WelcomeMessage:    "Hi! How can I help you today?",
			SystemPrompt:      "You are a helpful customer support assistant.",
			AIProvider:        defaultProvider,
			ModelName:         defaultModel,
			DailyMessageLimit: models.DefaultDailyLimit,
		}
		return tx.Create(&botSetting).Error
	})
	if txErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	token, err := generateJWT(user.ID, org.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:  token,
		OrgID:  org.ID.String(),
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
	})
}

// Login authenticates a user and returns a JWT.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := generateJWT(user.ID, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:  token,
		OrgID:  user.OrganizationID.String(),
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
	})
}

// Me returns the current user's info from the JWT.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	var user models.User
	if err := db.DB.Where("id = ? AND organization_id = ?", userID, tenantID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"org_id":  user.OrganizationID,
	})
}

// generateJWT creates a signed JWT with org_id and sub claims.
func generateJWT(userID, orgID uuid.UUID) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	claims := jwt.MapClaims{
		"sub":    userID.String(),
		"org_id": orgID.String(),
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
