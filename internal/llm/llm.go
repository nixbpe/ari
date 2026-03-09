package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Option func(*callConfig)

type callConfig struct {
	Temperature float64
	MaxTokens   int
	JSONOutput  bool
}

func WithTemperature(t float64) Option {
	return func(cfg *callConfig) { cfg.Temperature = t }
}

func WithMaxTokens(n int) Option {
	return func(cfg *callConfig) { cfg.MaxTokens = n }
}

func WithJSONOutput(enabled bool) Option {
	return func(cfg *callConfig) { cfg.JSONOutput = enabled }
}

type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}

type Provider interface {
	Complete(ctx context.Context, prompt string, opts ...Option) (string, error)
	Name() string
}

type EvalResult struct {
	Passed     bool
	Evidence   string
	Confidence float64
	Mode       string
}

type Evaluator interface {
	Evaluate(ctx context.Context, prompt string) (*EvalResult, error)
}

func NewProviderFromConfig(cfg Config) (Provider, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))

	switch provider {
	case "mock":
		return &MockProvider{}, nil
	case "openai", "anthropic", "ollama":
		return nil, fmt.Errorf("provider %q not implemented yet", provider)
	default:
		return nil, fmt.Errorf("unknown provider: %q", cfg.Provider)
	}
}

func ConfigFromEnv() Config {
	return Config{
		Provider: os.Getenv("ARI_LLM_PROVIDER"),
		APIKey:   os.Getenv("ARI_API_KEY"),
		Model:    os.Getenv("ARI_LLM_MODEL"),
		BaseURL:  os.Getenv("ARI_LLM_BASE_URL"),
	}
}
