package clients

import "context"

type AIClient interface {
	WithSystemPrompt(prompt string) AIClient
	WithUserPrompt(prompt string) AIClient
	WithTemperature(temp float32) AIClient
	Run(ctx context.Context) (AIClientResponse, error)
}

type AIClientResponse struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Content          string
}
