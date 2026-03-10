package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type LocalServicesSetupChecker struct{}

func (c *LocalServicesSetupChecker) ID() checker.CheckerID  { return "local_services_setup" }
func (c *LocalServicesSetupChecker) Pillar() checker.Pillar { return checker.PillarDevEnvironment }
func (c *LocalServicesSetupChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *LocalServicesSetupChecker) Name() string           { return "Local Services Setup" }
func (c *LocalServicesSetupChecker) Description() string {
	return "Checks that a docker-compose or compose file is present for running local services"
}
func (c *LocalServicesSetupChecker) Suggestion() string {
	return "Add a docker-compose.yml to define local service dependencies (databases, message queues, etc.)"
}

func (c *LocalServicesSetupChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}
	found, path := checker.FileExistsAny(repo, candidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No docker-compose or compose file found"
	}
	return result, nil
}
