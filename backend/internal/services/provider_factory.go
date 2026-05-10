package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"my-chatbot-backend/internal/models"
)

// ProviderFactory creates the right AIService based on BotSetting.
// Hybrid logic:
//   - If tenant has their own BYOK key → use it (no rate limit)
//   - If no BYOK key → use platform key with rate limit enforcement
func NewAIServiceFromSetting(ctx context.Context, setting *models.BotSetting) (AIService, bool, error) {
	usingByok := setting.APIKey != ""

	// Determine which API key to use
	apiKey := setting.APIKey // BYOK key (may be empty)

	switch setting.AIProvider {
	case models.AIProviderOpenAI:
		// Fallback to platform OpenAI key if no BYOK
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			return nil, false, fmt.Errorf("no OpenAI API key configured for this organization")
		}
		svc, err := NewOpenAIService(apiKey, setting.ModelName)
		return svc, usingByok, err

	case models.AIProviderGemini:
		fallthrough
	default:
		// Fallback to platform Gemini key if no BYOK
		if apiKey == "" {
			apiKey = os.Getenv("GEMINI_API_KEY")
		}
		if apiKey == "" {
			return nil, false, fmt.Errorf("no Gemini API key configured for this organization")
		}
		svc, err := NewGeminiService(ctx, apiKey, setting.ModelName)
		return svc, usingByok, err
	}
}

// CheckAndIncrementRateLimit checks if the tenant has exceeded their daily limit.
// Returns true if the request is allowed, false if rate limited.
// Only enforced when using the platform key (not BYOK).
func CheckAndIncrementRateLimit(setting *models.BotSetting, usingByok bool) bool {
	if usingByok {
		return true // BYOK users have no platform rate limit
	}
	if setting.DailyMessageLimit <= 0 {
		return true // No limit configured
	}
	if setting.DailyMessageCount >= setting.DailyMessageLimit {
		log.Printf("Rate limit exceeded for org %s: %d/%d", setting.OrganizationID, setting.DailyMessageCount, setting.DailyMessageLimit)
		return false
	}
	return true
}

// DefaultModelFor returns the default model name for a given provider.
func DefaultModelFor(provider string) string {
	switch provider {
	case models.AIProviderOpenAI:
		return "gpt-4o-mini"
	default:
		return "gemini-2.0-flash"
	}
}
