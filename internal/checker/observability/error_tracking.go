package observability

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ErrorTrackingChecker struct{}

func (c *ErrorTrackingChecker) ID() checker.CheckerID  { return "error_tracking" }
func (c *ErrorTrackingChecker) Pillar() checker.Pillar { return checker.PillarVerification }
func (c *ErrorTrackingChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *ErrorTrackingChecker) Name() string           { return "Error Tracking" }
func (c *ErrorTrackingChecker) Description() string {
	return "Checks for an error tracking service integration in dependencies"
}
func (c *ErrorTrackingChecker) Suggestion() string {
	return "Integrate an error tracking service (e.g., Sentry, Rollbar, Bugsnag) to capture and monitor errors in production"
}

func (c *ErrorTrackingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	packages := []string{
		"sentry",
		"rollbar",
		"bugsnag",
		"honeybadger",
	}

	found, evidence := checker.DepFileContains(repo, lang, packages)
	result.Passed = found
	if found {
		result.Evidence = "Found error tracking integration: " + evidence
	} else {
		result.Evidence = "No error tracking service found in dependencies"
	}
	return result, nil
}
