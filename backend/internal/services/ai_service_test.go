package services

import (
	"context"
	"testing"

	"my-chatbot-backend/internal/models"

	"github.com/stretchr/testify/require"
)

func TestNewAIServiceFromSetting_NormalizesLegacyGeminiModel(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "gemini-test-key")

	svc, usingByok, err := NewAIServiceFromSetting(context.Background(), &models.BotSetting{
		AIProvider: models.AIProviderGemini,
		ModelName:  "gemini-1.5-flash",
	})
	require.NoError(t, err)
	require.False(t, usingByok)

	geminiSvc, ok := svc.(*geminiService)
	require.True(t, ok)
	require.Equal(t, "gemini-2.0-flash", geminiSvc.modelName)
}

func TestNewAIServiceFromSetting_PreservesOpenAIModel(t *testing.T) {
	svc, usingByok, err := NewAIServiceFromSetting(context.Background(), &models.BotSetting{
		AIProvider: models.AIProviderOpenAI,
		ModelName:  "gpt-4o-mini",
		APIKey:     "openai-test-key",
	})
	require.NoError(t, err)
	require.True(t, usingByok)

	openAISvc, ok := svc.(*openAIService)
	require.True(t, ok)
	require.Equal(t, "gpt-4o-mini", openAISvc.modelName)
}
