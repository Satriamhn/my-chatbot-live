package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

// AIService handles streaming AI responses.
type AIService interface {
	StreamReply(ctx context.Context, systemPrompt, userMessage string, history []ChatTurn, out chan<- string) error
}

// ChatTurn represents a single turn in the conversation history.
type ChatTurn struct {
	Role    string // "user" or "model"
	Content string
}

type geminiService struct {
	client    *genai.Client
	modelName string
}

// NewAIService creates a platform-level Gemini service (uses GEMINI_API_KEY env).
// Used during server startup for the default platform service.
func NewAIService(ctx context.Context) (AIService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	return NewGeminiService(ctx, apiKey, "gemini-2.0-flash")
}

// NewGeminiService creates a Gemini AIService with a specific key and model.
func NewGeminiService(ctx context.Context, apiKey, modelName string) (AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}
	modelName = normalizeGeminiModelName(modelName)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &geminiService{client: client, modelName: modelName}, nil
}

func normalizeGeminiModelName(modelName string) string {
	switch strings.TrimSpace(modelName) {
	case "", "gemini-1.5-flash":
		return "gemini-2.0-flash"
	default:
		return modelName
	}
}

// StreamReply sends a streaming request to Gemini and writes tokens to the out channel.
func (s *geminiService) StreamReply(
	ctx context.Context,
	systemPrompt, userMessage string,
	history []ChatTurn,
	out chan<- string,
) error {
	defer close(out)

	var contents []*genai.Content
	for _, turn := range history {
		role := turn.Role
		if role == "bot" {
			role = "model"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: turn.Content}},
		})
	}
	contents = append(contents, &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: userMessage}},
	})

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}

	stream := s.client.Models.GenerateContentStream(ctx, s.modelName, contents, config)
	for response, err := range stream {
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}
		for _, part := range response.Candidates[0].Content.Parts {
			if part.Text != "" {
				select {
				case out <- part.Text:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
	return nil
}
