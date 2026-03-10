package analytics

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type TrackingPlanDocsChecker struct{}

func (c *TrackingPlanDocsChecker) ID() checker.CheckerID  { return "tracking_plan_docs" }
func (c *TrackingPlanDocsChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *TrackingPlanDocsChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *TrackingPlanDocsChecker) Name() string           { return "Tracking Plan Documentation" }
func (c *TrackingPlanDocsChecker) Description() string {
	return "Checks for a tracking plan document (docs/tracking-plan.md, avo.json, .avo/) defining analytics events"
}
func (c *TrackingPlanDocsChecker) Suggestion() string {
	return "Create a tracking plan at docs/tracking-plan.md or use Avo (avo.json) to define analytics events"
}

var trackingPlanFiles = []string{
	"docs/tracking-plan.md",
	"docs/analytics.md",
	"tracking-plan.json",
	"avo.json",
}

func (c *TrackingPlanDocsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if found, path := checker.FileExistsAny(repo, trackingPlanFiles); found {
		result.Passed = true
		result.Evidence = "Tracking plan found: " + path
		return result, nil
	}

	if _, err := fs.ReadDir(repo, ".avo"); err == nil {
		result.Passed = true
		result.Evidence = "Tracking plan found: .avo/"
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No tracking plan document found"
	return result, nil
}
