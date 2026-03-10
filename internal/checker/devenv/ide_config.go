package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type IDEConfigChecker struct{}

func (c *IDEConfigChecker) ID() checker.CheckerID  { return "ide_config" }
func (c *IDEConfigChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *IDEConfigChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *IDEConfigChecker) Name() string           { return "IDE Configuration" }
func (c *IDEConfigChecker) Description() string {
	return "Checks that IDE configuration files are committed for consistent editor settings across the team"
}
func (c *IDEConfigChecker) Suggestion() string {
	return "Add .vscode/settings.json, .vscode/extensions.json, or .editorconfig to standardize IDE settings for all contributors"
}

func (c *IDEConfigChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{
		".vscode/settings.json",
		".vscode/extensions.json",
		".editorconfig",
		".idea/workspace.xml",
	}
	found, path := checker.FileExistsAny(repo, candidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No IDE configuration files found"
	}
	return result, nil
}
