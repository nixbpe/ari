package checktesting

import (
	"bytes"
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type TestPerformanceTrackingChecker struct{}

func (c *TestPerformanceTrackingChecker) ID() checker.CheckerID  { return "test_performance_tracking" }
func (c *TestPerformanceTrackingChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *TestPerformanceTrackingChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *TestPerformanceTrackingChecker) Name() string           { return "Test Performance Tracking" }
func (c *TestPerformanceTrackingChecker) Description() string {
	return "Checks that test performance tracking (benchmarks) is configured"
}
func (c *TestPerformanceTrackingChecker) Suggestion() string {
	return "Add test performance tracking. For Go: add go test -bench to CI. For TS: add performance budgets"
}

func (c *TestPerformanceTrackingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		if bytes.Contains(content, []byte("go test -bench")) || bytes.Contains(content, []byte("go test -run=^$ -bench")) {
			result.Passed = true
			result.Evidence = "Found benchmark CI workflow — test performance tracking enabled"
			return result, nil
		}
		if bytes.Contains(content, []byte("jest --verbose")) || bytes.Contains(content, []byte("performance budget")) {
			result.Passed = true
			result.Evidence = "Found benchmark CI workflow — test performance tracking enabled"
			return result, nil
		}
	}

	if _, err := fs.Stat(repo, "benchmarks"); err == nil {
		result.Passed = true
		result.Evidence = "Found benchmark CI workflow — test performance tracking enabled"
		return result, nil
	}

	benchFiles, _ := fs.Glob(repo, "*_bench_test.go")
	benchFiles2, _ := fs.Glob(repo, "bench_test.go")
	benchFiles = append(benchFiles, benchFiles2...)

	benchFilesNested, _ := fs.Glob(repo, "**/*_bench_test.go")
	benchFiles = append(benchFiles, benchFilesNested...)

	if len(benchFiles) > 0 {
		result.Passed = true
		result.Evidence = "Found benchmark CI workflow — test performance tracking enabled"
		return result, nil
	}

	tsBenchFiles, _ := fs.Glob(repo, "*.bench.ts")
	tsBenchFilesNested, _ := fs.Glob(repo, "**/*.bench.ts")
	tsBenchFiles = append(tsBenchFiles, tsBenchFilesNested...)

	if len(tsBenchFiles) > 0 {
		result.Passed = true
		result.Evidence = "Found benchmark CI workflow — test performance tracking enabled"
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No test performance tracking found"
	return result, nil
}
