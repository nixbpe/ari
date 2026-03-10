package build

import (
	"bytes"
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type FeatureFlagInfrastructureChecker struct{}

func (c *FeatureFlagInfrastructureChecker) ID() checker.CheckerID {
	return "feature_flag_infrastructure"
}
func (c *FeatureFlagInfrastructureChecker) Pillar() checker.Pillar { return checker.PillarVerification }
func (c *FeatureFlagInfrastructureChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *FeatureFlagInfrastructureChecker) Name() string           { return "Feature Flag Infrastructure" }
func (c *FeatureFlagInfrastructureChecker) Description() string {
	return "Checks that a feature flag library is configured in the project"
}
func (c *FeatureFlagInfrastructureChecker) Suggestion() string {
	return "Add feature flags. Consider OpenFeature SDK for language-agnostic feature flags"
}

var goFeatureFlagLibs = []string{"launchdarkly", "unleash", "openfeature", "flipt", "flagsmith"}
var tsFeatureFlagLibs = []string{"@launchdarkly", "@openfeature", "unleash-client", "flagsmith"}

func (c *FeatureFlagInfrastructureChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if goMod, err := fs.ReadFile(repo, "go.mod"); err == nil {
		for _, lib := range goFeatureFlagLibs {
			if bytes.Contains(goMod, []byte(lib)) {
				result.Passed = true
				result.Evidence = "Found " + lib + " in go.mod — feature flag infrastructure configured"
				return result, nil
			}
		}
	}

	if pkgJSON, err := fs.ReadFile(repo, "package.json"); err == nil {
		for _, lib := range tsFeatureFlagLibs {
			if bytes.Contains(pkgJSON, []byte(lib)) {
				result.Passed = true
				result.Evidence = "Found " + lib + " in package.json — feature flag infrastructure configured"
				return result, nil
			}
		}
	}

	for _, buildFile := range []string{"pom.xml", "build.gradle"} {
		if content, err := fs.ReadFile(repo, buildFile); err == nil {
			for _, lib := range goFeatureFlagLibs {
				if bytes.Contains(content, []byte(lib)) {
					result.Passed = true
					result.Evidence = "Found " + lib + " in " + buildFile + " — feature flag infrastructure configured"
					return result, nil
				}
			}
		}
	}

	result.Passed = false
	result.Evidence = "No feature flag infrastructure found"
	return result, nil
}
