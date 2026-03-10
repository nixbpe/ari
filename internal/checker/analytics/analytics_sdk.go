package analytics

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

// AnalyticsSdkChecker detects analytics SDKs in dependency files.
type AnalyticsSdkChecker struct{}

func (c *AnalyticsSdkChecker) ID() checker.CheckerID  { return "analytics_sdk" }
func (c *AnalyticsSdkChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *AnalyticsSdkChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *AnalyticsSdkChecker) Name() string           { return "Analytics SDK" }
func (c *AnalyticsSdkChecker) Description() string {
	return "Checks that an analytics SDK (Segment, Mixpanel, Amplitude, PostHog, RudderStack) is present in dependencies"
}
func (c *AnalyticsSdkChecker) Suggestion() string {
	return "Add an analytics SDK (e.g. posthog-js, @segment/analytics-next, mixpanel-browser) to track product usage"
}

var analyticsSdkPackages = []string{
	"@segment/analytics-next",
	"mixpanel-browser",
	"amplitude-js",
	"@amplitude/analytics-browser",
	"posthog-js",
	"posthog",
	"@posthog/node",
	"rudder-sdk-js",
	"@rudderstack/",
}

func (c *AnalyticsSdkChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	found, evidence := checker.DepFileContains(repo, lang, analyticsSdkPackages)
	result.Passed = found
	if found {
		result.Evidence = "Analytics SDK found: " + evidence
	} else {
		result.Evidence = "No analytics SDK found in dependencies"
	}
	return result, nil
}
