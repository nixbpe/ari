package style

import (
	"context"
	"encoding/json"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type PreCommitHooksChecker struct{}

func NewPreCommitHooksChecker() *PreCommitHooksChecker {
	return &PreCommitHooksChecker{}
}

func (c *PreCommitHooksChecker) ID() checker.CheckerID  { return "pre_commit_hooks" }
func (c *PreCommitHooksChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *PreCommitHooksChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *PreCommitHooksChecker) Name() string           { return "Pre-commit Hooks" }
func (c *PreCommitHooksChecker) Description() string {
	return "Checks if pre-commit hooks are configured for the project"
}

const preCommitSuggestion = "Add pre-commit hooks. Install: pip install pre-commit && pre-commit install or use Husky for JS projects"

func (c *PreCommitHooksChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	for _, candidate := range []string{".pre-commit-config.yaml", "lefthook.yml", "lefthook.yaml"} {
		if _, err := fs.Stat(repo, candidate); err == nil {
			return c.pass(candidate), nil
		}
	}

	if _, err := fs.Stat(repo, ".husky"); err == nil {
		return c.pass(".husky/"), nil
	}

	if found, name := checkLintStaged(repo); found {
		return c.pass(name), nil
	}

	return &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Passed:     false,
		Evidence:   "No pre-commit hook configuration found",
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: preCommitSuggestion,
	}, nil
}

func (c *PreCommitHooksChecker) pass(file string) *checker.Result {
	return &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Passed:     true,
		Evidence:   "Found " + file + " — pre-commit hooks configured",
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: preCommitSuggestion,
	}
}

func checkLintStaged(repo fs.FS) (bool, string) {
	data, err := fs.ReadFile(repo, "package.json")
	if err != nil {
		return false, ""
	}

	var pkg struct {
		LintStaged interface{} `json:"lint-staged"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false, ""
	}

	if pkg.LintStaged != nil {
		return true, "package.json (lint-staged)"
	}
	return false, ""
}
