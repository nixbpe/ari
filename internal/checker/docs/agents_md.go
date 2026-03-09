package docs

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type AgentsMdChecker struct{}

func (c *AgentsMdChecker) ID() checker.CheckerID  { return "agents_md" }
func (c *AgentsMdChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
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
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	for _, name := range agentDocCandidates {
		data, err := fs.ReadFile(repo, name)
		if err != nil {
			continue
		}

		if len(data) < 100 {
			result.Passed = false
			result.Evidence = fmt.Sprintf("Found %s but content is too short (%d chars, need 100+)", name, len(data))
			return result, nil
		}

		result.Passed = true
		result.Evidence = fmt.Sprintf("Found %s — AI agent documentation present", name)
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No AI agent documentation found"
	return result, nil
}
