package main

import (
	// "bufio"
	"client/chatgpt"
	"client/xclient"
	"context"
	"fmt"
	"log"

	// "os"
	// "strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

func openAIClient(xc *xclient.XClient, ctx context.Context, typ, serviceMethod string, args *chatgpt.ChatGPTArgs, idx int) {
	var reply openai.ChatCompletionResponse
	var err error
	switch typ {
	case "call":
		err = xc.Call(ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%s %s error: %v", typ, serviceMethod, err)
	} else {
		log.Printf("G%d %s %s success: \nQ:%s\nA:%s \n", idx, typ, serviceMethod, args.Request.Messages[0].Content, reply.Choices[0].Message.Content)
	}
}

func call(registry string) {
	d := xclient.NewMRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RoundRobinSelect, nil)
	defer func() { _ = xc.Close() }()
	// send request & receive response
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := "2 + 2 = ?"
			d := openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: s,
					},
				},
			}
			openAIClient(xc, context.Background(), "call", "OpenAIClient.RequestChatGPT", &chatgpt.ChatGPTArgs{Request: d}, i)
		}(i)
	}
	wg.Wait()
}

func console(registry string) {
	d := xclient.NewMRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RoundRobinSelect, nil)
	defer func() { _ = xc.Close() }()
	var text string
	for {
		fmt.Print("请输入prompt: ")
		fmt.Scanln(&text)
		if text == "exit" {
			fmt.Println("程序已退出")
			break
		}
		d := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		}
		openAIClient(xc, context.Background(), "call", "OpenAIClient.RequestChatGPT", &chatgpt.ChatGPTArgs{Request: d}, 0)
	}
}

func main() {
	log.SetFlags(0)
	registryAddr := "http://localhost:9999/_mrpc_/registry"
	// call(registryAddr)
	console(registryAddr)

}
