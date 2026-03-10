package docs

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type AgentsMdChecker struct {
	Evaluator llm.Evaluator
}

func (c *AgentsMdChecker) ID() checker.CheckerID  { return "agents_md" }
func (c *AgentsMdChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
func (c *AgentsMdChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *AgentsMdChecker) Name() string           { return "AI Agent Documentation" }
func (c *AgentsMdChecker) Description() string {
	return "Checks for AI agent documentation files (AGENTS.md, CLAUDE.md, .cursor/rules, etc.)"
}
func (c *AgentsMdChecker) Suggestion() string {
	return "Create AGENTS.md documenting: build commands, test commands, architecture overview, coding conventions for AI agents"
}

var agentDocCandidates = []string{
	"AGENTS.md",
	"CLAUDE.md",
	".cursor/rules",
	".github/copilot-instructions.md",
}

func (c *AgentsMdChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		"Evaluate the quality of this AI agent documentation for a %s project.\n"+
			"Rule-based finding: %s\nContent (truncated):\n%s\n\n"+
			"A good agent doc (AGENTS.md/CLAUDE.md) should have:\n"+
			"1. Build commands (exact, backtick-wrapped, copy-paste ready)\n"+
			"2. Test commands (how to run tests)\n"+
			"3. Architecture overview (project structure, key modules)\n"+
			"4. Coding conventions (naming, patterns, style rules)\n"+
			"5. Specific, actionable instructions (not vague prose)\n\n"+
			"Does this doc provide enough guidance for an AI coding agent to work effectively in this repo?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation of what's present/missing", "confidence": 0.0}`,
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

func (c *AgentsMdChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, name := range agentDocCandidates {
		data, err := fs.ReadFile(repo, name)
		if err != nil {
			continue
		}

		content := string(data)
		if len(content) > 4000 {
			content = content[:4000]
		}

		if len(data) < 100 {
			return false, fmt.Sprintf("Found %s but content is too short (%d chars, need 100+)", name, len(data)), content
		}

		return true, fmt.Sprintf("Found %s — AI agent documentation present", name), content
	}

	return false, "No AI agent documentation found", ""
}
