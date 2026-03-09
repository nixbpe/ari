package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type BuildPerformanceTrackingChecker struct{}

func (c *BuildPerformanceTrackingChecker) ID() checker.CheckerID  { return "build_performance_tracking" }
func (c *BuildPerformanceTrackingChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *BuildPerformanceTrackingChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *BuildPerformanceTrackingChecker) Name() string           { return "Build Performance Tracking" }
func (c *BuildPerformanceTrackingChecker) Description() string {
	return "Checks that build caching and performance tracking tools are configured"
}
func (c *BuildPerformanceTrackingChecker) Suggestion() string {
	return "Add build caching. For Go: use GitHub Actions cache for Go modules. For JS: use Turborepo"
}

func (c *BuildPerformanceTrackingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if _, err := fs.Stat(repo, "turbo.json"); err == nil {
		result.Passed = true
		result.Evidence = "Found turbo.json — build caching enabled"
		return result, nil
	}

	if _, err := fs.Stat(repo, "nx.json"); err == nil {
		result.Passed = true
		result.Evidence = "Found nx.json — build caching enabled"
		return result, nil
	}

	for _, f := range []string{"WORKSPACE", "BUILD", ".bazelrc"} {
		if _, err := fs.Stat(repo, f); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — Bazel build system configured", f)
			return result, nil
		}
	}

	workflows, _ := fs.Glob(repo, ".github/workflows/*.yml")
	yamlFiles, _ := fs.Glob(repo, ".github/workflows/*.yaml")
	workflows = append(workflows, yamlFiles...)

	for _, wf := range workflows {
		content, err := fs.ReadFile(repo, wf)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("cache")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found cache step in CI workflow %s — build caching enabled", wf)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No build performance tracking found"
	return result, nil
}
