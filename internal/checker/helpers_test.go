package checker_test

import (
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
)

func TestCIWorkflowContainsYamlWithKeyword(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{
			Data: []byte("name: CI\njobs:\n  lint:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: gitleaks\n"),
		},
	}
	found, msg := checker.CIWorkflowContains(repo, []string{"gitleaks"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if msg == "" {
		t.Errorf("expected non-empty message, got empty string")
	}
	if !contains(msg, "gitleaks") {
		t.Errorf("expected message to contain 'gitleaks', got: %s", msg)
	}
}

func TestCIWorkflowContainsYamlWithoutKeyword(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{
			Data: []byte("name: CI\njobs:\n  test:\n    runs-on: ubuntu-latest\n"),
		},
	}
	found, msg := checker.CIWorkflowContains(repo, []string{"gitleaks", "security"})
	if found {
		t.Errorf("expected found=false, got true; message: %s", msg)
	}
	if msg != "" {
		t.Errorf("expected empty message, got: %s", msg)
	}
}

func TestCIWorkflowContainsNoWorkflowsDir(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	found, msg := checker.CIWorkflowContains(repo, []string{"gitleaks"})
	if found {
		t.Errorf("expected found=false, got true; message: %s", msg)
	}
	if msg != "" {
		t.Errorf("expected empty message, got: %s", msg)
	}
}

func TestDepFileContainsGoModWithPackage(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{
			Data: []byte("module example.com/foo\n\nrequire (\n\tgo.uber.org/zap v1.24.0\n)"),
		},
	}
	found, msg := checker.DepFileContains(repo, checker.LanguageGo, []string{"go.uber.org/zap"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if msg == "" {
		t.Errorf("expected non-empty message, got empty string")
	}
	if !contains(msg, "go.uber.org/zap") {
		t.Errorf("expected message to contain 'go.uber.org/zap', got: %s", msg)
	}
}

func TestDepFileContainsPackageJsonWithPackage(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"name":"app","dependencies":{"react":"^18.0.0","lodash":"^4.17.0"}}`),
		},
	}
	found, msg := checker.DepFileContains(repo, checker.LanguageTypeScript, []string{"react"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if msg == "" {
		t.Errorf("expected non-empty message, got empty string")
	}
	if !contains(msg, "react") {
		t.Errorf("expected message to contain 'react', got: %s", msg)
	}
}

func TestDepFileContainsNeitherDepFileExists(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# Project")},
	}
	found, msg := checker.DepFileContains(repo, checker.LanguageGo, []string{"go.uber.org/zap"})
	if found {
		t.Errorf("expected found=false, got true; message: %s", msg)
	}
	if msg != "" {
		t.Errorf("expected empty message, got: %s", msg)
	}
}

func TestDepFileContainsUnknownLanguageFallback(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"dependencies":{"express":"^4.18.0"}}`),
		},
	}
	found, msg := checker.DepFileContains(repo, checker.LanguageUnknown, []string{"express"})
	if !found {
		t.Errorf("expected found=true with fallback, got false")
	}
	if msg == "" {
		t.Errorf("expected non-empty message, got empty string")
	}
	if !contains(msg, "express") {
		t.Errorf("expected message to contain 'express', got: %s", msg)
	}
}

func TestDepFileContainsJavaWithPomXml(t *testing.T) {
	repo := fstest.MapFS{
		"pom.xml": &fstest.MapFile{
			Data: []byte(`<project><dependencies><dependency><artifactId>junit</artifactId></dependency></dependencies></project>`),
		},
	}
	found, msg := checker.DepFileContains(repo, checker.LanguageJava, []string{"junit"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if msg == "" {
		t.Errorf("expected non-empty message, got empty string")
	}
	if !contains(msg, "junit") {
		t.Errorf("expected message to contain 'junit', got: %s", msg)
	}
}

func TestFileExistsAnyFirstCandidateExists(t *testing.T) {
	repo := fstest.MapFS{
		".goreleaser.yml": &fstest.MapFile{Data: []byte("version: 2")},
	}
	found, path := checker.FileExistsAny(repo, []string{".goreleaser.yml", ".goreleaser.yaml"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if path != ".goreleaser.yml" {
		t.Errorf("expected path='.goreleaser.yml', got: %s", path)
	}
}

func TestFileExistsAnySecondCandidateExists(t *testing.T) {
	repo := fstest.MapFS{
		".goreleaser.yaml": &fstest.MapFile{Data: []byte("version: 2")},
	}
	found, path := checker.FileExistsAny(repo, []string{".goreleaser.yml", ".goreleaser.yaml"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if path != ".goreleaser.yaml" {
		t.Errorf("expected path='.goreleaser.yaml', got: %s", path)
	}
}

func TestFileExistsAnyNoneExist(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	found, path := checker.FileExistsAny(repo, []string{".goreleaser.yml", ".goreleaser.yaml"})
	if found {
		t.Errorf("expected found=false, got true; path: %s", path)
	}
	if path != "" {
		t.Errorf("expected empty path, got: %s", path)
	}
}

func TestFileContentContainsKeywordFound(t *testing.T) {
	repo := fstest.MapFS{
		"Makefile": &fstest.MapFile{
			Data: []byte(".PHONY: release\nrelease:\n\tgoreleaser release\n"),
		},
	}
	found, keyword := checker.FileContentContains(repo, "Makefile", []string{"release", "build"})
	if !found {
		t.Errorf("expected found=true, got false")
	}
	if keyword != "release" {
		t.Errorf("expected keyword='release', got: %s", keyword)
	}
}

func TestFileContentContainsKeywordNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"Makefile": &fstest.MapFile{
			Data: []byte(".PHONY: test\ntest:\n\tgo test ./...\n"),
		},
	}
	found, keyword := checker.FileContentContains(repo, "Makefile", []string{"release", "deploy"})
	if found {
		t.Errorf("expected found=false, got true; keyword: %s", keyword)
	}
	if keyword != "" {
		t.Errorf("expected empty keyword, got: %s", keyword)
	}
}

func TestFileContentContainsFileMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	found, keyword := checker.FileContentContains(repo, "Makefile", []string{"release"})
	if found {
		t.Errorf("expected found=false, got true; keyword: %s", keyword)
	}
	if keyword != "" {
		t.Errorf("expected empty keyword, got: %s", keyword)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != ""
}
