package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type CodeownersChecker struct{}

func (c *CodeownersChecker) ID() checker.CheckerID  { return "codeowners" }
func (c *CodeownersChecker) Pillar() checker.Pillar { return checker.PillarSecurity }
func (c *CodeownersChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *CodeownersChecker) Name() string           { return "CODEOWNERS File" }
func (c *CodeownersChecker) Description() string {
	return "Checks that a CODEOWNERS file exists to define code ownership"
}
func (c *CodeownersChecker) Suggestion() string {
	return "Create .github/CODEOWNERS to assign ownership to code paths for automated review requests"
}

var codeownersCandidates = []string{
	"CODEOWNERS",
	".github/CODEOWNERS",
	"docs/CODEOWNERS",
}

func (c *CodeownersChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	found, path := checker.FileExistsAny(repo, codeownersCandidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No CODEOWNERS file found (checked CODEOWNERS, .github/CODEOWNERS, docs/CODEOWNERS)"
	}

	return result, nil
}
