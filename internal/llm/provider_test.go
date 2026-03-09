package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestOpenAIProvider(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/v1/chat/completions")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("Authorization = %q, want %q", got, "Bearer test-key")
		}

		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		if body.Model != "gpt-4o-mini" {
			t.Fatalf("model = %q, want %q", body.Model, "gpt-4o-mini")
		}
		if len(body.Messages) != 1 || body.Messages[0].Content != "hello" || body.Messages[0].Role != "user" {
			t.Fatalf("messages = %+v, want single user message with content hello", body.Messages)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"openai-response"}}]}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider(Config{APIKey: "test-key", BaseURL: server.URL})
	got, err := p.Complete(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if got != "openai-response" {
		t.Fatalf("Complete() = %q, want %q", got, "openai-response")
	}
}

func TestOpenAIProvider429Retry(t *testing.T) {
	t.Parallel()

	originalSleep := openAISleep
	openAISleep = func(context.Context, time.Duration) error { return nil }
	t.Cleanup(func() { openAISleep = originalSleep })

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count := atomic.AddInt32(&calls, 1)
		if count < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate limited"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	p := NewOpenAIProvider(Config{APIKey: "test-key", BaseURL: server.URL})
	got, err := p.Complete(context.Background(), "retry")
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if got != "ok" {
		t.Fatalf("Complete() = %q, want %q", got, "ok")
	}
	if gotCalls := atomic.LoadInt32(&calls); gotCalls != 3 {
		t.Fatalf("calls = %d, want %d", gotCalls, 3)
	}
}

func TestAnthropicProvider(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/v1/messages")
		}
		if got := r.Header.Get("x-api-key"); got != "anthropic-key" {
			t.Fatalf("x-api-key = %q, want %q", got, "anthropic-key")
		}
		if got := r.Header.Get("anthropic-version"); got != "2023-06-01" {
			t.Fatalf("anthropic-version = %q, want %q", got, "2023-06-01")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"content":[{"text":"anthropic-response"}]}`))
	}))
	defer server.Close()

	p := NewAnthropicProvider(Config{APIKey: "anthropic-key", BaseURL: server.URL})
	got, err := p.Complete(context.Background(), "hello")
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if got != "anthropic-response" {
		t.Fatalf("Complete() = %q, want %q", got, "anthropic-response")
	}
}

func TestOllamaConnectionRefused(t *testing.T) {
	t.Parallel()

	provider := NewOllamaProvider(Config{BaseURL: "http://127.0.0.1:1"})
	_, err := provider.Complete(context.Background(), "hello")
	if err == nil {
		t.Fatal("Complete() error = nil, want connection refused error")
	}

	if !strings.Contains(err.Error(), "Ollama is not running") {
		t.Fatalf("error = %q, want message containing %q", err.Error(), "Ollama is not running")
	}
}

func TestFactorySelectsOpenAI(t *testing.T) {
	t.Parallel()

	p, err := NewProviderFromConfig(Config{Provider: "openai", APIKey: "key"})
	if err != nil {
		t.Fatalf("NewProviderFromConfig() error = %v", err)
	}

	if _, ok := p.(*OpenAIProvider); !ok {
		t.Fatalf("provider type = %T, want *OpenAIProvider", p)
	}
}

func TestFactorySelectsAnthropic(t *testing.T) {
	t.Parallel()

	p, err := NewProviderFromConfig(Config{Provider: "anthropic", APIKey: "key"})
	if err != nil {
		t.Fatalf("NewProviderFromConfig() error = %v", err)
	}

	if _, ok := p.(*AnthropicProvider); !ok {
		t.Fatalf("provider type = %T, want *AnthropicProvider", p)
	}
}

func TestFactorySelectsOllama(t *testing.T) {
	t.Parallel()

	p, err := NewProviderFromConfig(Config{Provider: "ollama"})
	if err != nil {
		t.Fatalf("NewProviderFromConfig() error = %v", err)
	}

	if _, ok := p.(*OllamaProvider); !ok {
		t.Fatalf("provider type = %T, want *OllamaProvider", p)
	}
}
