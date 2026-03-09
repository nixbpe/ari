package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNamingConsistencyPrompt(t *testing.T) {
	prompt := NamingConsistencyPrompt("Go", "main.go, scanner.go")
	if !strings.Contains(prompt, "Go") {
		t.Error("prompt missing language")
	}
	if !strings.Contains(prompt, "main.go") {
		t.Error("prompt missing sample files")
	}
	if !strings.Contains(prompt, "passed") {
		t.Error("prompt missing JSON schema hint")
	}
}

func TestCodeModularizationPrompt(t *testing.T) {
	prompt := CodeModularizationPrompt("TypeScript", "Found src/ directory")
	if !strings.Contains(prompt, "TypeScript") {
		t.Error("prompt missing language")
	}
	if !strings.Contains(prompt, "Found src/") {
		t.Error("prompt missing rule evidence")
	}
}

func TestProviderEvaluatorSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"choices": [{"message": {"content": "{\"passed\": true, \"evidence\": \"looks good\", \"confidence\": 0.9}"}}]
		}`))
	}))
	defer srv.Close()

	provider := NewOpenAIProvider(Config{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})
	eval := &ProviderEvaluator{Provider: provider}

	result, err := eval.Evaluate(context.Background(), "test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Error("expected passed=true")
	}
	if result.Evidence != "looks good" {
		t.Errorf("expected evidence 'looks good', got %q", result.Evidence)
	}
	if result.Mode != "llm" {
		t.Errorf("expected mode 'llm', got %q", result.Mode)
	}
}

func TestProviderEvaluatorNilProvider(t *testing.T) {
	eval := &ProviderEvaluator{Provider: nil}
	_, err := eval.Evaluate(context.Background(), "test")
	if err == nil {
		t.Error("expected error for nil provider")
	}
}

func TestExtractJSON(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`{"passed": true}`, `{"passed": true}`},
		{`Here is the result: {"passed": false, "evidence": "bad"} done`, `{"passed": false, "evidence": "bad"}`},
		{`no json here`, `no json here`},
	}
	for _, tc := range cases {
		got := extractJSON(tc.input)
		if got != tc.expected {
			t.Errorf("extractJSON(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestLLMWiringMode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"choices": [{"message": {"content": "{\"passed\": true, \"evidence\": \"llm says ok\", \"confidence\": 0.95}"}}]
		}`))
	}))
	defer srv.Close()

	provider := NewOpenAIProvider(Config{APIKey: "test-key", BaseURL: srv.URL})
	eval := &FallbackEvaluator{
		Primary:  &ProviderEvaluator{Provider: provider},
		Fallback: &RuleBasedEvaluator{},
	}

	result, err := eval.Evaluate(context.Background(), "check naming")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != "llm" {
		t.Errorf("expected mode 'llm', got %q", result.Mode)
	}
}

func TestNoLLMFallback(t *testing.T) {
	eval := &FallbackEvaluator{
		Primary:  nil,
		Fallback: &RuleBasedEvaluator{Rules: []Rule{{Pattern: "naming", Passed: true, Evidence: "rule says ok"}}},
	}

	result, err := eval.Evaluate(context.Background(), "check naming consistency")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != "rule-based" {
		t.Errorf("expected mode 'rule-based', got %q", result.Mode)
	}
}
