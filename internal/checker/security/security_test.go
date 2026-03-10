package security

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
)

var ctx = context.Background()

func TestSecurityPolicyFound(t *testing.T) {
	repo := fstest.MapFS{
		"SECURITY.md": &fstest.MapFile{Data: []byte("# Security Policy\n## Reporting Vulnerabilities")},
	}
	c := &SecurityPolicyChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestSecurityPolicyGithubDir(t *testing.T) {
	repo := fstest.MapFS{
		".github/SECURITY.md": &fstest.MapFile{Data: []byte("# Security")},
	}
	c := &SecurityPolicyChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for .github/SECURITY.md, got false; evidence: %s", r.Evidence)
	}
}

func TestSecurityPolicyMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &SecurityPolicyChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestGitignoreComprehensivePass(t *testing.T) {
	repo := fstest.MapFS{
		".gitignore": &fstest.MapFile{Data: []byte(".env\n*.pem\n*.key\nnode_modules\n*.log\n")},
	}
	c := &GitignoreComprehensiveChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (5 patterns), got false; evidence: %s", r.Evidence)
	}
}

func TestGitignoreComprehensiveTooFewPatterns(t *testing.T) {
	repo := fstest.MapFS{
		".gitignore": &fstest.MapFile{Data: []byte(".env\n*.pem\n")},
	}
	c := &GitignoreComprehensiveChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (only 2 patterns), got true; evidence: %s", r.Evidence)
	}
}

func TestGitignoreMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &GitignoreComprehensiveChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no .gitignore), got true; evidence: %s", r.Evidence)
	}
}

func TestCodeownersFound(t *testing.T) {
	repo := fstest.MapFS{
		".github/CODEOWNERS": &fstest.MapFile{Data: []byte("* @owner")},
	}
	c := &CodeownersChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestCodeownersMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &CodeownersChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestDepUpdateDependabot(t *testing.T) {
	repo := fstest.MapFS{
		".github/dependabot.yml": &fstest.MapFile{Data: []byte("version: 2\nupdates:\n  - package-ecosystem: gomod")},
	}
	c := &DepUpdateAutomationChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (dependabot.yml), got false; evidence: %s", r.Evidence)
	}
}

func TestDepUpdateRenovate(t *testing.T) {
	repo := fstest.MapFS{
		"renovate.json": &fstest.MapFile{Data: []byte(`{"$schema": "https://docs.renovatebot.com/renovate-schema.json"}`)},
	}
	c := &DepUpdateAutomationChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (renovate.json), got false; evidence: %s", r.Evidence)
	}
}

func TestDepUpdateMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &DepUpdateAutomationChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestSecretScanningWithFile(t *testing.T) {
	repo := fstest.MapFS{
		".gitleaks.toml": &fstest.MapFile{Data: []byte("title = \"gitleaks\"")},
	}
	c := &SecretScanningConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (.gitleaks.toml), got false; evidence: %s", r.Evidence)
	}
}

func TestSecretScanningWithCI(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{Data: []byte("- uses: gitleaks/gitleaks-action")},
	}
	c := &SecretScanningConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (CI gitleaks), got false; evidence: %s", r.Evidence)
	}
}

func TestSecretScanningNone(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &SecretScanningConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestSASTWithFile(t *testing.T) {
	repo := fstest.MapFS{
		".semgrep.yml": &fstest.MapFile{Data: []byte("rules: []")},
	}
	c := &SASTConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (.semgrep.yml), got false; evidence: %s", r.Evidence)
	}
}

func TestSASTWithCI(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/security.yml": &fstest.MapFile{Data: []byte("- uses: github/codeql-action/analyze@v2")},
	}
	c := &SASTConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (CI codeql), got false; evidence: %s", r.Evidence)
	}
}

func TestSASTNone(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &SASTConfigChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestDepAuditCIPass(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{Data: []byte("run: govulncheck ./...")},
	}
	c := &DependencyAuditCIChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (govulncheck), got false; evidence: %s", r.Evidence)
	}
}

func TestDepAuditCITrivy(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/scan.yml": &fstest.MapFile{Data: []byte("- uses: aquasecurity/trivy-action@master")},
	}
	c := &DependencyAuditCIChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (trivy), got false; evidence: %s", r.Evidence)
	}
}

func TestDepAuditCIMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &DependencyAuditCIChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestAllCheckersReturnSecurityPillar(t *testing.T) {
	checkers := []checker.Checker{
		&SecurityPolicyChecker{},
		&GitignoreComprehensiveChecker{},
		&CodeownersChecker{},
		&DepUpdateAutomationChecker{},
		&SecretScanningConfigChecker{},
		&SASTConfigChecker{},
		&DependencyAuditCIChecker{},
	}
	for _, c := range checkers {
		if c.Pillar() != checker.PillarConstraints {
			t.Errorf("%s: expected PillarConstraints, got %v", c.ID(), c.Pillar())
		}
	}
}
