package observability

import (
	"context"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type HealthChecksChecker struct{}

func (c *HealthChecksChecker) ID() checker.CheckerID  { return "health_checks" }
func (c *HealthChecksChecker) Pillar() checker.Pillar { return checker.PillarVerification }
func (c *HealthChecksChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *HealthChecksChecker) Name() string           { return "Health Checks" }
func (c *HealthChecksChecker) Description() string {
	return "Checks that health check endpoints are defined in source files"
}
func (c *HealthChecksChecker) Suggestion() string {
	return "Add a health check endpoint (e.g., GET /healthz, /health, /ready, /ping) to your service"
}

var healthPatterns = []string{"/health", "/healthz", "/ready", "/ping", "/live"}

func (c *HealthChecksChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	found := false
	evidence := ""

	_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".ts") &&
			!strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".java") &&
			!strings.HasSuffix(path, ".py") {
			return nil
		}
		data, readErr := fs.ReadFile(repo, path)
		if readErr != nil {
			return nil
		}
		content := string(data)
		for _, pattern := range healthPatterns {
			if strings.Contains(content, pattern) {
				found = true
				evidence = path + " contains health route " + pattern
				return fs.SkipAll
			}
		}
		return nil
	})

	result.Passed = found
	if found {
		result.Evidence = "Found health check: " + evidence
	} else {
		result.Evidence = "No health check endpoints found in source files"
	}
	return result, nil
}
