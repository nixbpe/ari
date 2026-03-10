package observability

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type RunbooksDocumentedChecker struct {
	Evaluator llm.Evaluator
}

func (c *RunbooksDocumentedChecker) ID() checker.CheckerID  { return "runbooks_documented" }
func (c *RunbooksDocumentedChecker) Pillar() checker.Pillar { return checker.PillarObservability }
func (c *RunbooksDocumentedChecker) Level() checker.Level   { return checker.LevelAutonomous }
func (c *RunbooksDocumentedChecker) Name() string           { return "Runbooks Documented" }
func (c *RunbooksDocumentedChecker) Description() string {
	return "Checks that operational runbooks are present and contain meaningful content"
}
func (c *RunbooksDocumentedChecker) Suggestion() string {
	return "Create runbooks in docs/runbooks/, runbooks/, or RUNBOOK.md describing incident response procedures"
}

var runbookCandidates = []string{
	"docs/runbooks",
	"runbooks",
	"RUNBOOK.md",
	"docs/ops",
}

func (c *RunbooksDocumentedChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		"Evaluate the quality of these runbooks for a %s project.\n"+
			"Rule-based finding: %s\nRunbook content (truncated):\n%s\n\n"+
			"Do these runbooks provide meaningful operational procedures for incident response?\n"+
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

func (c *RunbooksDocumentedChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, candidate := range runbookCandidates {
		if info, err := fs.Stat(repo, candidate); err == nil {
			if info.IsDir() {
				entries, dirErr := fs.ReadDir(repo, candidate)
				if dirErr == nil && len(entries) > 0 {
					return true, fmt.Sprintf("Found runbooks directory %s with %d entries", candidate, len(entries)), candidate
				}
			} else {
				data, readErr := fs.ReadFile(repo, candidate)
				if readErr == nil && len(data) > 50 {
					content := string(data)
					if len(content) > 4000 {
						content = content[:4000]
					}
					return true, fmt.Sprintf("Found %s (%d bytes)", candidate, len(data)), content
				}
			}
		}
	}
	return false, "No runbooks found (checked docs/runbooks/, runbooks/, RUNBOOK.md, docs/ops/)", ""
}
