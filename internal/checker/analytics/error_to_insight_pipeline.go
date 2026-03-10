package analytics

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ErrorToInsightPipelineChecker struct{}

func (c *ErrorToInsightPipelineChecker) ID() checker.CheckerID { return "error_to_insight_pipeline" }
func (c *ErrorToInsightPipelineChecker) Pillar() checker.Pillar {
	return checker.PillarProductAnalytics
}
func (c *ErrorToInsightPipelineChecker) Level() checker.Level { return checker.LevelOptimized }
func (c *ErrorToInsightPipelineChecker) Name() string         { return "Error-to-Insight Pipeline" }
func (c *ErrorToInsightPipelineChecker) Description() string {
	return "Checks for error reporting pipelines: Sentry CLI in CI workflows, sentry.properties config, or automated issue creation"
}
func (c *ErrorToInsightPipelineChecker) Suggestion() string {
	return "Add error-to-insight pipeline: integrate sentry-cli in CI, add .sentry.properties, or use gh issue create to turn errors into tracked issues"
}

var sentryConfigFiles = []string{".sentry.properties", "sentry.properties"}

var sentryAndIssueCIKeywords = []string{
	"sentry-cli",
	"sentry upload",
	"@sentry/wizard",
	"gh issue create",
	"create issue",
}

func (c *ErrorToInsightPipelineChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if found, path := checker.FileExistsAny(repo, sentryConfigFiles); found {
		result.Passed = true
		result.Evidence = "Sentry config found: " + path
		return result, nil
	}

	if found, evidence := checker.CIWorkflowContains(repo, sentryAndIssueCIKeywords); found {
		result.Passed = true
		result.Evidence = "Error pipeline in CI: " + evidence
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No error-to-insight pipeline found (no sentry config, no sentry-cli/gh issue create in CI)"
	return result, nil
}
