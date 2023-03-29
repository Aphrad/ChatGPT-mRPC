package chatgpt

import (
	"sync"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	mu     sync.Mutex
	ASAK   string
	Client *openai.Client
}

type ChatGPTArgs struct {
	Request openai.ChatCompletionRequest
}
