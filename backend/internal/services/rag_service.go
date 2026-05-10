package services

import (
	"context"
	"fmt"
	"my-chatbot-backend/internal/models"
	"strings"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// RAGService defines the operations for Retrieval-Augmented Generation.
type RAGService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	RetrieveContext(ctx context.Context, orgID uuid.UUID, queryEmbedding []float32, limit int) (string, error)
	AssembleSystemPrompt(basePrompt string, retrievedContext string, userMessage string) string
	ProcessRAG(ctx context.Context, orgID uuid.UUID, userMessage string, basePrompt string) (string, error)
}

type ragService struct {
	db *gorm.DB
}

// NewRAGService creates a new instance of RAGService.
func NewRAGService(db *gorm.DB) RAGService {
	return &ragService{db: db}
}

// GenerateEmbedding creates a vector embedding for the input text.
func (s *ragService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Mock implementation for embedding generation
	// In a real scenario, this would call an external API (e.g., OpenAI embeddings)
	// We return a mock vector of 1536 dimensions for pgvector
	mockVector := make([]float32, 1536)
	for i := range mockVector {
		// Just a dummy representation based on string length to simulate variety
		mockVector[i] = float32(len(text)) * 0.001
	}
	return mockVector, nil
}

// RetrieveContext performs a similarity search on KnowledgeBaseItem using pgvector with tenant isolation.
func (s *ragService) RetrieveContext(ctx context.Context, orgID uuid.UUID, queryEmbedding []float32, limit int) (string, error) {
	if orgID == uuid.Nil {
		return "", fmt.Errorf("invalid organization ID")
	}

	var items []models.KnowledgeBaseItem

	// Perform similarity search using cosine distance (<=>)
	// We enforce tenant isolation using organization_id filter.
	err := s.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Where("status = ?", "ready").
		Order(gorm.Expr("embedding <=> ?", pgvector.NewVector(queryEmbedding))).
		Limit(limit).
		Find(&items).Error

	if err != nil {
		return "", fmt.Errorf("failed to retrieve knowledge base items: %w", err)
	}

	if len(items) == 0 {
		return "", nil
	}

	var builder strings.Builder
	for i, item := range items {
		fmt.Fprintf(&builder, "--- Source %d ---\n", i+1)
		builder.WriteString(item.Content)
		builder.WriteString("\n\n")
	}

	return builder.String(), nil
}

// AssembleSystemPrompt assembles the retrieved text into a clean system prompt context.
func (s *ragService) AssembleSystemPrompt(basePrompt string, retrievedContext string, userMessage string) string {
	if retrievedContext == "" {
		return fmt.Sprintf("%s\n\nUser Message:\n%s", basePrompt, userMessage)
	}

	return fmt.Sprintf("%s\n\nContext information is below:\n---------------------\n%s\n---------------------\nGiven the context information and not prior knowledge, answer the query.\n\nUser Message:\n%s", basePrompt, retrievedContext, userMessage)
}

// ProcessRAG orchestrates the entire RAG pipeline: Embedding -> Retrieval -> Prompt Assembly
func (s *ragService) ProcessRAG(ctx context.Context, orgID uuid.UUID, userMessage string, basePrompt string) (string, error) {
	// 1. Generate Embedding
	embedding, err := s.GenerateEmbedding(ctx, userMessage)
	if err != nil {
		return "", fmt.Errorf("embedding generation failed: %w", err)
	}

	// 2. Retrieve Context (using top 3 items)
	retrievedContext, err := s.RetrieveContext(ctx, orgID, embedding, 3)
	if err != nil {
		return "", fmt.Errorf("context retrieval failed: %w", err)
	}

	// 3. Assemble Prompt
	finalPrompt := s.AssembleSystemPrompt(basePrompt, retrievedContext, userMessage)

	return finalPrompt, nil
}
