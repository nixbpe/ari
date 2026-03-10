package analytics

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type ProductMetricsDocsChecker struct {
	Evaluator llm.Evaluator
}

func (c *ProductMetricsDocsChecker) ID() checker.CheckerID  { return "product_metrics_docs" }
func (c *ProductMetricsDocsChecker) Pillar() checker.Pillar { return checker.PillarProductAnalytics }
func (c *ProductMetricsDocsChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *ProductMetricsDocsChecker) Name() string           { return "Product Metrics Documentation" }
func (c *ProductMetricsDocsChecker) Description() string {
	return "Checks that product metrics, KPIs, or north-star metrics are documented"
}
func (c *ProductMetricsDocsChecker) Suggestion() string {
	return "Create docs/metrics.md or docs/kpis.md documenting north-star metric, KPIs, and conversion/retention targets"
}

var productMetricsCandidates = []string{
	"docs/metrics.md",
	"docs/kpis.md",
	"docs/north-star.md",
	"docs/product-metrics.md",
}

var productMetricsKeywords = []string{"metric", "kpi", "north star", "conversion", "retention"}

func (c *ProductMetricsDocsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence, content := c.ruleCheck(repo)

	if c.Evaluator == nil || content == "" {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate the quality of this product metrics documentation for a %s project.\n"+
			"Rule-based finding: %s\nDocument content (truncated):\n%s\n\n"+
			"Does this document clearly define product metrics, KPIs, north-star metric, or success criteria?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang.String(), evidence, content,
	)
	evalResult, err := c.Evaluator.Evaluate(ctx, prompt)
	if err != nil || evalResult == nil {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	result.Passed = evalResult.Passed
	result.Evidence = evalResult.Evidence
	result.Mode = evalResult.Mode
	if result.Mode == "" {
		result.Mode = "llm"
	}
	return result, nil
}

func (c *ProductMetricsDocsChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	found, path := checker.FileExistsAny(repo, productMetricsCandidates)
	if !found {
		return false, "No product metrics documentation found", ""
	}

	data, err := fs.ReadFile(repo, path)
	if err != nil {
		return false, "Found " + path + " but could not read it", ""
	}

	content := string(data)
	if len(content) > 4000 {
		content = content[:4000]
	}

	lower := strings.ToLower(content)
	for _, kw := range productMetricsKeywords {
		if strings.Contains(lower, kw) {
			return true, fmt.Sprintf("Found %s with product metrics content (keyword: %q)", path, kw), content
		}
	}

	return false, fmt.Sprintf("Found %s but missing metrics keywords (metric, kpi, north star, conversion, retention)", path), content
}
