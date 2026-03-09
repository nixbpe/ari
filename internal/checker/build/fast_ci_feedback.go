package build

import (
	"context"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type FastCIFeedbackChecker struct{}

func (c *FastCIFeedbackChecker) ID() checker.CheckerID  { return "fast_ci_feedback" }
func (c *FastCIFeedbackChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *FastCIFeedbackChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *FastCIFeedbackChecker) Name() string           { return "Fast CI Feedback" }
func (c *FastCIFeedbackChecker) Description() string {
	return "Checks for CI/CD configuration (GitHub Actions, GitLab CI, Jenkins, CircleCI)"
}
func (c *FastCIFeedbackChecker) Suggestion() string {
	return "Set up CI/CD. For GitHub: create .github/workflows/ci.yml with test and lint steps"
}

func (c *FastCIFeedbackChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	type ciEntry struct {
		path    string
		isDir   bool
		message string
	}

	candidates := []ciEntry{
		{".github/workflows", true, "Found GitHub Actions workflows — CI configured"},
		{".gitlab-ci.yml", false, "Found .gitlab-ci.yml — CI configured"},
		{"Jenkinsfile", false, "Found Jenkinsfile — CI configured"},
		{".circleci/config.yml", false, "Found .circleci/config.yml — CI configured"},
	}

	for _, entry := range candidates {
		info, err := fs.Stat(repo, entry.path)
		if err == nil {
			if entry.isDir && !info.IsDir() {
				continue
			}
			result.Passed = true
			result.Evidence = entry.message
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No CI configuration found"
	return result, nil
}
