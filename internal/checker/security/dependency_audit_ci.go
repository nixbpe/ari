package security

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DependencyAuditCIChecker struct{}

func (c *DependencyAuditCIChecker) ID() checker.CheckerID  { return "dependency_audit_ci" }
func (c *DependencyAuditCIChecker) Pillar() checker.Pillar { return checker.PillarSecurity }
func (c *DependencyAuditCIChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *DependencyAuditCIChecker) Name() string           { return "Dependency Audit in CI" }
func (c *DependencyAuditCIChecker) Description() string {
	return "Checks that dependency vulnerability auditing runs in CI (npm audit, govulncheck, trivy, snyk, etc.)"
}
func (c *DependencyAuditCIChecker) Suggestion() string {
	return "Add dependency auditing to CI: npm audit, govulncheck, cargo audit, trivy, snyk, grype, or osv-scanner"
}

var depAuditKeywords = []string{
	"npm audit",
	"govulncheck",
	"cargo audit",
	"trivy",
	"snyk",
	"grype",
	"osv-scanner",
}

func (c *DependencyAuditCIChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	found, evidence := checker.CIWorkflowContains(repo, depAuditKeywords)
	result.Passed = found
	if found {
		result.Evidence = "Dependency audit in CI: " + evidence
	} else {
		result.Evidence = "No dependency audit step found in CI workflows"
	}

	return result, nil
}
