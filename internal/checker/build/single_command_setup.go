package build

import (
	"bytes"
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type SingleCommandSetupChecker struct{}

func (c *SingleCommandSetupChecker) ID() checker.CheckerID  { return "single_command_setup" }
func (c *SingleCommandSetupChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *SingleCommandSetupChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *SingleCommandSetupChecker) Name() string           { return "Single Command Setup" }
func (c *SingleCommandSetupChecker) Description() string {
	return "Checks for a single command to set up the development environment"
}
func (c *SingleCommandSetupChecker) Suggestion() string {
	return "Add a single setup command. Create a Makefile with make setup or a script/setup script"
}

func (c *SingleCommandSetupChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if content, err := fs.ReadFile(repo, "Makefile"); err == nil {
		if bytes.Contains(content, []byte("setup:")) || bytes.Contains(content, []byte("install:")) {
			result.Passed = true
			result.Evidence = "Found Makefile with setup target — single command setup available"
			return result, nil
		}
	}

	simpleFiles := []string{
		"script/setup",
		"script/bootstrap",
		"docker-compose.yml",
		".devcontainer/devcontainer.json",
	}

	for _, f := range simpleFiles {
		if _, err := fs.Stat(repo, f); err == nil {
			result.Passed = true
			result.Evidence = "Found " + f + " — single command setup available"
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No single command setup found"
	return result, nil
}
