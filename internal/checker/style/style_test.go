package style_test

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/checker/style"
)

func TestLintConfigGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod":        &fstest.MapFile{Data: []byte("module example.com/foo")},
		".golangci.yml": &fstest.MapFile{Data: []byte("run:\n  timeout: 5m\n")},
	}
	c := &style.LintConfigChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "golangci") {
		t.Errorf("evidence should contain 'golangci', got: %s", res.Evidence)
	}
}

func TestLintConfigGoMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &style.LintConfigChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestLintConfigTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json":   &fstest.MapFile{Data: []byte("{}")},
		".eslintrc.json": &fstest.MapFile{Data: []byte("{}")},
	}
	c := &style.LintConfigChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestLintConfigJava(t *testing.T) {
	repo := fstest.MapFS{
		"checkstyle.xml": &fstest.MapFile{Data: []byte("<checkstyle/>")},
	}
	c := &style.LintConfigChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageJava)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFormatterGoAlwaysPasses(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &style.FormatterChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFormatterTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte("{}")},
		".prettierrc":  &fstest.MapFile{Data: []byte("{}")},
	}
	c := &style.FormatterChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFormatterTSMissing(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte("{}")},
	}
	c := &style.FormatterChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestTypeCheckGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &style.TypeCheckChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestTypeCheckTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json":  &fstest.MapFile{Data: []byte("{}")},
		"tsconfig.json": &fstest.MapFile{Data: []byte("{}")},
	}
	c := &style.TypeCheckChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestTypeCheckJSOnly(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte("{}")},
	}
	c := &style.TypeCheckChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "JavaScript without TypeScript") {
		t.Errorf("evidence should contain 'JavaScript without TypeScript', got: %s", res.Evidence)
	}
}

func TestTypeCheckJava(t *testing.T) {
	repo := fstest.MapFS{
		"pom.xml": &fstest.MapFile{Data: []byte("<project/>")},
	}
	c := &style.TypeCheckChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageJava)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}
