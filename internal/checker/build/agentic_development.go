package build

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type AgenticDevelopmentChecker struct {
	Evaluator llm.Evaluator
}

func (c *AgenticDevelopmentChecker) ID() checker.CheckerID  { return "agentic_development" }
func (c *AgenticDevelopmentChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *AgenticDevelopmentChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *AgenticDevelopmentChecker) Name() string           { return "Agentic Development Support" }
func (c *AgenticDevelopmentChecker) Description() string {
	return "Checks for AI agent documentation files that guide automated development tools"
}
func (c *AgenticDevelopmentChecker) Suggestion() string {
	return "Create AGENTS.md documenting: build commands, test commands, architecture overview, coding conventions for AI agents"
}

var agenticDocCandidates = []struct {
	path   string
	isFile bool
}{
	{"AGENTS.md", true},
	{"CLAUDE.md", true},
	{".cursor/rules", false},
	{".github/copilot-instructions.md", true},
}

func (c *AgenticDevelopmentChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		"Evaluate if this AI agent documentation effectively supports agentic development for a %s project.\n"+
			"Rule-based finding: %s\nContent (truncated):\n%s\n\n"+
			"For effective agentic development, the documentation should:\n"+
			"1. Provide step-by-step development workflows (not just commands)\n"+
			"2. Include constraints or prohibited actions (what NOT to do)\n"+
			"3. Describe environment setup requirements\n"+
			"4. Be project-specific (not a generic template)\n"+
			"5. Be actionable enough for an AI agent to follow without human help\n\n"+
			"Does this documentation enable an AI coding agent to independently develop features in this repo?\n"+
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

func (c *AgenticDevelopmentChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	for _, candidate := range agenticDocCandidates {
		if candidate.isFile {
			data, err := fs.ReadFile(repo, candidate.path)
			if err != nil {
				continue
			}
			content := string(data)
			if len(content) > 4000 {
				content = content[:4000]
			}
			return true, fmt.Sprintf("Found %s — AI agent documentation present", candidate.path), content
		}

		// Directory check (e.g., .cursor/rules)
		if _, err := fs.Stat(repo, candidate.path); err == nil {
			// Try to read files from the directory for LLM content
			var content string
			entries, dirErr := fs.ReadDir(repo, candidate.path)
			if dirErr == nil {
				for _, e := range entries {
					if e.IsDir() {
						continue
					}
					data, err := fs.ReadFile(repo, candidate.path+"/"+e.Name())
					if err == nil {
						content += fmt.Sprintf("--- %s ---\n%s\n", e.Name(), string(data))
						if len(content) > 4000 {
							content = content[:4000]
							break
						}
					}
				}
			}
			return true, fmt.Sprintf("Found %s — AI agent documentation present", candidate.path), content
		}
	}

	// Also check for .claude directory
	if _, err := fs.Stat(repo, ".claude"); err == nil {
		return true, "Found .claude — AI agent documentation present", ""
	}

	return false, "No AI agent documentation found", ""
}
