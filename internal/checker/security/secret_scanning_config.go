package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type SecretScanningConfigChecker struct{}

func (c *SecretScanningConfigChecker) ID() checker.CheckerID  { return "secret_scanning_config" }
func (c *SecretScanningConfigChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *SecretScanningConfigChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *SecretScanningConfigChecker) Name() string           { return "Secret Scanning Config" }
func (c *SecretScanningConfigChecker) Description() string {
	return "Checks that secret scanning is configured via a config file or CI workflow step"
}
func (c *SecretScanningConfigChecker) Suggestion() string {
	return "Add secret scanning: configure .gitleaks.toml or add gitleaks/detect-secrets/trufflehog to CI workflows"
}

var secretScanningFileCandidates = []string{
	".gitleaks.toml",
	".gitleaks.yaml",
	".detect-secrets.yaml",
}

var secretScanningKeywords = []string{"gitleaks", "detect-secrets", "trufflehog"}

func (c *SecretScanningConfigChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	fileFound, path := checker.FileExistsAny(repo, secretScanningFileCandidates)
	if fileFound {
		result.Passed = true
		result.Evidence = "Found " + path
		return result, nil
	}

	ciFound, evidence := checker.CIWorkflowContains(repo, secretScanningKeywords)
	result.Passed = ciFound
	if ciFound {
		result.Evidence = "Secret scanning in CI: " + evidence
	} else {
		result.Evidence = "No secret scanning config or CI step found"
	}

	return result, nil
}
