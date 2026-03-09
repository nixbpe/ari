package checktesting

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type TestCoverageThresholdsChecker struct{}

func (c *TestCoverageThresholdsChecker) ID() checker.CheckerID  { return "test_coverage_thresholds" }
func (c *TestCoverageThresholdsChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *TestCoverageThresholdsChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *TestCoverageThresholdsChecker) Name() string           { return "Test Coverage Thresholds" }
func (c *TestCoverageThresholdsChecker) Description() string {
	return "Checks for test coverage threshold configuration in CI, build files, or test framework config"
}
func (c *TestCoverageThresholdsChecker) Suggestion() string {
	return "Add coverage thresholds. For Go: add -coverprofile to CI. For TS: add coverageThreshold to jest config"
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
		checkers = []func(fs.FS) (bool, string){checkGoCoverage}
	case checker.LanguageTypeScript:
		checkers = []func(fs.FS) (bool, string){checkTSCoverage}
	case checker.LanguageJava:
		checkers = []func(fs.FS) (bool, string){checkJavaCoverage}
	default:
		checkers = []func(fs.FS) (bool, string){checkGoCoverage, checkTSCoverage, checkJavaCoverage}
	}

	for _, fn := range checkers {
		if found, evidence := fn(repo); found {
			result.Passed = true
			result.Evidence = evidence
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No test coverage threshold configuration found"
	return result, nil
}

func checkGoCoverage(repo fs.FS) (bool, string) {
	// Check CI workflow files for coverage flags
	for _, pattern := range []string{".github/workflows/*.yml", ".github/workflows/*.yaml"} {
		workflows, _ := fs.Glob(repo, pattern)
		for _, wf := range workflows {
			content, err := fs.ReadFile(repo, wf)
			if err != nil {
				continue
			}
			if bytes.Contains(content, []byte("-coverprofile")) {
				return true, fmt.Sprintf("Found coverage configuration in CI workflow %s", wf)
			}
			if bytes.Contains(content, []byte("-cover ")) || bytes.Contains(content, []byte("-cover\n")) {
				return true, fmt.Sprintf("Found coverage configuration in CI workflow %s", wf)
			}
		}
	}

	// Check Makefile for coverage target
	if content, err := fs.ReadFile(repo, "Makefile"); err == nil {
		if bytes.Contains(content, []byte("coverage")) {
			return true, "Found coverage target in Makefile"
		}
	}

	return false, ""
}

func checkTSCoverage(repo fs.FS) (bool, string) {
	// Check package.json for jest --coverage in scripts
	if content, err := fs.ReadFile(repo, "package.json"); err == nil {
		if bytes.Contains(content, []byte("jest --coverage")) || bytes.Contains(content, []byte("--coverage")) {
			return true, "Found coverage configuration in package.json"
		}
	}

	// Check jest config files for coverageThreshold
	for _, cfg := range []string{"jest.config.js", "jest.config.ts", "jest.config.mjs"} {
		if content, err := fs.ReadFile(repo, cfg); err == nil {
			if bytes.Contains(content, []byte("coverageThreshold")) {
				return true, fmt.Sprintf("Found coverageThreshold in %s", cfg)
			}
			if bytes.Contains(content, []byte("coverage")) {
				return true, fmt.Sprintf("Found coverage configuration in %s", cfg)
			}
		}
	}

	return false, ""
}

func checkJavaCoverage(repo fs.FS) (bool, string) {
	// Check pom.xml for JaCoCo plugin
	if content, err := fs.ReadFile(repo, "pom.xml"); err == nil {
		if bytes.Contains(content, []byte("jacoco")) {
			return true, "Found JaCoCo plugin in pom.xml"
		}
	}

	// Check build.gradle / build.gradle.kts for JaCoCo
	for _, bf := range []string{"build.gradle", "build.gradle.kts"} {
		if content, err := fs.ReadFile(repo, bf); err == nil {
			if bytes.Contains(content, []byte("jacocoTestCoverageVerification")) || bytes.Contains(content, []byte("jacoco")) {
				return true, fmt.Sprintf("Found JaCoCo configuration in %s", bf)
			}
		}
	}

	return false, ""
}
