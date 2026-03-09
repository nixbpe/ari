package build

import (
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
	"time"

	"github.com/bbik/ari/internal/checker"
)

type DeploymentFrequencyChecker struct {
	RepoPath  string
	GitRunner func(args ...string) ([]byte, error)
}

func (c *DeploymentFrequencyChecker) ID() checker.CheckerID  { return "deployment_frequency" }
func (c *DeploymentFrequencyChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *DeploymentFrequencyChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *DeploymentFrequencyChecker) Name() string           { return "Deployment Frequency" }
func (c *DeploymentFrequencyChecker) Description() string {
	return "Checks how frequently the project deploys by examining git tag history"
}
func (c *DeploymentFrequencyChecker) Suggestion() string {
	return "Increase deployment frequency. Aim for at least one release per 90 days"
}

func (c *DeploymentFrequencyChecker) runner() func(args ...string) ([]byte, error) {
	if c.GitRunner != nil {
		return c.GitRunner
	}
	return func(args ...string) ([]byte, error) {
		return exec.Command("git", args...).Output()
	}
}

func (c *DeploymentFrequencyChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	run := c.runner()

	gitArgs := []string{"tag", "--sort=-creatordate", "--format=%(refname:short)|%(creatordate:iso8601)", "-l"}
	if c.RepoPath != "" {
		gitArgs = append([]string{"-C", c.RepoPath}, gitArgs...)
	}

	out, err := run(gitArgs...)
	if err != nil {
		result.Skipped = true
		result.SkipReason = "not a git repository or no tags"
		return result, nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var firstValid string
	var firstTag string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			firstTag = strings.TrimSpace(parts[0])
			firstValid = strings.TrimSpace(parts[1])
			break
		}
	}

	if firstValid == "" {
		result.Passed = false
		result.Evidence = "No releases found"
		return result, nil
	}

	tagDate, parseErr := time.Parse("2006-01-02 15:04:05 -0700", firstValid)
	if parseErr != nil {
		tagDate, parseErr = time.Parse("2006-01-02 15:04:05 +0000", firstValid)
	}
	if parseErr != nil {
		result.Passed = false
		result.Evidence = fmt.Sprintf("Found tag %s but could not parse date: %s", firstTag, firstValid)
		return result, nil
	}

	daysSince := int(time.Since(tagDate).Hours() / 24)
	if daysSince <= 90 {
		result.Passed = true
		result.Evidence = fmt.Sprintf("Last release: %s (%d days ago) — regular deployment cadence", firstTag, daysSince)
	} else {
		result.Passed = false
		result.Evidence = fmt.Sprintf("No releases in last 90 days (last: %s, %d days ago)", firstTag, daysSince)
	}
	return result, nil
}
