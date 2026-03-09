package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ReleaseNotesAutomationChecker struct{}

func (c *ReleaseNotesAutomationChecker) ID() checker.CheckerID  { return "release_notes_automation" }
func (c *ReleaseNotesAutomationChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *ReleaseNotesAutomationChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *ReleaseNotesAutomationChecker) Name() string           { return "Release Notes Automation" }
func (c *ReleaseNotesAutomationChecker) Description() string {
	return "Checks that release notes automation is configured"
}
func (c *ReleaseNotesAutomationChecker) Suggestion() string {
	return "Automate release notes. Use git-cliff, release-please, or semantic-release"
}

var releaseConfigFiles = []string{
	"CHANGELOG.md",
	"cliff.toml",
	"release-please-config.json",
	".releaserc",
	".releaserc.json",
	".releaserc.yml",
}

var releaseDevDeps = []string{"semantic-release", "standard-version"}

func (c *ReleaseNotesAutomationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	for _, f := range releaseConfigFiles {
		if _, err := fs.Stat(repo, f); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — release notes documented", f)
			return result, nil
		}
	}

	if pkgJSON, err := fs.ReadFile(repo, "package.json"); err == nil {
		for _, dep := range releaseDevDeps {
			if bytes.Contains(pkgJSON, []byte(dep)) {
				result.Passed = true
				result.Evidence = "Found " + dep + " in package.json — release notes automation configured"
				return result, nil
			}
		}
	}

	result.Passed = false
	result.Evidence = "No release notes automation found"
	return result, nil
}
