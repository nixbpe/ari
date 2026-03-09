package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

var openAISleep = sleepWithContext

func NewOpenAIProvider(cfg Config) *OpenAIProvider {
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}

	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) Complete(ctx context.Context, prompt string, opts ...Option) (string, error) {
	cfg := callConfig{Temperature: 0.1}
	for _, opt := range opts {
		opt(&cfg)
	}

	type requestBody struct {
		Model       string  `json:"model"`
		Messages    []any   `json:"messages"`
		Temperature float64 `json:"temperature"`
	}

	payload, err := json.Marshal(requestBody{
		Model: p.model,
		Messages: []any{
			map[string]string{"role": "user", "content": prompt},
		},
		Temperature: cfg.Temperature,
	})
	if err != nil {
		return "", fmt.Errorf("marshal openai request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(payload))
		if err != nil {
			return "", fmt.Errorf("build openai request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+p.apiKey)

		resp, err := p.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("openai request failed: %w", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return "", fmt.Errorf("read openai response: %w", readErr)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			var parsed struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}
			if err := json.Unmarshal(body, &parsed); err != nil {
				return "", fmt.Errorf("parse openai response: %w", err)
			}
			if len(parsed.Choices) == 0 {
				return "", fmt.Errorf("openai response missing choices")
			}
			return parsed.Choices[0].Message.Content, nil
		case http.StatusUnauthorized:
			return "", fmt.Errorf("OpenAI API key is invalid or missing")
		case http.StatusTooManyRequests:
			lastErr = fmt.Errorf("openai rate limited: %s", strings.TrimSpace(string(body)))
			if attempt == 3 {
				break
			}
			backoff := time.Duration(1<<attempt) * time.Second
			if err := openAISleep(ctx, backoff); err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("openai request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
		}
	}

	if lastErr != nil {
		return "", lastErr
	}

	return "", fmt.Errorf("openai request failed after retries")
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
