package observability

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type StructuredLoggingChecker struct{}

func (c *StructuredLoggingChecker) ID() checker.CheckerID  { return "structured_logging" }
func (c *StructuredLoggingChecker) Pillar() checker.Pillar { return checker.PillarObservability }
func (c *StructuredLoggingChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *StructuredLoggingChecker) Name() string           { return "Structured Logging" }
func (c *StructuredLoggingChecker) Description() string {
	return "Checks for a structured logging library in dependencies"
}
func (c *StructuredLoggingChecker) Suggestion() string {
	return "Add a structured logging library (e.g., zap, zerolog, slog for Go; winston/pino for TypeScript; logback/log4j2 for Java)"
}

func (c *StructuredLoggingChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	packages := []string{
		"go.uber.org/zap",
		"rs/zerolog",
		"log/slog",
		"winston",
		"pino",
		"bunyan",
		"logback",
		"log4j2",
	}

	found, evidence := checker.DepFileContains(repo, lang, packages)
	result.Passed = found
	if found {
		result.Evidence = "Found structured logging: " + evidence
	} else {
		result.Evidence = "No structured logging library found in dependencies"
	}
	return result, nil
}
