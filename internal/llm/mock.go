package llm

import "context"

type MockProvider struct {
	CompleteFunc func(ctx context.Context, prompt string, opts ...Option) (string, error)
}

func (m *MockProvider) Complete(ctx context.Context, prompt string, opts ...Option) (string, error) {
	if m != nil && m.CompleteFunc != nil {
		return m.CompleteFunc(ctx, prompt, opts...)
	}

	return `{"passed": true, "evidence": "mock", "confidence": 0.9}`, nil
}

func (m *MockProvider) Name() string {
	return "mock"
}
