package docs

import (
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
	"time"

	"github.com/nixbpe/ari/internal/checker"
)

// DocumentationFreshnessChecker checks whether documentation files have been
// updated recently using git history. It uses a GitRunner field for testability.
type DocumentationFreshnessChecker struct {
	// RepoPath is the filesystem path passed to "git -C <RepoPath>" when using the default runner.
	RepoPath string

	// GitRunner runs git commands and returns stdout. Defaults to exec.Command("git", ...).Output().
	// Override in tests to avoid requiring an actual git repository.
	GitRunner func(args ...string) ([]byte, error)
}

func (c *DocumentationFreshnessChecker) ID() checker.CheckerID  { return "documentation_freshness" }
func (c *DocumentationFreshnessChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *DocumentationFreshnessChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *DocumentationFreshnessChecker) Name() string           { return "Documentation Freshness" }
func (c *DocumentationFreshnessChecker) Description() string {
	return "Checks that documentation files (README, CONTRIBUTING, AGENTS) have been updated within 180 days"
}
func (c *DocumentationFreshnessChecker) Suggestion() string {
	return "Update documentation — it hasn't been updated in over 180 days"
}

func (c *DocumentationFreshnessChecker) gitRunner() func(args ...string) ([]byte, error) {
	if c.GitRunner != nil {
		return c.GitRunner
	}
	repoPath := c.RepoPath
	if repoPath == "" {
		repoPath = "."
	}
	return func(args ...string) ([]byte, error) {
		allArgs := append([]string{"-C", repoPath}, args...)
		return exec.Command("git", allArgs...).Output()
	}
}

const freshnessThresholdDays = 180
const gitDateFormat = "2006-01-02 15:04:05 -0700"

func (c *DocumentationFreshnessChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	run := c.gitRunner()
	out, err := run("log", "--format=%ci", "-1", "--", "README.md", "CONTRIBUTING.md", "AGENTS.md", "CLAUDE.md")
	if err != nil {
		result.Skipped = true
		result.SkipReason = "not a git repository"
		return result, nil
	}

	dateStr := strings.TrimSpace(string(out))
	if dateStr == "" {
		result.Passed = false
		result.Evidence = "Documentation files not tracked in git"
		return result, nil
	}

	// git outputs dates like "2025-01-15 10:30:00 +0000"; parse only the first line
	firstLine := strings.SplitN(dateStr, "\n", 2)[0]
	firstLine = strings.TrimSpace(firstLine)

	t, parseErr := time.Parse(gitDateFormat, firstLine)
	if parseErr != nil {
		result.Passed = false
		result.Evidence = fmt.Sprintf("Could not parse git date: %s", firstLine)
		return result, nil
	}

	days := int(time.Since(t).Hours() / 24)
	if days <= freshnessThresholdDays {
		result.Passed = true
		result.Evidence = fmt.Sprintf("Documentation last updated %d days ago — documentation is fresh", days)
	} else {
		result.Passed = false
		result.Evidence = "Documentation not updated in over 180 days"
	}
	return result, nil
}
