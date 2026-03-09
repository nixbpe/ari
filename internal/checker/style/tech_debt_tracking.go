package style

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type TechDebtTrackingChecker struct{}

func (c *TechDebtTrackingChecker) ID() checker.CheckerID  { return "tech_debt_tracking" }
func (c *TechDebtTrackingChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *TechDebtTrackingChecker) Level() checker.Level   { return checker.LevelAutonomous }
func (c *TechDebtTrackingChecker) Name() string           { return "Tech Debt Tracking" }
func (c *TechDebtTrackingChecker) Description() string {
	return "Checks for tech debt tracking and TODO enforcement"
}
func (c *TechDebtTrackingChecker) Suggestion() string {
	return "Add tech debt tracking. Use linter rules to enforce TODO(JIRA-123) format"
}

func (c *TechDebtTrackingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	workflows, _ := fs.Glob(repo, ".github/workflows/*.yml")
	for _, wf := range workflows {
		content, err := fs.ReadFile(repo, wf)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("TODO")) || bytes.Contains(content, []byte("grep -r TODO")) {
			result.Passed = true
			result.Evidence = "Found TODO tracking enforcement in CI workflow"
			return result, nil
		}
	}

	switch lang {
	case checker.LanguageGo:
		for _, f := range []string{".golangci.yml", ".golangci.yaml"} {
			content, err := fs.ReadFile(repo, f)
			if err == nil {
				if bytes.Contains(content, []byte("godot")) || bytes.Contains(content, []byte("nolintlint")) {
					result.Passed = true
					result.Evidence = fmt.Sprintf("Found TODO tracking linter in %s", f)
					return result, nil
				}
			}
		}

	case checker.LanguageTypeScript:
		for _, f := range []string{
			".eslintrc", ".eslintrc.js", ".eslintrc.json", ".eslintrc.yml",
			"eslint.config.js", "eslint.config.mjs",
		} {
			content, err := fs.ReadFile(repo, f)
			if err == nil && bytes.Contains(content, []byte("no-warning-comments")) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found no-warning-comments rule in %s", f)
				return result, nil
			}
		}
	}

	result.Passed = false
	result.Evidence = "No tech debt tracking configured"
	return result, nil
}
