package devenv

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type DevcontainerQualityChecker struct {
	Evaluator llm.Evaluator
}

func (c *DevcontainerQualityChecker) ID() checker.CheckerID  { return "devcontainer_quality" }
func (c *DevcontainerQualityChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *DevcontainerQualityChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *DevcontainerQualityChecker) Name() string           { return "Dev Container Quality" }
func (c *DevcontainerQualityChecker) Description() string {
	return "Evaluates the quality of devcontainer configuration including postCreateCommand, features, and extensions"
}
func (c *DevcontainerQualityChecker) Suggestion() string {
	return "Enhance your devcontainer.json with postCreateCommand for setup automation, features for tooling, and extensions for IDE setup"
}

var devcontainerFiles = []string{".devcontainer/devcontainer.json", ".devcontainer.json"}

func (c *DevcontainerQualityChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence, content := c.ruleCheck(repo)

	if c.Evaluator == nil || content == "" {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate the quality of this devcontainer.json for a %s project.\n"+
			"Rule-based finding: %s\ndevcontainer.json content:\n%s\n\n"+
			"Does this devcontainer have postCreateCommand, features for tooling, and extensions for IDE setup?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang.String(), evidence, content,
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

func (c *DevcontainerQualityChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, path := range devcontainerFiles {
		data, err := fs.ReadFile(repo, path)
		if err != nil {
			continue
		}

		content := string(data)
		if len(content) > 4000 {
			content = content[:4000]
		}

		found, _ := checker.FileContentContains(repo, path, []string{"postCreateCommand"})
		if found {
			return true, fmt.Sprintf("Found %s with postCreateCommand", path), content
		}
		return false, fmt.Sprintf("Found %s but missing postCreateCommand", path), content
	}
	return false, "No devcontainer configuration found", ""
}
