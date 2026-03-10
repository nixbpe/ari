package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type EnvTemplateChecker struct{}

func (c *EnvTemplateChecker) ID() checker.CheckerID  { return "env_template" }
func (c *EnvTemplateChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *EnvTemplateChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *EnvTemplateChecker) Name() string           { return "Environment Template" }
func (c *EnvTemplateChecker) Description() string {
	return "Checks that a .env.example or similar template documents required environment variables"
}
func (c *EnvTemplateChecker) Suggestion() string {
	return "Create a .env.example file documenting required environment variables"
}

func (c *EnvTemplateChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{".env.example", ".env.template", ".env.sample", ".env.local.example"}
	found, path := checker.FileExistsAny(repo, candidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No .env.example or similar template found"
	}
	return result, nil
}
