package checktesting

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"

	"github.com/nixbpe/ari/internal/checker"
)

type TestCoverageThresholdsChecker struct{}

func (c *TestCoverageThresholdsChecker) ID() checker.CheckerID  { return "test_coverage_thresholds" }
func (c *TestCoverageThresholdsChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *TestCoverageThresholdsChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *TestCoverageThresholdsChecker) Name() string           { return "Test Coverage Thresholds" }
func (c *TestCoverageThresholdsChecker) Description() string {
	return "Checks for test coverage threshold configuration of ≥80% in CI, build files, or test framework config"
}
func (c *TestCoverageThresholdsChecker) Suggestion() string {
	return "Add coverage thresholds ≥80%. For Go: add .testcoverage.yml with threshold.total: 80. For TS: set coverageThreshold in jest config"
}

func (c *TestCoverageThresholdsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	var checkers []func(fs.FS) (bool, string)
	switch lang {
	case checker.LanguageGo:
		checkers = []func(fs.FS) (bool, string){checkGoCoverageThreshold}
	case checker.LanguageTypeScript:
		checkers = []func(fs.FS) (bool, string){checkTSCoverageThreshold}
	case checker.LanguageJava:
		checkers = []func(fs.FS) (bool, string){checkJavaCoverageThreshold}
	default:
		checkers = []func(fs.FS) (bool, string){checkGoCoverageThreshold, checkTSCoverageThreshold, checkJavaCoverageThreshold}
	}

	for _, fn := range checkers {
		if found, evidence := fn(repo); found {
			result.Passed = true
			result.Evidence = evidence
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No test coverage threshold ≥80% configured"
	return result, nil
}

// reTotalThreshold matches "total: 80" in .testcoverage.yml (go-test-coverage nested format).
var reTotalThreshold = regexp.MustCompile(`(?m)^\s*total\s*:\s*(\d+)`)

// reScalarThreshold matches "threshold: 80" in .testcoverage.yml (scalar format).
var reScalarThreshold = regexp.MustCompile(`(?m)^\s*threshold\s*:\s*(\d+)`)

// reCLIThreshold matches --threshold=80 or --min-coverage=80 patterns in CI scripts.
var reCLIThreshold = regexp.MustCompile(`(?:--threshold|--min-coverage|--min-coverage-overall)=(\d+)`)

func extractThresholdFromYAML(content []byte) (int, bool) {
	for _, re := range []*regexp.Regexp{reTotalThreshold, reScalarThreshold} {
		if m := re.FindSubmatch(content); len(m) > 1 {
			if v, err := strconv.Atoi(string(m[1])); err == nil {
				return v, true
			}
		}
	}
	return 0, false
}

func checkGoCoverageThreshold(repo fs.FS) (bool, string) {
	for _, pattern := range []string{".github/workflows/*.yml", ".github/workflows/*.yaml"} {
		workflows, _ := fs.Glob(repo, pattern)
		for _, wf := range workflows {
			content, err := fs.ReadFile(repo, wf)
			if err != nil {
				continue
			}
			if m := reCLIThreshold.FindSubmatch(content); len(m) > 1 {
				if v, err := strconv.Atoi(string(m[1])); err == nil && v >= 80 {
					return true, fmt.Sprintf("Found %d%% coverage threshold enforced in CI %s", v, wf)
				}
			}
			if bytes.Contains(content, []byte("go-test-coverage")) {
				cfg, err := fs.ReadFile(repo, ".testcoverage.yml")
				if err != nil {
					continue
				}
				if v, ok := extractThresholdFromYAML(cfg); ok && v >= 80 {
					return true, fmt.Sprintf("Found go-test-coverage enforcing %d%% threshold via .testcoverage.yml in CI %s", v, wf)
				}
			}
		}
	}
	return false, ""
}

// reJestThreshold matches "lines: 80" / "branches: 80" etc. inside coverageThreshold blocks.
var reJestThreshold = regexp.MustCompile(`(?:lines|branches|functions|statements)\s*:\s*(\d+)`)

func checkTSCoverageThreshold(repo fs.FS) (bool, string) {
	for _, cfg := range []string{"jest.config.js", "jest.config.ts", "jest.config.mjs", "jest.config.cjs"} {
		content, err := fs.ReadFile(repo, cfg)
		if err != nil {
			continue
		}
		if !bytes.Contains(content, []byte("coverageThreshold")) {
			continue
		}
		matches := reJestThreshold.FindAllSubmatch(content, -1)
		if len(matches) == 0 {
			continue
		}
		allAbove80 := true
		for _, m := range matches {
			if v, err := strconv.Atoi(string(m[1])); err != nil || v < 80 {
				allAbove80 = false
				break
			}
		}
		if allAbove80 {
			return true, fmt.Sprintf("Found coverageThreshold ≥80%% in %s", cfg)
		}
	}
	return false, ""
}

// reJacocoThreshold matches JaCoCo minimum values like <minimum>0.80</minimum>.
var reJacocoThreshold = regexp.MustCompile(`<minimum>(0\.\d+)</minimum>`)

func checkJavaCoverageThreshold(repo fs.FS) (bool, string) {
	for _, bf := range []string{"pom.xml", "build.gradle", "build.gradle.kts"} {
		content, err := fs.ReadFile(repo, bf)
		if err != nil {
			continue
		}
		if !bytes.Contains(content, []byte("jacoco")) {
			continue
		}
		if m := reJacocoThreshold.FindSubmatch(content); len(m) > 1 {
			if v, err := strconv.ParseFloat(string(m[1]), 64); err == nil && v >= 0.80 {
				return true, fmt.Sprintf("Found JaCoCo coverage threshold %.0f%% in %s", v*100, bf)
			}
		}
	}
	return false, ""
}
