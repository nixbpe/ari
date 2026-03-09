package style_test

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/checker/style"
)

func TestCyclomaticComplexityGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - govet
    - gocyclo
    - errcheck
`)},
	}

	c := &style.CyclomaticComplexityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", result.Evidence)
	}
	if !strings.Contains(result.Evidence, "gocyclo") {
		t.Errorf("expected evidence to mention gocyclo, got: %s", result.Evidence)
	}
}

func TestCyclomaticComplexityGoMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - govet
    - errcheck
`)},
	}

	c := &style.CyclomaticComplexityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Passed {
		t.Errorf("expected Passed=false when no complexity tool in config")
	}
}

func TestCyclomaticComplexityNoConfig(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
	}

	c := &style.CyclomaticComplexityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Passed {
		t.Errorf("expected Passed=false when no golangci config exists")
	}
	if result.Suggestion == "" {
		t.Errorf("expected non-empty suggestion")
	}
}

func TestCyclomaticComplexityGoGocognit(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - gocognit
`)},
	}

	c := &style.CyclomaticComplexityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for gocognit, got false")
	}
}

func TestCyclomaticComplexityTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"name":"app"}`)},
		".eslintrc.json": &fstest.MapFile{Data: []byte(`{
  "rules": {
    "complexity": ["error", 10]
  }
}`)},
	}

	c := &style.CyclomaticComplexityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for TS with complexity rule, got false; evidence: %s", result.Evidence)
	}
}

func TestDeadCodeGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - govet
    - deadcode
    - errcheck
`)},
	}

	c := &style.DeadCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", result.Evidence)
	}
	if !strings.Contains(result.Evidence, "deadcode") {
		t.Errorf("expected evidence to mention deadcode, got: %s", result.Evidence)
	}
}

func TestDeadCodeGoUnparam(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - unparam
`)},
	}

	c := &style.DeadCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for unparam, got false")
	}
}

func TestDeadCodeTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{
  "name": "my-app",
  "devDependencies": {
    "typescript": "^5.0.0",
    "knip": "^3.0.0"
  }
}`)},
	}

	c := &style.DeadCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for TS with knip, got false; evidence: %s", result.Evidence)
	}
	if !strings.Contains(result.Evidence, "knip") {
		t.Errorf("expected evidence to mention knip, got: %s", result.Evidence)
	}
}

func TestDeadCodeTSMissing(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{
  "name": "my-app",
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}`)},
	}

	c := &style.DeadCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Passed {
		t.Errorf("expected Passed=false when no dead code tool present")
	}
}

func TestDuplicateCodeGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".golangci.yml": &fstest.MapFile{Data: []byte(`
linters:
  enable:
    - govet
    - dupl
    - errcheck
`)},
	}

	c := &style.DuplicateCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", result.Evidence)
	}
	if !strings.Contains(result.Evidence, "dupl") {
		t.Errorf("expected evidence to mention dupl, got: %s", result.Evidence)
	}
}

func TestDuplicateCodeJscpdJson(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod":      &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
		".jscpd.json": &fstest.MapFile{Data: []byte(`{"threshold": 5}`)},
	}

	c := &style.DuplicateCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for .jscpd.json, got false; evidence: %s", result.Evidence)
	}
}

func TestDuplicateCodePackageJsonScript(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{
  "name": "my-app",
  "scripts": {
    "check-duplication": "jscpd src/"
  }
}`)},
	}

	c := &style.DuplicateCodeDetectionChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Passed {
		t.Errorf("expected Passed=true for jscpd in scripts, got false; evidence: %s", result.Evidence)
	}
}

func TestAnalysisToolsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/app\n\ngo 1.21\n")},
	}

	cyclomatic := &style.CyclomaticComplexityChecker{}
	deadCode := &style.DeadCodeDetectionChecker{}
	duplicate := &style.DuplicateCodeDetectionChecker{}

	cr, err := cyclomatic.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("cyclomatic.Check error: %v", err)
	}
	if cr.Passed {
		t.Errorf("cyclomatic: expected Passed=false with no config")
	}
	if cr.Suggestion == "" {
		t.Errorf("cyclomatic: expected non-empty suggestion")
	}

	dr, err := deadCode.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("deadCode.Check error: %v", err)
	}
	if dr.Passed {
		t.Errorf("deadCode: expected Passed=false with no config")
	}
	if dr.Suggestion == "" {
		t.Errorf("deadCode: expected non-empty suggestion")
	}

	dup, err := duplicate.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("duplicate.Check error: %v", err)
	}
	if dup.Passed {
		t.Errorf("duplicate: expected Passed=false with no config")
	}
	if dup.Suggestion == "" {
		t.Errorf("duplicate: expected non-empty suggestion")
	}
}
