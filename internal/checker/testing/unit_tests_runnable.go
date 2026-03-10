package checktesting

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type UnitTestsRunnableChecker struct {
	Evaluator llm.Evaluator
}

func (c *UnitTestsRunnableChecker) ID() checker.CheckerID  { return "unit_tests_runnable" }
func (c *UnitTestsRunnableChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *UnitTestsRunnableChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *UnitTestsRunnableChecker) Name() string           { return "Unit Tests Runnable" }
func (c *UnitTestsRunnableChecker) Description() string {
	return "Checks that test commands are documented in project files"
}
func (c *UnitTestsRunnableChecker) Suggestion() string {
	return "Document test commands in README.md and add a Makefile test target"
}

func (c *UnitTestsRunnableChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
	}

	passed, evidence, content := c.ruleCheck(repo, lang)

	if c.Evaluator == nil || content == "" {
		result.Passed = passed
		result.Evidence = evidence
		result.Mode = "rule-based"
		return result, nil
	}

	prompt := fmt.Sprintf(
		"Evaluate the test documentation quality for a %s project.\n"+
			"Rule-based finding: %s\nDocumentation content (truncated):\n%s\n\n"+
			"Good test documentation should:\n"+
			"1. Include exact test commands (backtick-wrapped, copy-paste ready)\n"+
			"2. Explain how to run specific test subsets (single file, single test)\n"+
			"3. Mention any required setup or environment variables\n"+
			"4. Document common test flags or options\n"+
			"5. Be clear enough for an AI agent to run tests without guessing\n\n"+
			"Does this documentation adequately explain how to run the project's tests?\n"+
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

func (c *UnitTestsRunnableChecker) ruleCheck(repo fs.FS, lang checker.Language) (bool, string, string) {
	testCmds := testCommandsForLang(lang)
	docFiles := []string{"README.md", "CONTRIBUTING.md", "Makefile", "CLAUDE.md", "AGENTS.md"}

	var collectedContent string

	for _, docFile := range docFiles {
		content, err := fs.ReadFile(repo, docFile)
		if err != nil {
			continue
		}

		snippet := string(content)
		if len(snippet) > 3000 {
			snippet = snippet[:3000]
		}
		collectedContent += fmt.Sprintf("--- %s ---\n%s\n", docFile, snippet)

		for _, cmd := range testCmds {
			if bytes.Contains(content, []byte(cmd)) {
				if len(collectedContent) > 3000 {
					collectedContent = collectedContent[:3000]
				}
				return true, fmt.Sprintf("Test command documented in %s: %s", docFile, cmd), collectedContent
			}
		}
	}

	if lang == checker.LanguageTypeScript || lang == checker.LanguageUnknown {
		if content, err := fs.ReadFile(repo, "package.json"); err == nil {
			if bytes.Contains(content, []byte(`"test"`)) {
				if len(collectedContent) > 3000 {
					collectedContent = collectedContent[:3000]
				}
				return true, "Test script found in package.json", collectedContent
			}
		}
	}

	if content, err := fs.ReadFile(repo, "Makefile"); err == nil {
		if bytes.Contains(content, []byte("test:")) {
			if len(collectedContent) > 3000 {
				collectedContent = collectedContent[:3000]
			}
			return true, "Found Makefile with test target", collectedContent
		}
	}

	if len(collectedContent) > 3000 {
		collectedContent = collectedContent[:3000]
	}
	return false, "No test command documentation found", collectedContent
}

func testCommandsForLang(lang checker.Language) []string {
	switch lang {
	case checker.LanguageGo:
		return []string{"go test", "make test"}
	case checker.LanguageTypeScript:
		return []string{"npm test", "npm run test", "yarn test"}
	case checker.LanguageJava:
		return []string{"mvn test", "./gradlew test"}
	default:
		return []string{"go test", "make test", "npm test", "npm run test", "yarn test", "mvn test", "./gradlew test"}
	}
}
