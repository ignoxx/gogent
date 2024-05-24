package gogent

import (
	"github.com/ignoxx/gogent/clients"
	"github.com/ignoxx/gogent/clients/openai"
	oai "github.com/sashabaranov/go-openai"
)

type Provider int

const (
	ProviderOpenAI Provider = iota + 1
)

type Model string

const (
	ModelOpenAIGpt35 = oai.GPT3Dot5Turbo0613
	ModelOpenAIGpt4o = oai.GPT4o
)

type LLM struct {
	provider Provider
	model    Model
	apiKey   string

	CompletionTokensUsed int
	PromptTokensUsed     int
	TotalTokensUsed      int
}

type LLMOpt func(*LLM) *LLM

func NewLLM(opts ...LLMOpt) *LLM {
	llm := &LLM{}

	for _, opt := range opts {
		llm = opt(llm)
	}

	if llm.apiKey == "" {
		panic("missing api-key!")
	}

	return llm
}

func WithProvider(provider Provider) LLMOpt {
	return func(llm *LLM) *LLM {
		llm.provider = provider
		return llm
	}
}

func WithModel(model Model) LLMOpt {
	return func(llm *LLM) *LLM {
		llm.model = model
		return llm
	}
}

func WithApiKey(apiKey string) LLMOpt {
	return func(llm *LLM) *LLM {
		llm.apiKey = apiKey
		return llm
	}
}

func (llm *LLM) Client() clients.AIClient {
	switch llm.provider {
	case ProviderOpenAI:
		return openai.NewOpenAIClient(llm.apiKey)
	}

	panic("no client matched")
}
