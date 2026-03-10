package style

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type LintConfigChecker struct{}

func (c *LintConfigChecker) ID() checker.CheckerID  { return "lint_config" }
func (c *LintConfigChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *LintConfigChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *LintConfigChecker) Name() string           { return "Lint Configuration" }
func (c *LintConfigChecker) Description() string    { return "Checks for linter configuration file" }
func (c *LintConfigChecker) Suggestion() string {
	return "Add a linter configuration. For Go: golangci-lint init. For TS: npm init @eslint/config"
}

func (c *LintConfigChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := lintConfigCandidates(lang)

	for _, entry := range candidates {
		if _, err := fs.Stat(repo, entry.file); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — %s configured", entry.file, entry.tool)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No linter configuration found"
	return result, nil
}

type lintEntry struct {
	file string
	tool string
}

func lintConfigCandidates(lang checker.Language) []lintEntry {
	switch lang {
	case checker.LanguageGo:
		return []lintEntry{
			{".golangci.yml", "golangci-lint"},
			{".golangci.yaml", "golangci-lint"},
			{".golangci.json", "golangci-lint"},
			{".golangci.toml", "golangci-lint"},
		}
	case checker.LanguageTypeScript:
		return []lintEntry{
			{".eslintrc", "ESLint"},
			{".eslintrc.js", "ESLint"},
			{".eslintrc.json", "ESLint"},
			{".eslintrc.yml", "ESLint"},
			{"eslint.config.js", "ESLint"},
			{"eslint.config.mjs", "ESLint"},
			{"biome.json", "Biome"},
			{"biome.jsonc", "Biome"},
		}
	case checker.LanguageJava:
		return []lintEntry{
			{"checkstyle.xml", "Checkstyle"},
			{"pmd.xml", "PMD"},
			{".spotbugs.xml", "SpotBugs"},
		}
	default:
		return nil
	}
}
