package chatgpt

import (
	"context"
	"log"
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

func NewOpenAIClient(token string) *OpenAIClient {
	return &OpenAIClient{
		ASAK:   token,
		Client: openai.NewClient(token)}
}

func getOpenAIClient() *OpenAIClient {
	idx := _OpenAISched.r.Intn(_num)
	_OpenAIClient[idx].mu.Lock()
	return _OpenAIClient[idx]
}

func releaseOpenAIClient(c *OpenAIClient) {
	c.mu.Unlock()
}

func (s *OpenAIClient) RequestChatGPT(args ChatGPTArgs, reply *openai.ChatCompletionResponse) error {
	client := getOpenAIClient()
	defer releaseOpenAIClient(client)
	resp, err := client.Client.CreateChatCompletion(
		context.Background(),
		args.Request,
	)
	if err != nil {
		log.Panicln(err)
		return err
	}
	*reply = resp
	return nil
}
