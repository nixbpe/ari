package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type SASTConfigChecker struct{}

func (c *SASTConfigChecker) ID() checker.CheckerID  { return "sast_config" }
func (c *SASTConfigChecker) Pillar() checker.Pillar { return checker.PillarSecurity }
func (c *SASTConfigChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *SASTConfigChecker) Name() string           { return "SAST Configuration" }
func (c *SASTConfigChecker) Description() string {
	return "Checks that static application security testing (SAST) is configured via config file or CI"
}
func (c *SASTConfigChecker) Suggestion() string {
	return "Add SAST tooling: configure .semgrep.yml or add CodeQL/Semgrep/SonarQube to CI workflows"
}

var sastFileCandidates = []string{
	".semgrep.yml",
	".semgrep.yaml",
	"sonar-project.properties",
	".sonarcloud.properties",
}

var sastKeywords = []string{"codeql", "semgrep", "sonarqube", "sonar"}

func (c *SASTConfigChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	fileFound, path := checker.FileExistsAny(repo, sastFileCandidates)
	if fileFound {
		result.Passed = true
		result.Evidence = "Found " + path
		return result, nil
	}

	ciFound, evidence := checker.CIWorkflowContains(repo, sastKeywords)
	result.Passed = ciFound
	if ciFound {
		result.Evidence = "SAST in CI: " + evidence
	} else {
		result.Evidence = "No SAST config or CI step found"
	}

	return result, nil
}
