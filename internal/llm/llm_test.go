package llm

import (
	"context"
	"errors"
	"testing"
)

type testEvaluator struct {
	evaluateFunc func(ctx context.Context, prompt string) (*EvalResult, error)
}

func (t testEvaluator) Evaluate(ctx context.Context, prompt string) (*EvalResult, error) {
	return t.evaluateFunc(ctx, prompt)
}

func TestFallbackUsesPrimary(t *testing.T) {
	fallbackCalled := false
	evaluator := &FallbackEvaluator{
		Primary: testEvaluator{evaluateFunc: func(ctx context.Context, prompt string) (*EvalResult, error) {
			return &EvalResult{Passed: true, Evidence: "primary", Confidence: 0.9}, nil
		}},
		Fallback: testEvaluator{evaluateFunc: func(ctx context.Context, prompt string) (*EvalResult, error) {
			fallbackCalled = true
			return &EvalResult{Passed: false, Evidence: "fallback", Confidence: 0.5}, nil
		}},
	}

	got, err := evaluator.Evaluate(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if fallbackCalled {
		t.Fatal("fallback should not be called when primary succeeds")
	}
	if got == nil {
		t.Fatal("Evaluate() returned nil result")
	}
	if got.Mode != "llm" {
		t.Fatalf("Mode = %q, want %q", got.Mode, "llm")
	}
	if got.Evidence != "primary" {
		t.Fatalf("Evidence = %q, want %q", got.Evidence, "primary")
	}
}

func TestFallbackOnError(t *testing.T) {
	evaluator := &FallbackEvaluator{
		Primary: testEvaluator{evaluateFunc: func(ctx context.Context, prompt string) (*EvalResult, error) {
			return nil, errors.New("primary failed")
		}},
		Fallback: testEvaluator{evaluateFunc: func(ctx context.Context, prompt string) (*EvalResult, error) {
			return &EvalResult{Passed: true, Evidence: "fallback", Confidence: 0.7}, nil
		}},
	}

	got, err := evaluator.Evaluate(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got == nil {
		t.Fatal("Evaluate() returned nil result")
	}
	if got.Mode != "rule-based" {
		t.Fatalf("Mode = %q, want %q", got.Mode, "rule-based")
	}
	if got.Evidence != "fallback" {
		t.Fatalf("Evidence = %q, want %q", got.Evidence, "fallback")
	}
}

func TestFallbackNilPrimary(t *testing.T) {
	evaluator := &FallbackEvaluator{
		Fallback: testEvaluator{evaluateFunc: func(ctx context.Context, prompt string) (*EvalResult, error) {
			return &EvalResult{Passed: false, Evidence: "fallback-only", Confidence: 0.6}, nil
		}},
	}

	got, err := evaluator.Evaluate(context.Background(), "prompt")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got == nil {
		t.Fatal("Evaluate() returned nil result")
	}
	if got.Mode != "rule-based" {
		t.Fatalf("Mode = %q, want %q", got.Mode, "rule-based")
	}
	if got.Evidence != "fallback-only" {
		t.Fatalf("Evidence = %q, want %q", got.Evidence, "fallback-only")
	}
}

func TestRuleBasedEvaluator(t *testing.T) {
	evaluator := &RuleBasedEvaluator{
		Rules: []Rule{
			{Pattern: "must-have", Passed: true, Evidence: "found must-have"},
		},
	}

	got, err := evaluator.Evaluate(context.Background(), "this contains must-have text")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if got == nil {
		t.Fatal("Evaluate() returned nil result")
	}
	if !got.Passed {
		t.Fatal("Passed = false, want true")
	}
	if got.Evidence != "found must-have" {
		t.Fatalf("Evidence = %q, want %q", got.Evidence, "found must-have")
	}
	if got.Mode != "rule-based" {
		t.Fatalf("Mode = %q, want %q", got.Mode, "rule-based")
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("ARI_LLM_PROVIDER", "anthropic")
	t.Setenv("ARI_API_KEY", "test-key")
	t.Setenv("ARI_LLM_MODEL", "test-model")
	t.Setenv("ARI_LLM_BASE_URL", "https://example.com")

	cfg := ConfigFromEnv()

	if cfg.Provider != "anthropic" {
		t.Fatalf("Provider = %q, want %q", cfg.Provider, "anthropic")
	}
	if cfg.APIKey != "test-key" {
		t.Fatalf("APIKey = %q, want %q", cfg.APIKey, "test-key")
	}
	if cfg.Model != "test-model" {
		t.Fatalf("Model = %q, want %q", cfg.Model, "test-model")
	}
	if cfg.BaseURL != "https://example.com" {
		t.Fatalf("BaseURL = %q, want %q", cfg.BaseURL, "https://example.com")
	}
}
