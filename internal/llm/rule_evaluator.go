package llm

import (
	"context"
	"strings"
)

type Rule struct {
	Pattern  string
	Passed   bool
	Evidence string
}

type RuleBasedEvaluator struct {
	Rules []Rule
}

func (e *RuleBasedEvaluator) Evaluate(_ context.Context, prompt string) (*EvalResult, error) {
	for _, rule := range e.Rules {
		if strings.Contains(prompt, rule.Pattern) {
			return &EvalResult{
				Passed:     rule.Passed,
				Evidence:   rule.Evidence,
				Confidence: 0.8,
				Mode:       "rule-based",
			}, nil
		}
	}

	return defaultRuleResult(), nil
}
