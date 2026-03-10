package taskdiscovery

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type BacklogStructureDocsChecker struct {
	Evaluator llm.Evaluator
}

func (c *BacklogStructureDocsChecker) ID() checker.CheckerID  { return "backlog_structure_docs" }
func (c *BacklogStructureDocsChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *BacklogStructureDocsChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *BacklogStructureDocsChecker) Name() string           { return "Backlog Structure Documentation" }
func (c *BacklogStructureDocsChecker) Description() string {
	return "Checks that backlog/process documentation exists describing how work is prioritized, triaged, and structured"
}
func (c *BacklogStructureDocsChecker) Suggestion() string {
	return "Document your backlog process in AGENTS.md or docs/process.md: how issues are triaged, priority levels, sprint/milestone structure"
}

var backlogDocCandidates = []string{
	"AGENTS.md",
	"docs/process.md",
	"docs/contributing.md",
}

var backlogKeywords = []string{"priority", "triage", "backlog", "sprint", "milestone"}

func (c *BacklogStructureDocsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		"Evaluate the backlog and process documentation for a %s project.\n"+
			"Rule-based finding: %s\nContent (truncated):\n%s\n\n"+
			"Good backlog structure documentation should describe:\n"+
			"1. How issues/tasks are prioritized (e.g. labels, milestones, priority levels)\n"+
			"2. A triage process for incoming requests\n"+
			"3. Backlog grooming or sprint planning processes\n"+
			"4. Clear definitions of done or acceptance criteria\n\n"+
			"Does this documentation provide meaningful guidance on how work is structured and prioritized?\n"+
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

func (c *BacklogStructureDocsChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, name := range backlogDocCandidates {
		data, err := fs.ReadFile(repo, name)
		if err != nil {
			continue
		}

		content := string(data)
		if len(content) > 4000 {
			content = content[:4000]
		}

		contentLower := strings.ToLower(content)
		for _, keyword := range backlogKeywords {
			if strings.Contains(contentLower, keyword) {
				return true, fmt.Sprintf("Found backlog process documentation in %s (contains %q)", name, keyword), content
			}
		}
	}

	return false, "No backlog structure documentation found with process keywords (priority/triage/backlog/sprint/milestone)", ""
}
