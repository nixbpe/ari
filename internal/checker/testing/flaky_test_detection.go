package checktesting

import (
	"bytes"
	"context"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type FlakyTestDetectionChecker struct{}

func (c *FlakyTestDetectionChecker) ID() checker.CheckerID  { return "flaky_test_detection" }
func (c *FlakyTestDetectionChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *FlakyTestDetectionChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *FlakyTestDetectionChecker) Name() string           { return "Flaky Test Detection" }
func (c *FlakyTestDetectionChecker) Description() string {
	return "Checks that flaky test handling (retry/rerun) is configured in CI or test frameworks"
}
func (c *FlakyTestDetectionChecker) Suggestion() string {
	return "Add flaky test handling. For Go: use --count flag for stress testing. For TS: add retries to Playwright config"
}

var retryKeywords = [][]byte{
	[]byte("retry"),
	[]byte("rerun"),
	[]byte("flaky"),
	[]byte("--count="),
}

func (c *FlakyTestDetectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	workflows, _ := fs.Glob(repo, ".github/workflows/*.yml")
	yamlFiles, _ := fs.Glob(repo, ".github/workflows/*.yaml")
	workflows = append(workflows, yamlFiles...)

	for _, wf := range workflows {
		content, err := fs.ReadFile(repo, wf)
		if err != nil {
			continue
		}
		for _, kw := range retryKeywords {
			if bytes.Contains(content, kw) {
				result.Passed = true
				result.Evidence = "Found test retry configuration in CI workflow"
				return result, nil
			}
		}
	}

	for _, playwrightCfg := range []string{"playwright.config.ts", "playwright.config.js"} {
		content, err := fs.ReadFile(repo, playwrightCfg)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("retries:")) || bytes.Contains(content, []byte("--retry")) {
			result.Passed = true
			result.Evidence = "Found test retry configuration in CI workflow"
			return result, nil
		}
	}

	for _, jestCfg := range []string{"jest.config.ts", "jest.config.js", "jest.config.json"} {
		content, err := fs.ReadFile(repo, jestCfg)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("--retry")) {
			result.Passed = true
			result.Evidence = "Found test retry configuration in CI workflow"
			return result, nil
		}
	}

	if content, err := fs.ReadFile(repo, "pom.xml"); err == nil {
		if bytes.Contains(content, []byte("rerunFailingTestsCount")) {
			result.Passed = true
			result.Evidence = "Found test retry configuration in CI workflow"
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No flaky test handling configured"
	return result, nil
}
