package observability

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type MetricsCollectionChecker struct{}

func (c *MetricsCollectionChecker) ID() checker.CheckerID  { return "metrics_collection" }
func (c *MetricsCollectionChecker) Pillar() checker.Pillar { return checker.PillarVerification }
func (c *MetricsCollectionChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *MetricsCollectionChecker) Name() string           { return "Metrics Collection" }
func (c *MetricsCollectionChecker) Description() string {
	return "Checks for metrics collection library in dependencies"
}
func (c *MetricsCollectionChecker) Suggestion() string {
	return "Add a metrics library (e.g., prometheus/client_golang for Go; prom-client for Node.js; Micrometer for Java) to expose application metrics"
}

func (c *MetricsCollectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	packages := []string{
		"prometheus/client_golang",
		"prom-client",
		"micrometer",
		"statsd",
	}

	found, evidence := checker.DepFileContains(repo, lang, packages)
	result.Passed = found
	if found {
		result.Evidence = "Found metrics collection: " + evidence
	} else {
		result.Evidence = "No metrics collection library found in dependencies"
	}
	return result, nil
}
