package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DevcontainerChecker struct{}

func (c *DevcontainerChecker) ID() checker.CheckerID  { return "devcontainer" }
func (c *DevcontainerChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *DevcontainerChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *DevcontainerChecker) Name() string           { return "Dev Container" }
func (c *DevcontainerChecker) Description() string {
	return "Checks that a devcontainer configuration is present for reproducible development environments"
}
func (c *DevcontainerChecker) Suggestion() string {
	return "Add a .devcontainer/devcontainer.json to enable reproducible development environments via VS Code or GitHub Codespaces"
}

func (c *DevcontainerChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{".devcontainer/devcontainer.json", ".devcontainer.json"}
	found, path := checker.FileExistsAny(repo, candidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No devcontainer configuration found"
	}
	return result, nil
}
