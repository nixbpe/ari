package observability

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type AlertingConfiguredChecker struct{}

func (c *AlertingConfiguredChecker) ID() checker.CheckerID  { return "alerting_configured" }
func (c *AlertingConfiguredChecker) Pillar() checker.Pillar { return checker.PillarObservability }
func (c *AlertingConfiguredChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *AlertingConfiguredChecker) Name() string           { return "Alerting Configured" }
func (c *AlertingConfiguredChecker) Description() string {
	return "Checks that alerting configuration files are present in the repository"
}
func (c *AlertingConfiguredChecker) Suggestion() string {
	return "Add alerting configuration (e.g., alertmanager.yml, Prometheus alert rules in monitoring/rules/, or .sentry.properties)"
}

func (c *AlertingConfiguredChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	configPaths := []string{
		"alertmanager.yml",
		"alertmanager.yaml",
		".sentry.properties",
		"datadog.yaml",
	}
	found, evidence := checker.FileExistsAny(repo, configPaths)
	if found {
		result.Passed = true
		result.Evidence = "Found alerting configuration: " + evidence
		return result, nil
	}

	alertDirs := []string{"monitoring/rules", "alerts"}
	for _, dir := range alertDirs {
		if _, err := fs.ReadDir(repo, dir); err == nil {
			result.Passed = true
			result.Evidence = "Found alerting configuration directory: " + dir
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No alerting configuration found (checked alertmanager.yml, monitoring/rules/, alerts/, .sentry.properties, datadog.yaml)"
	return result, nil
}
