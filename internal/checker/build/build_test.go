package build_test

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/checker/build"
)

func TestDepsPinnedGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo\n\ngo 1.21\n")},
		"go.sum": &fstest.MapFile{Data: []byte("example.com/dep v1.0.0 h1:abc123\n")},
	}
	c := &build.DepsPinnedChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "go.sum") {
		t.Errorf("evidence should mention go.sum, got: %s", res.Evidence)
	}
}

func TestDepsPinnedGoMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo\n\ngo 1.21\n")},
	}
	c := &build.DepsPinnedChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestDepsPinnedTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json":      &fstest.MapFile{Data: []byte("{}")},
		"package-lock.json": &fstest.MapFile{Data: []byte("{\"lockfileVersion\": 2}")},
	}
	c := &build.DepsPinnedChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "package-lock.json") {
		t.Errorf("evidence should mention package-lock.json, got: %s", res.Evidence)
	}
}

func TestBuildCmdDocGo(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project\n\n## Build\n\nRun `go build ./...` to compile.\n")},
	}
	c := &build.BuildCmdDocChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "go build") {
		t.Errorf("evidence should mention 'go build', got: %s", res.Evidence)
	}
}

func TestBuildCmdDocMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project\n\nThis project does amazing things.\n")},
	}
	c := &build.BuildCmdDocChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Suggestion, "Document build commands") {
		t.Errorf("suggestion should mention documenting build commands, got: %s", res.Suggestion)
	}
}

func TestSingleCommandSetup(t *testing.T) {
	repo := fstest.MapFS{
		"Makefile": &fstest.MapFile{Data: []byte(".PHONY: setup\nsetup:\n\tgo mod download\n\tnpm install\n")},
	}
	c := &build.SingleCommandSetupChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "Makefile") {
		t.Errorf("evidence should mention Makefile, got: %s", res.Evidence)
	}
}

func TestSingleCommandSetupMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo\n\ngo 1.21\n")},
	}
	c := &build.SingleCommandSetupChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestAgenticDevelopment(t *testing.T) {
	repo := fstest.MapFS{
		"CLAUDE.md": &fstest.MapFile{Data: []byte("# Claude instructions")},
	}
	c := &build.AgenticDevelopmentChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "CLAUDE.md") {
		t.Errorf("evidence should mention CLAUDE.md, got: %s", res.Evidence)
	}
}

func TestAgenticDevelopmentAgentsMd(t *testing.T) {
	repo := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Agents instructions")},
	}
	c := &build.AgenticDevelopmentChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "AGENTS.md") {
		t.Errorf("evidence should mention AGENTS.md, got: %s", res.Evidence)
	}
}

func TestAgenticDevMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &build.AgenticDevelopmentChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Suggestion, "AGENTS.md") {
		t.Errorf("suggestion should mention AGENTS.md, got: %s", res.Suggestion)
	}
}

func TestVCSCliTools(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{Data: []byte("on: push")},
	}
	c := &build.VCSCliToolsChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, ".github/") {
		t.Errorf("evidence should mention .github/, got: %s", res.Evidence)
	}
}

func TestVCSCliToolsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &build.VCSCliToolsChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestAutomatedPRReview(t *testing.T) {
	repo := fstest.MapFS{
		"CODEOWNERS": &fstest.MapFile{Data: []byte("* @team")},
	}
	c := &build.AutomatedPRReviewChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "CODEOWNERS") {
		t.Errorf("evidence should mention CODEOWNERS, got: %s", res.Evidence)
	}
}

func TestAutomatedPRReviewMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &build.AutomatedPRReviewChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}
