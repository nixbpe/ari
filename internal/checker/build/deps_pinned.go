package build

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DepsPinnedChecker struct{}

func (c *DepsPinnedChecker) ID() checker.CheckerID  { return "deps_pinned" }
func (c *DepsPinnedChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *DepsPinnedChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *DepsPinnedChecker) Name() string           { return "Dependencies Pinned" }
func (c *DepsPinnedChecker) Description() string {
	return "Checks for a dependency lock file to ensure reproducible builds"
}
func (c *DepsPinnedChecker) Suggestion() string {
	return "Pin dependencies. For Go: commit go.sum. For npm: npm install and commit package-lock.json"
}

func (c *DepsPinnedChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	lockFiles := lockFilesForLang(lang)
	for _, lf := range lockFiles {
		if _, err := fs.Stat(repo, lf.file); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — %s", lf.file, lf.description)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No dependency lock file found"
	return result, nil
}

type lockFileEntry struct {
	file        string
	description string
}

func lockFilesForLang(lang checker.Language) []lockFileEntry {
	switch lang {
	case checker.LanguageGo:
		return []lockFileEntry{
			{"go.sum", "dependencies pinned with exact versions"},
		}
	case checker.LanguageTypeScript:
		return []lockFileEntry{
			{"package-lock.json", "dependencies pinned via npm"},
			{"yarn.lock", "dependencies pinned via yarn"},
			{"pnpm-lock.yaml", "dependencies pinned via pnpm"},
			{"bun.lockb", "dependencies pinned via bun"},
		}
	case checker.LanguageJava:
		return []lockFileEntry{
			{"gradle.lockfile", "dependencies pinned via Gradle"},
			{"pom.xml", "dependencies managed via Maven"},
		}
	default:
		return []lockFileEntry{
			{"go.sum", "dependencies pinned with exact versions"},
			{"package-lock.json", "dependencies pinned via npm"},
			{"yarn.lock", "dependencies pinned via yarn"},
			{"pnpm-lock.yaml", "dependencies pinned via pnpm"},
		}
	}
}
