package ai

import (
	"context"
	"fmt"

	"github.com/FuradWho/TgRadar-Go/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	cfg    *config.Config
}

func NewClient(cfg *config.Config) *Client {
	aiConfig := openai.DefaultConfig(cfg.AI.APIKey)
	if cfg.AI.BaseURL != "" {
		aiConfig.BaseURL = cfg.AI.BaseURL
	}

	return &Client{
		client: openai.NewClientWithConfig(aiConfig),
		cfg:    cfg,
	}
}

// Analyze performs AI analysis on chat logs
func (c *Client) Analyze(ctx context.Context, chatLog string) (string, error) {
	// Crafted Prompt (Prompt Engineering)
	systemPrompt := `You are a professional community sentiment analyst. Your task is to read a chat log and generate a briefing.
Please output strictly in the following format:
1. **Core Topics**: Summarize the 3 most discussed topics (one sentence each).
2. **Sentiment**: Judge the overall sentiment (Positive/Panic/Wait-and-see/Angry).
3. **High-Value Info**: Extract any specific project names, code error messages, URLs, or breaking news mentioned.
4. **Brief Summary**: Summarize what happened in under 50 words.

Ignore small talk (like greetings, emojis). If the content is too sparse to analyze, please state so directly.`

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.cfg.AI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Here is the recent chat log:\n\n%s", chatLog),
				},
			},
			// Control output length
			MaxTokens: 500,
			// Lower temperature for more objective results
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI response is empty")
	}

	return resp.Choices[0].Message.Content, nil
}
