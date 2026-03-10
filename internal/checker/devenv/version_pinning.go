package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type VersionPinningChecker struct{}

func (c *VersionPinningChecker) ID() checker.CheckerID  { return "version_pinning" }
func (c *VersionPinningChecker) Pillar() checker.Pillar { return checker.PillarDevEnvironment }
func (c *VersionPinningChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *VersionPinningChecker) Name() string           { return "Runtime Version Pinning" }
func (c *VersionPinningChecker) Description() string {
	return "Checks that the runtime/SDK version is pinned via .nvmrc, .python-version, .tool-versions, or similar"
}
func (c *VersionPinningChecker) Suggestion() string {
	return "Pin your runtime version using .nvmrc (Node), .python-version (Python), .tool-versions (asdf/mise), or mise.toml"
}

func (c *VersionPinningChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	candidates := []string{
		".nvmrc",
		".node-version",
		".python-version",
		".tool-versions",
		"mise.toml",
		"rust-toolchain.toml",
		".java-version",
		"global.json",
	}
	found, path := checker.FileExistsAny(repo, candidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No runtime version pinning file found"
	}
	return result, nil
}
