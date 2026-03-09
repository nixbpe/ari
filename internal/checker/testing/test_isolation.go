package checktesting

import (
	"bytes"
	"context"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type TestIsolationChecker struct{}

func (c *TestIsolationChecker) ID() checker.CheckerID  { return "test_isolation" }
func (c *TestIsolationChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *TestIsolationChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *TestIsolationChecker) Name() string           { return "Test Isolation" }
func (c *TestIsolationChecker) Description() string {
	return "Checks that tests are configured for isolation and parallel execution"
}
func (c *TestIsolationChecker) Suggestion() string {
	return "Add test isolation. For Go: add t.Parallel() to test functions and use -race flag in CI"
}

func (c *TestIsolationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	switch lang {
	case checker.LanguageGo:
		found := false
		fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, "_test.go") {
				return nil
			}
			content, readErr := fs.ReadFile(repo, path)
			if readErr != nil {
				return nil
			}
			if bytes.Contains(content, []byte("t.Parallel()")) {
				found = true
			}
			return nil
		})
		if found {
			result.Passed = true
			result.Evidence = "Found t.Parallel() usage in test files — test isolation enabled"
		} else {
			result.Passed = false
			result.Evidence = "No test isolation configuration found"
		}

	case checker.LanguageTypeScript:
		if content, err := fs.ReadFile(repo, "jest.config.js"); err == nil {
			if bytes.Contains(content, []byte("maxWorkers")) {
				result.Passed = true
				result.Evidence = "Found maxWorkers in jest.config.js — test isolation configured"
				return result, nil
			}
		}
		if content, err := fs.ReadFile(repo, "jest.config.ts"); err == nil {
			if bytes.Contains(content, []byte("maxWorkers")) {
				result.Passed = true
				result.Evidence = "Found maxWorkers in jest.config.ts — test isolation configured"
				return result, nil
			}
		}
		if content, err := fs.ReadFile(repo, "playwright.config.ts"); err == nil {
			if bytes.Contains(content, []byte("--parallel")) || bytes.Contains(content, []byte("workers")) {
				result.Passed = true
				result.Evidence = "Found parallel configuration in playwright.config.ts"
				return result, nil
			}
		}
		result.Passed = false
		result.Evidence = "No test isolation configuration found"

	case checker.LanguageJava:
		found := false
		fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".java") {
				return nil
			}
			content, readErr := fs.ReadFile(repo, path)
			if readErr != nil {
				return nil
			}
			if bytes.Contains(content, []byte("@Execution(CONCURRENT)")) {
				found = true
			}
			return nil
		})
		if !found {
			if content, err := fs.ReadFile(repo, "pom.xml"); err == nil {
				if bytes.Contains(content, []byte("parallel")) {
					found = true
				}
			}
		}
		if found {
			result.Passed = true
			result.Evidence = "Found concurrent test execution configuration"
		} else {
			result.Passed = false
			result.Evidence = "No test isolation configuration found"
		}

	default:
		result.Passed = false
		result.Evidence = "No test isolation configuration found"
	}

	return result, nil
}
