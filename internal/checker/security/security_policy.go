package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type SecurityPolicyChecker struct{}

func (c *SecurityPolicyChecker) ID() checker.CheckerID  { return "security_policy" }
func (c *SecurityPolicyChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *SecurityPolicyChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *SecurityPolicyChecker) Name() string           { return "Security Policy" }
func (c *SecurityPolicyChecker) Description() string {
	return "Checks that a SECURITY.md file exists describing how to report vulnerabilities"
}
func (c *SecurityPolicyChecker) Suggestion() string {
	return "Create SECURITY.md (or .github/SECURITY.md) with: vulnerability reporting process, response timeline, supported versions"
}

var securityPolicyCandidates = []string{
	"SECURITY.md",
	".github/SECURITY.md",
	"docs/SECURITY.md",
}

func (c *SecurityPolicyChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	found, path := checker.FileExistsAny(repo, securityPolicyCandidates)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No SECURITY.md found (checked SECURITY.md, .github/SECURITY.md, docs/SECURITY.md)"
	}

	return result, nil
}
