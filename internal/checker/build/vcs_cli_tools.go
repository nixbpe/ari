package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type VCSCliToolsChecker struct{}

func (c *VCSCliToolsChecker) ID() checker.CheckerID  { return "vcs_cli_tools" }
func (c *VCSCliToolsChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *VCSCliToolsChecker) Level() checker.Level   { return checker.LevelAutonomous }
func (c *VCSCliToolsChecker) Name() string           { return "VCS CLI Tools" }
func (c *VCSCliToolsChecker) Description() string {
	return "Checks for VCS CLI tooling such as GitHub CLI (gh) configured in the project"
}
func (c *VCSCliToolsChecker) Suggestion() string {
	return "Install GitHub CLI: brew install gh and document its usage in AGENTS.md"
}

func (c *VCSCliToolsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if _, err := fs.Stat(repo, ".github"); err == nil {
		result.Passed = true
		result.Evidence = "Found .github/ directory — GitHub CLI workflow configured"
		return result, nil
	}

	candidates := []string{"Makefile", "README.md", "CONTRIBUTING.md", "CLAUDE.md", "AGENTS.md"}
	for _, file := range candidates {
		data, err := fs.ReadFile(repo, file)
		if err != nil {
			continue
		}
		if bytes.Contains(data, []byte("gh ")) || bytes.Contains(data, []byte("`gh`")) || bytes.Contains(data, []byte("github.com/cli/cli")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found 'gh' mentioned in %s — GitHub CLI usage documented", file)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No VCS CLI tool usage found"
	return result, nil
}
