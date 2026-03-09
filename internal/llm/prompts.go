package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// NamingConsistencyPrompt builds a structured prompt for evaluating naming consistency.
func NamingConsistencyPrompt(lang, sampleFiles string) string {
	return fmt.Sprintf(
		"Evaluate naming consistency for a %s project.\nSample files: %s\n\n"+
			"Check if naming conventions are consistently followed.\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang, sampleFiles,
	)
}

// CodeModularizationPrompt builds a structured prompt for evaluating code modularization.
func CodeModularizationPrompt(lang, ruleEvidence string) string {
	return fmt.Sprintf(
		"Evaluate code modularization for a %s repository.\nRule-based finding: %s\n\n"+
			"Does this repository enforce module boundaries?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang, ruleEvidence,
	)
}

// ProviderEvaluator adapts a Provider to the Evaluator interface.
// It calls Provider.Complete and parses the JSON response.
type ProviderEvaluator struct {
	Provider Provider
}

func (e *ProviderEvaluator) Evaluate(ctx context.Context, prompt string) (*EvalResult, error) {
	if e == nil || e.Provider == nil {
		return nil, fmt.Errorf("no provider configured")
	}

	resp, err := e.Provider.Complete(ctx, prompt, WithJSONOutput(true), WithTemperature(0.1))
	if err != nil {
		return nil, fmt.Errorf("LLM complete: %w", err)
	}

	var parsed struct {
		Passed     bool    `json:"passed"`
		Evidence   string  `json:"evidence"`
		Confidence float64 `json:"confidence"`
	}

	jsonStr := extractJSON(resp)
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("parse LLM response %q: %w", jsonStr, err)
	}

	return &EvalResult{
		Passed:     parsed.Passed,
		Evidence:   parsed.Evidence,
		Confidence: parsed.Confidence,
		Mode:       "llm",
	}, nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
