package taskdiscovery

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ContributingGuideChecker struct{}

func (c *ContributingGuideChecker) ID() checker.CheckerID  { return "contributing_guide" }
func (c *ContributingGuideChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *ContributingGuideChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *ContributingGuideChecker) Name() string           { return "Contributing Guide" }
func (c *ContributingGuideChecker) Description() string {
	return "Checks that a CONTRIBUTING.md or equivalent guide exists to help contributors understand the workflow"
}
func (c *ContributingGuideChecker) Suggestion() string {
	return "Create a CONTRIBUTING.md that describes: how to report issues, how to submit PRs, coding standards, and the development workflow"
}

var contributingCandidates = []string{
	"CONTRIBUTING.md",
	"CONTRIBUTING.rst",
	"docs/contributing.md",
	".github/CONTRIBUTING.md",
}

func (c *ContributingGuideChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	passed, path := checker.FileExistsAny(repo, contributingCandidates)
	if passed {
		result.Passed = true
		result.Evidence = "Found contributing guide: " + path
	} else {
		result.Passed = false
		result.Evidence = "No contributing guide found (checked CONTRIBUTING.md, CONTRIBUTING.rst, docs/contributing.md, .github/CONTRIBUTING.md)"
	}

	return result, nil
}
