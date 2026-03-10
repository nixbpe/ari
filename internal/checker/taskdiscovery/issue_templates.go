package taskdiscovery

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type IssueTemplatesChecker struct{}

func (c *IssueTemplatesChecker) ID() checker.CheckerID  { return "issue_templates" }
func (c *IssueTemplatesChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *IssueTemplatesChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *IssueTemplatesChecker) Name() string           { return "Issue Templates" }
func (c *IssueTemplatesChecker) Description() string {
	return "Checks that issue templates exist in .github/ISSUE_TEMPLATE/ to guide contributors in reporting bugs and requesting features"
}
func (c *IssueTemplatesChecker) Suggestion() string {
	return "Create .github/ISSUE_TEMPLATE/ with at least a bug_report.md and feature_request.md template"
}

func (c *IssueTemplatesChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	entries, err := fs.ReadDir(repo, ".github/ISSUE_TEMPLATE")
	if err != nil {
		result.Passed = false
		result.Evidence = "No .github/ISSUE_TEMPLATE directory found"
		return result, nil
	}

	fileCount := 0
	for _, e := range entries {
		if !e.IsDir() {
			fileCount++
		}
	}

	if fileCount >= 1 {
		result.Passed = true
		result.Evidence = fmt.Sprintf("Found %d issue template(s) in .github/ISSUE_TEMPLATE/", fileCount)
	} else {
		result.Passed = false
		result.Evidence = ".github/ISSUE_TEMPLATE/ exists but is empty"
	}

	return result, nil
}
