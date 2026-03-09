package llm

import "context"

// Option is a functional option for LLM completions
type Option func(*options)

type options struct {
	temperature float64
	maxTokens   int
	jsonOutput  bool
}

// WithTemperature sets the sampling temperature
func WithTemperature(t float64) Option {
	return func(o *options) { o.temperature = t }
}

// WithMaxTokens sets the maximum token limit
func WithMaxTokens(n int) Option {
	return func(o *options) { o.maxTokens = n }
}

// WithJSONOutput requests structured JSON output
func WithJSONOutput(enabled bool) Option {
	return func(o *options) { o.jsonOutput = enabled }
}

// Provider is an LLM API provider
type Provider interface {
	Complete(ctx context.Context, prompt string, opts ...Option) (string, error)
	Name() string
}

// EvalResult holds the result of an LLM evaluation
type EvalResult struct {
	Passed     bool
	Evidence   string
	Confidence float64
}

// Evaluator evaluates criteria using LLM
type Evaluator interface {
	Evaluate(ctx context.Context, prompt string) (*EvalResult, error)
}

// Config holds LLM provider configuration
type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}
