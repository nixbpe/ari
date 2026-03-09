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

type AnthropicProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

func NewAnthropicProvider(cfg Config) *AnthropicProvider {
	model := cfg.Model
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	return &AnthropicProvider{
		apiKey:  cfg.APIKey,
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *AnthropicProvider) Name() string { return "anthropic" }

func (p *AnthropicProvider) Complete(ctx context.Context, prompt string, opts ...Option) (string, error) {
	cfg := callConfig{MaxTokens: 1024}
	for _, opt := range opts {
		opt(&cfg)
	}

	type requestBody struct {
		Model     string `json:"model"`
		MaxTokens int    `json:"max_tokens"`
		Messages  []any  `json:"messages"`
	}

	payload, err := json.Marshal(requestBody{
		Model:     p.model,
		MaxTokens: cfg.MaxTokens,
		Messages: []any{
			map[string]string{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal anthropic request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("build anthropic request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read anthropic response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("parse anthropic response: %w", err)
	}
	if len(parsed.Content) == 0 {
		return "", fmt.Errorf("anthropic response missing content")
	}

	return parsed.Content[0].Text, nil
}
