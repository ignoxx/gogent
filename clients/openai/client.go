package openai

import (
	"context"
	"strings"

	"github.com/ignoxx/gogent/clients"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	*openai.Client
	UserPrompt   string
	SystemPrompt string
	Temperature  float32
}

func NewOpenAIClient(apiKey string) clients.AIClient {
	client := openai.NewClient(apiKey)

	return &Client{
		Client: client,
	}
}

func (c *Client) WithTemperature(temp float32) clients.AIClient {
	c.Temperature = temp
	return c
}

func (c *Client) WithSystemPrompt(prompt string) clients.AIClient {
	c.SystemPrompt = prompt
	return c
}

func (c *Client) WithUserPrompt(prompt string) clients.AIClient {
	c.UserPrompt = prompt
	return c
}

func (c *Client) Run(ctx context.Context, model string) (clients.AIClientResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: c.SystemPrompt},
		{Role: openai.ChatMessageRoleUser, Content: c.UserPrompt},
	}

	resp, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Temperature: c.Temperature,
		Messages:    messages,
	})

	if err != nil {
		return clients.AIClientResponse{}, err
	}

	return parse(resp), err
}

func parse(resp openai.ChatCompletionResponse) clients.AIClientResponse {
	var sb strings.Builder
	for _, message := range resp.Choices {
		sb.WriteString(message.Message.Content)
	}

	return clients.AIClientResponse{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
		Content:          sb.String(),
	}
}
