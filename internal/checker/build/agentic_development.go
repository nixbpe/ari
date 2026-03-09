package build

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type AgenticDevelopmentChecker struct{}

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

func (c *AgenticDevelopmentChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{
		"AGENTS.md",
		"CLAUDE.md",
		".cursor/rules",
		".github/copilot-instructions.md",
		".claude",
	}

	for _, candidate := range candidates {
		if _, err := fs.Stat(repo, candidate); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — AI agent documentation present", candidate)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No AI agent documentation found"
	return result, nil
}
