package checktesting

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

// IntegrationTestsExistChecker checks whether integration or end-to-end tests exist in the repository.
type IntegrationTestsExistChecker struct{}

func (c *IntegrationTestsExistChecker) ID() checker.CheckerID  { return "integration_tests_exist" }
func (c *IntegrationTestsExistChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *IntegrationTestsExistChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *IntegrationTestsExistChecker) Name() string           { return "Integration Tests Exist" }
func (c *IntegrationTestsExistChecker) Description() string {
	return "Checks for integration or end-to-end tests (directories, config files, or build tags)"
}
func (c *IntegrationTestsExistChecker) Suggestion() string {
	return "Add integration tests. Create an e2e/ or integration/ directory with end-to-end tests"
}

func (c *IntegrationTestsExistChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	// 1. Check for well-known integration/e2e directories
	integrationDirs := []string{
		"integration",
		"e2e",
		"acceptance",
		"test/integration",
		"tests/e2e",
	}
	for _, dir := range integrationDirs {
		info, err := fs.Stat(repo, dir)
		if err == nil && info.IsDir() {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found integration tests in %s/ directory", dir)
			return result, nil
		}
	}

	// 2. Check TypeScript e2e framework config files
	tsConfigs := []string{
		"playwright.config.ts",
		"cypress.config.ts",
		"cypress.config.js",
	}
	for _, cfg := range tsConfigs {
		if _, err := fs.Stat(repo, cfg); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — integration test framework configured", cfg)
			return result, nil
		}
	}

	// 3. Check Java src/integrationTest/ directory
	if info, err := fs.Stat(repo, "src/integrationTest"); err == nil && info.IsDir() {
		result.Passed = true
		result.Evidence = "Found integration tests in src/integrationTest/ directory"
		return result, nil
	}

	// 4. Walk files: Go build tags, TestIntegration functions, Java IT files
	_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || result.Passed {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// Java: files ending in IT.java
		if strings.HasSuffix(path, "IT.java") {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found Java integration test file %s", path)
			return fs.SkipAll
		}

		// Go: check file content for build tag or TestIntegration function
		if strings.HasSuffix(path, ".go") {
			content, readErr := fs.ReadFile(repo, path)
			if readErr != nil {
				return nil
			}
			if bytes.Contains(content, []byte("//go:build integration")) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found //go:build integration tag in %s", path)
				return fs.SkipAll
			}
			if bytes.Contains(content, []byte("func TestIntegration")) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found TestIntegration function in %s", path)
				return fs.SkipAll
			}
		}

		return nil
	})

	if !result.Passed {
		result.Evidence = "No integration tests found"
	}
	return result, nil
}
