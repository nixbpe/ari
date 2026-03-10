package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DepUpdateAutomationChecker struct{}

func (c *DepUpdateAutomationChecker) ID() checker.CheckerID  { return "dep_update_automation" }
func (c *DepUpdateAutomationChecker) Pillar() checker.Pillar { return checker.PillarSecurity }
func (c *DepUpdateAutomationChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *DepUpdateAutomationChecker) Name() string           { return "Dependency Update Automation" }
func (c *DepUpdateAutomationChecker) Description() string {
	return "Checks that automated dependency updates are configured (Dependabot or Renovate)"
}
func (c *DepUpdateAutomationChecker) Suggestion() string {
	return "Configure Dependabot (.github/dependabot.yml) or Renovate (renovate.json) for automated dependency updates"
}

var depUpdateCandidates = []string{
	".github/dependabot.yml",
	".github/dependabot.yaml",
	"renovate.json",
	".renovaterc",
	"renovate.json5",
	".github/renovate.json",
}

func (c *DepUpdateAutomationChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	found, path := checker.FileExistsAny(repo, depUpdateCandidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No dependency update automation found (checked dependabot.yml, renovate.json, .renovaterc)"
	}

	return result, nil
}
