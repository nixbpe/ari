package style

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/llm"
)

type CodeModularizationChecker struct {
	Evaluator llm.Evaluator
}

func (c *CodeModularizationChecker) ID() checker.CheckerID  { return "code_modularization" }
func (c *CodeModularizationChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *CodeModularizationChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *CodeModularizationChecker) Name() string           { return "Code Modularization" }
func (c *CodeModularizationChecker) Description() string {
	return "Checks for module boundary enforcement patterns"
}
func (c *CodeModularizationChecker) Suggestion() string {
	return "Enforce module boundaries. For Go: use internal/ package pattern. For TS: use NX module boundaries"
}

func (c *CodeModularizationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence := c.ruleCheck(repo, lang)

	if c.Evaluator == nil {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate code modularization for a %s repository.\nRule-based finding: %s\nDoes this repository enforce module boundaries?",
		lang.String(), evidence,
	)
	evalResult, err := c.Evaluator.Evaluate(ctx, prompt)
	if err != nil || evalResult == nil {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	result.Passed = evalResult.Passed
	result.Evidence = evalResult.Evidence
	result.Mode = evalResult.Mode
	if result.Mode == "" {
		result.Mode = "llm"
	}
	return result, nil
}

func (c *CodeModularizationChecker) ruleCheck(repo fs.FS, lang checker.Language) (bool, string) {
	switch lang {
	case checker.LanguageGo:
		if _, err := fs.Stat(repo, "internal"); err == nil {
			return true, "Found internal/ directory — Go module boundary pattern used"
		}
		return false, "No module boundary enforcement found"

	case checker.LanguageTypeScript:
		eslintFiles := []string{
			".eslintrc", ".eslintrc.js", ".eslintrc.json", ".eslintrc.yml",
			"eslint.config.js", "eslint.config.mjs",
		}
		for _, f := range eslintFiles {
			content, err := fs.ReadFile(repo, f)
			if err == nil && bytes.Contains(content, []byte("@nx/enforce-module-boundaries")) {
				return true, fmt.Sprintf("Found @nx/enforce-module-boundaries in %s", f)
			}
		}
		if _, err := fs.Stat(repo, "src"); err == nil {
			return true, "Found src/ directory structure — module separation pattern used"
		}
		return false, "No module boundary enforcement found"

	case checker.LanguageJava:
		if _, err := fs.Stat(repo, "module-info.java"); err == nil {
			return true, "Found module-info.java — Java module system used"
		}
		if matches, _ := fs.Glob(repo, "*/module-info.java"); len(matches) > 0 {
			return true, fmt.Sprintf("Found %s — Java module system used", matches[0])
		}
		return false, "No module boundary enforcement found"

	default:
		return false, "No module boundary enforcement found"
	}
}
