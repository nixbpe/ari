package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type BuildCmdDocChecker struct {
	Evaluator llm.Evaluator
}

func (c *BuildCmdDocChecker) ID() checker.CheckerID  { return "build_cmd_doc" }
func (c *BuildCmdDocChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *BuildCmdDocChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *BuildCmdDocChecker) Name() string           { return "Build Command Documentation" }
func (c *BuildCmdDocChecker) Description() string {
	return "Checks that build commands are documented in project files"
}
func (c *BuildCmdDocChecker) Suggestion() string {
	return "Document build commands in README.md. Add a Getting Started section with build instructions"
}

func (c *BuildCmdDocChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence, docContent := c.ruleCheck(repo, lang)

	if c.Evaluator == nil {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate whether build commands are documented for a %s repository.\n"+
			"Rule-based finding: %s\nDocumentation content (truncated):\n%s\n\n"+
			"Are build/compile instructions clearly documented?\n"+
			`Respond with JSON only: {"passed": bool, "evidence": "brief explanation", "confidence": 0.0}`,
		lang.String(), evidence, docContent,
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

func (c *BuildCmdDocChecker) ruleCheck(repo fs.FS, lang checker.Language) (bool, string, string) {
	docFiles := []string{"README.md", "CONTRIBUTING.md", "Makefile", "CLAUDE.md", "AGENTS.md"}
	buildCmds := buildCommandsForLang(lang)

	type docEntry struct {
		name string
		data []byte
	}

	var snippets strings.Builder
	var docs []docEntry

	for _, docFile := range docFiles {
		content, err := fs.ReadFile(repo, docFile)
		if err != nil {
			continue
		}
		docs = append(docs, docEntry{docFile, content})
		snippet := string(content)
		if len(snippet) > 2000 {
			snippet = snippet[:2000]
		}
		fmt.Fprintf(&snippets, "--- %s ---\n%s\n", docFile, snippet)
	}

	for _, d := range docs {
		for _, cmd := range buildCmds {
			if bytes.Contains(d.data, []byte(cmd)) {
				return true, fmt.Sprintf("Build command documented in %s: %s", d.name, cmd), snippets.String()
			}
		}
	}

	return false, "No build command documentation found", snippets.String()
}

func buildCommandsForLang(lang checker.Language) []string {
	switch lang {
	case checker.LanguageGo:
		return []string{"go build", "make build", "./dev build"}
	case checker.LanguageTypeScript:
		return []string{"npm run build", "yarn build", "pnpm build"}
	case checker.LanguageJava:
		return []string{"mvn", "gradle", "./gradlew"}
	default:
		return []string{"go build", "make build", "npm run build", "yarn build", "./gradlew", "mvn"}
	}
}
