package taskdiscovery

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type PRTemplateChecker struct{}

func (c *PRTemplateChecker) ID() checker.CheckerID  { return "pr_template" }
func (c *PRTemplateChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *PRTemplateChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *PRTemplateChecker) Name() string           { return "Pull Request Template" }
func (c *PRTemplateChecker) Description() string {
	return "Checks that a pull request template exists to standardize PR descriptions and checklists"
}
func (c *PRTemplateChecker) Suggestion() string {
	return "Create .github/PULL_REQUEST_TEMPLATE.md with sections for: description, type of change, testing, and checklist"
}

var prTemplateCandidates = []string{
	".github/PULL_REQUEST_TEMPLATE.md",
	".github/pull_request_template.md",
	"PULL_REQUEST_TEMPLATE.md",
}

func (c *PRTemplateChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	passed, path := checker.FileExistsAny(repo, prTemplateCandidates)
	if passed {
		result.Passed = true
		result.Evidence = "Found PR template: " + path
	} else {
		result.Passed = false
		result.Evidence = "No pull request template found (checked .github/PULL_REQUEST_TEMPLATE.md, .github/pull_request_template.md, PULL_REQUEST_TEMPLATE.md)"
	}

	return result, nil
}
