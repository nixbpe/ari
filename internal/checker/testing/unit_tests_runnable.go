package checktesting

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type UnitTestsRunnableChecker struct{}

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
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	testCmds := testCommandsForLang(lang)
	docFiles := []string{"README.md", "CONTRIBUTING.md", "Makefile", "CLAUDE.md", "AGENTS.md"}

	for _, docFile := range docFiles {
		content, err := fs.ReadFile(repo, docFile)
		if err != nil {
			continue
		}
		for _, cmd := range testCmds {
			if bytes.Contains(content, []byte(cmd)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Test command documented in %s: %s", docFile, cmd)
				return result, nil
			}
		}
	}

	if lang == checker.LanguageTypeScript || lang == checker.LanguageUnknown {
		if content, err := fs.ReadFile(repo, "package.json"); err == nil {
			if bytes.Contains(content, []byte(`"test"`)) {
				result.Passed = true
				result.Evidence = "Test script found in package.json"
				return result, nil
			}
		}
	}

	if content, err := fs.ReadFile(repo, "Makefile"); err == nil {
		if bytes.Contains(content, []byte("test:")) {
			result.Passed = true
			result.Evidence = "Found Makefile with test target"
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No test command documentation found"
	return result, nil
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
