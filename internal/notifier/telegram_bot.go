package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const telegramMaxMessageLen = 3500

type TelegramBot struct {
	token  string
	chatID int64
	client *http.Client
}

func NewTelegramBot(token string, chatID int64) *TelegramBot {
	return &TelegramBot{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type telegramSendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func (t *TelegramBot) Send(ctx context.Context, text string) error {
	if t == nil || t.token == "" || t.chatID == 0 || text == "" {
		return nil
	}

	chunks := splitText(text, telegramMaxMessageLen)
	for _, chunk := range chunks {
		if err := t.sendOnce(ctx, chunk); err != nil {
			time.Sleep(500 * time.Millisecond)
			if retryErr := t.sendOnce(ctx, chunk); retryErr != nil {
				return retryErr
			}
		}
	}

	return nil
}

func (t *TelegramBot) sendOnce(ctx context.Context, text string) error {
	payload := telegramSendMessageRequest{
		ChatID: t.chatID,
		Text:   text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram bot send failed: %s", resp.Status)
	}

	return nil
}

func splitText(text string, maxLen int) []string {
	if maxLen <= 0 {
		return []string{text}
	}

	runes := []rune(text)
	if len(runes) <= maxLen {
		return []string{text}
	}

	var chunks []string
	for start := 0; start < len(runes); start += maxLen {
		end := start + maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}
