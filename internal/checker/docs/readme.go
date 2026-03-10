package docs

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type ReadmeChecker struct {
	Evaluator llm.Evaluator
}

func (c *ReadmeChecker) ID() checker.CheckerID  { return "readme" }
func (c *ReadmeChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *ReadmeChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *ReadmeChecker) Name() string           { return "README Exists" }
func (c *ReadmeChecker) Description() string {
	return "Checks that a README file exists and has meaningful content (50+ chars with headings)"
}
func (c *ReadmeChecker) Suggestion() string {
	return "Create a README.md with: project description, installation, usage, contributing guide"
}

var readmeCandidates = []string{"README.md", "README", "README.rst", "README.txt"}

func (c *ReadmeChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence, readmeContent := c.ruleCheck(repo)

	if c.Evaluator == nil || readmeContent == "" {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate the quality of this README for a %s project.\n"+
			"Rule-based finding: %s\nREADME content (truncated):\n%s\n\n"+
			"Does this README provide a meaningful project description, installation, and usage instructions?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang.String(), evidence, readmeContent,
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

func (c *ReadmeChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, name := range readmeCandidates {
		data, err := fs.ReadFile(repo, name)
		if err != nil {
			continue
		}

		content := string(data)
		if len(content) > 4000 {
			content = content[:4000]
		}

		if len(data) < 50 {
			return false, fmt.Sprintf("%s exists but has insufficient content", name), content
		}

		hasHeading := bytes.Contains(data, []byte("#")) || bytes.Contains(data, []byte("====="))
		if !hasHeading {
			return false, fmt.Sprintf("%s exists but has insufficient content", name), content
		}

		return true, fmt.Sprintf("Found %s (%d bytes) with headings", name, len(data)), content
	}

	return false, "No README found", ""
}
