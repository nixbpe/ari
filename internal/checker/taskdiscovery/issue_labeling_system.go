package taskdiscovery

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type IssueLabelingSystemChecker struct{}

func (c *IssueLabelingSystemChecker) ID() checker.CheckerID  { return "issue_labeling_system" }
func (c *IssueLabelingSystemChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *IssueLabelingSystemChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *IssueLabelingSystemChecker) Name() string           { return "Issue Labeling System" }
func (c *IssueLabelingSystemChecker) Description() string {
	return "Checks that a labels configuration file exists to standardize issue and PR labeling"
}
func (c *IssueLabelingSystemChecker) Suggestion() string {
	return "Create .github/labels.yml defining labels for: bug, feature, documentation, priority levels, and component areas"
}

var labelingCandidates = []string{
	".github/labels.yml",
	".github/labels.yaml",
	".github/labels.json",
	"labels.yml",
}

func (c *IssueLabelingSystemChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	passed, path := checker.FileExistsAny(repo, labelingCandidates)
	if passed {
		result.Passed = true
		result.Evidence = "Found labels configuration: " + path
	} else {
		result.Passed = false
		result.Evidence = "No labels configuration found (checked .github/labels.yml, .github/labels.yaml, .github/labels.json, labels.yml)"
	}

	return result, nil
}
