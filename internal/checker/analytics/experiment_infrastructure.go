package analytics

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ExperimentInfrastructureChecker struct{}

func (c *ExperimentInfrastructureChecker) ID() checker.CheckerID { return "experiment_infrastructure" }
func (c *ExperimentInfrastructureChecker) Pillar() checker.Pillar {
	return checker.PillarContextIntent
}
func (c *ExperimentInfrastructureChecker) Level() checker.Level { return checker.LevelStandardized }
func (c *ExperimentInfrastructureChecker) Name() string         { return "Experiment Infrastructure" }
func (c *ExperimentInfrastructureChecker) Description() string {
	return "Checks for A/B testing or feature flag SDKs (GrowthBook, Statsig, Optimizely, Split, LaunchDarkly)"
}
func (c *ExperimentInfrastructureChecker) Suggestion() string {
	return "Add an experimentation SDK (e.g. @growthbook/growthbook, statsig-js, @launchdarkly/) to run A/B tests"
}

var experimentPackages = []string{
	"@growthbook/growthbook",
	"growthbook",
	"statsig-js",
	"@statsig/",
	"@optimizely/optimizely-sdk",
	"@optimizely/react-sdk",
	"@splitsoftware/splitio",
	"@launchdarkly/",
}

func (c *ExperimentInfrastructureChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	found, evidence := checker.DepFileContains(repo, lang, experimentPackages)
	result.Passed = found
	if found {
		result.Evidence = "Experiment SDK found: " + evidence
	} else {
		result.Evidence = "No experiment/A/B testing SDK found in dependencies"
	}
	return result, nil
}
