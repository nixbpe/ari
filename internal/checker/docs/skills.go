package docs

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type SkillsChecker struct {
	Evaluator llm.Evaluator
}

func (c *SkillsChecker) ID() checker.CheckerID  { return "skills" }
func (c *SkillsChecker) Pillar() checker.Pillar { return checker.PillarContextIntent }
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
		"Evaluate the quality of these AI skill definitions for a %s project.\n"+
			"Rule-based finding: %s\nSkill content (truncated):\n%s\n\n"+
			"A good AI skill should have:\n"+
			"1. Clear description of what the skill does\n"+
			"2. Specific trigger conditions (when should the agent use this skill)\n"+
			"3. Actionable step-by-step instructions\n"+
			"4. Project-specific content (not generic boilerplate)\n"+
			"5. Proper metadata (name, description fields in YAML frontmatter if applicable)\n\n"+
			"Are these skills well-defined enough for an AI agent to use effectively?\n"+
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

func (c *SkillsChecker) ruleCheck(repo fs.FS) (bool, string, string) {
	var content string
	count := 0

	entries, err := fs.ReadDir(repo, ".claude/skills")
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			count++
			data, readErr := fs.ReadFile(repo, ".claude/skills/"+e.Name())
			if readErr == nil {
				content += fmt.Sprintf("--- %s ---\n%s\n", e.Name(), string(data))
				if len(content) > 3000 {
					content = content[:3000]
					break
				}
			}
		}
		if count > 0 {
			return true, fmt.Sprintf("Found %d skill(s) in .claude/skills/", count), content
		}
	}

	cursorEntries, cursorErr := fs.ReadDir(repo, ".cursor")
	if cursorErr == nil {
		for _, e := range cursorEntries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".mdc") {
				count++
				data, readErr := fs.ReadFile(repo, ".cursor/"+e.Name())
				if readErr == nil {
					content += fmt.Sprintf("--- %s ---\n%s\n", e.Name(), string(data))
					if len(content) > 3000 {
						content = content[:3000]
						break
					}
				}
			}
		}
		if count > 0 {
			return true, fmt.Sprintf("Found %d .mdc skill file(s) in .cursor/", count), content
		}
	}

	return false, "No AI skills configured", ""
}
