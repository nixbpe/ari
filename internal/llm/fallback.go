package llm

import "context"

type FallbackEvaluator struct {
	Primary  Evaluator
	Fallback Evaluator
}

func (e *FallbackEvaluator) Evaluate(ctx context.Context, prompt string) (*EvalResult, error) {
	if e == nil {
		return defaultRuleResult(), nil
	}

	if e.Primary != nil {
		res, err := e.Primary.Evaluate(ctx, prompt)
		if err == nil {
			if res == nil {
				res = defaultRuleResult()
			}
			res.Mode = "llm"
			return res, nil
		}
	}

	if e.Fallback != nil {
		res, err := e.Fallback.Evaluate(ctx, prompt)
		if err == nil {
			if res == nil {
				res = defaultRuleResult()
			}
			res.Mode = "rule-based"
			return res, nil
		}
	}

	return defaultRuleResult(), nil
}

func defaultRuleResult() *EvalResult {
	return &EvalResult{
		Passed:     false,
		Evidence:   "no matching rule",
		Confidence: 0.5,
		Mode:       "rule-based",
	}
}
