package docs

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/bbik/ari/internal/checker"
)

type SkillsChecker struct{}

func (c *SkillsChecker) ID() checker.CheckerID  { return "skills" }
func (c *SkillsChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *SkillsChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *SkillsChecker) Name() string           { return "AI Skills Configured" }
func (c *SkillsChecker) Description() string {
	return "Checks for AI skill files in .claude/skills/ or .cursor/ that teach agents project workflows"
}
func (c *SkillsChecker) Suggestion() string {
	return "Create AI skills in .claude/skills/ to teach agents project-specific workflows"
}

func (c *SkillsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	count := 0

	// Check .claude/skills/
	entries, err := fs.ReadDir(repo, ".claude/skills")
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				count++
			}
		}
		if count > 0 {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %d skill(s) in .claude/skills/", count)
			return result, nil
		}
	}

	// Check .cursor/ for .mdc files
	cursorEntries, cursorErr := fs.ReadDir(repo, ".cursor")
	if cursorErr == nil {
		for _, e := range cursorEntries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".mdc") {
				count++
			}
		}
		if count > 0 {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %d .mdc skill file(s) in .cursor/", count)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No AI skills configured"
	return result, nil
}
