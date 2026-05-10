package services

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type openAIService struct {
	client    *openai.Client
	modelName string
}

// NewOpenAIService creates a new OpenAI-backed AIService.
func NewOpenAIService(apiKey, modelName string) (AIService, error) {
	key := apiKey
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if key == "" {
		return nil, fmt.Errorf("no OpenAI API key available")
	}
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	return &openAIService{
		client:    openai.NewClient(key),
		modelName: modelName,
	}, nil
}

// StreamReply sends a streaming request to OpenAI and writes tokens to the out channel.
func (s *openAIService) StreamReply(
	ctx context.Context,
	systemPrompt, userMessage string,
	history []ChatTurn,
	out chan<- string,
) error {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
	}

	for _, turn := range history {
		role := openai.ChatMessageRoleUser
		if turn.Role == "model" || turn.Role == "bot" {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: turn.Content,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	req := openai.ChatCompletionRequest{
		Model:    s.modelName,
		Messages: messages,
		Stream:   true,
	}

	stream, err := s.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("openai stream error: %w", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			// io.EOF means stream ended normally
			break
		}
		if len(response.Choices) > 0 {
			token := response.Choices[0].Delta.Content
			if token != "" {
				select {
				case out <- token:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
	close(out)
	return nil
}
