package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
