package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type AutomatedPRReviewChecker struct{}

func (c *AutomatedPRReviewChecker) ID() checker.CheckerID  { return "automated_pr_review" }
func (c *AutomatedPRReviewChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *AutomatedPRReviewChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *AutomatedPRReviewChecker) Name() string           { return "Automated PR Review" }
func (c *AutomatedPRReviewChecker) Description() string {
	return "Checks for automated PR review configuration including CODEOWNERS, bots, and CI PR triggers"
}
func (c *AutomatedPRReviewChecker) Suggestion() string {
	return "Set up automated PR review. Add CODEOWNERS and CI checks that run on pull requests"
}

func (c *AutomatedPRReviewChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	staticCandidates := []struct {
		path     string
		evidence string
	}{
		{"CODEOWNERS", "Found CODEOWNERS — automated review configured"},
		{".github/CODEOWNERS", "Found .github/CODEOWNERS — automated review configured"},
		{".github/bots", "Found .github/bots/ — automated bot review configured"},
		{"renovate.json", "Found renovate.json — Renovate bot configured"},
		{".renovaterc", "Found .renovaterc — Renovate bot configured"},
	}

	for _, c := range staticCandidates {
		if _, err := fs.Stat(repo, c.path); err == nil {
			result.Passed = true
			result.Evidence = c.evidence
			return result, nil
		}
	}

	workflowDir := ".github/workflows"
	entries, err := fs.ReadDir(repo, workflowDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			data, err := fs.ReadFile(repo, fmt.Sprintf("%s/%s", workflowDir, entry.Name()))
			if err != nil {
				continue
			}
			if bytes.Contains(data, []byte("pull_request:")) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found PR-triggered CI workflow in %s/%s", workflowDir, entry.Name())
				return result, nil
			}
		}
	}

	result.Passed = false
	result.Evidence = "No automated PR review configured"
	return result, nil
}
