package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	go_openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client   *go_openai.Client
	ctx      context.Context
	messages []go_openai.ChatCompletionMessage
}

func NewClient(token string) *Client {
	client := go_openai.NewClient(token)
	ctx := context.Background()
	messages := make([]go_openai.ChatCompletionMessage, 0)
	return &Client{client, ctx, messages}
}

func (client *Client) stream(input string) (chan string, error) {
	channel := make(chan string)

	new_request := go_openai.ChatCompletionMessage{
		Role:    go_openai.ChatMessageRoleUser,
		Content: input,
	}
	client.messages = append(client.messages, new_request)

	req := go_openai.ChatCompletionRequest{
		Model:     go_openai.GPT4,
		MaxTokens: 500,
		Messages:  client.messages,
		Stream:    true,
	}
	stream, err := client.client.CreateChatCompletionStream(client.ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return nil, err
	}

	go func() {
		defer stream.Close()
		defer close(channel)

		var builder strings.Builder

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
			}

			content := response.Choices[0].Delta.Content

			channel <- content
			builder.WriteString(content)
		}

		new_response := go_openai.ChatCompletionMessage{
			Role:    go_openai.ChatMessageRoleAssistant,
			Content: builder.String(),
		}
		client.messages = append(client.messages, new_response)
	}()

	return channel, nil
}

func (client *Client) AddPrompt(prompt string) {
	client.messages = append(client.messages, go_openai.ChatCompletionMessage{
		Role:    go_openai.ChatMessageRoleUser,
		Content: prompt,
	})
	client.messages = append(client.messages, go_openai.ChatCompletionMessage{
		Role:    go_openai.ChatMessageRoleAssistant,
		Content: "Yes.",
	})
}
