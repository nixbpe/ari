package observability

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DistributedTracingChecker struct{}

func (c *DistributedTracingChecker) ID() checker.CheckerID  { return "distributed_tracing" }
func (c *DistributedTracingChecker) Pillar() checker.Pillar { return checker.PillarVerification }
func (c *DistributedTracingChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *DistributedTracingChecker) Name() string           { return "Distributed Tracing" }
func (c *DistributedTracingChecker) Description() string {
	return "Checks for distributed tracing instrumentation in dependencies"
}
func (c *DistributedTracingChecker) Suggestion() string {
	return "Add distributed tracing (e.g., OpenTelemetry, Jaeger, Zipkin) to observe request flows across services"
}

func (c *DistributedTracingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	packages := []string{
		"go.opentelemetry.io",
		"opentelemetry",
		"@opentelemetry/",
		"jaeger",
		"zipkin",
	}

	found, evidence := checker.DepFileContains(repo, lang, packages)
	result.Passed = found
	if found {
		result.Evidence = "Found distributed tracing: " + evidence
	} else {
		result.Evidence = "No distributed tracing library found in dependencies"
	}
	return result, nil
}
