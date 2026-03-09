package docs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/nixbpe/ari/internal/checker"
)

var ctx = context.Background()

// ── README ────────────────────────────────────────────────────────────────────

func TestReadmeGood(t *testing.T) {
	content := "# Title\n" + strings.Repeat("x", 100)
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte(content)},
	}
	c := &ReadmeChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestReadmeMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &ReadmeChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

func TestReadmeTooShort(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("Hello")},
	}
	c := &ReadmeChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (too short), got true; evidence: %s", r.Evidence)
	}
}

// ── AGENTS.MD ─────────────────────────────────────────────────────────────────

func TestAgentsMd(t *testing.T) {
	repo := fstest.MapFS{
		"CLAUDE.md": &fstest.MapFile{Data: []byte(strings.Repeat("a", 150))},
	}
	c := &AgentsMdChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestAgentsMdMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &AgentsMdChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

// ── DOCUMENTATION FRESHNESS ───────────────────────────────────────────────────

func TestDocFreshnessSkipNoGit(t *testing.T) {
	repo := fstest.MapFS{}
	c := &DocumentationFreshnessChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return nil, errors.New("not a git repository")
		},
	}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Skipped {
		t.Errorf("expected Skipped=true, got false; evidence: %s", r.Evidence)
	}
}

func TestDocFreshnessRecent(t *testing.T) {
	recentDate := time.Now().AddDate(0, 0, -10).Format(gitDateFormat)
	repo := fstest.MapFS{}
	c := &DocumentationFreshnessChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return []byte(recentDate + "\n"), nil
		},
	}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (recent), got false; evidence: %s", r.Evidence)
	}
}

func TestDocFreshnessStale(t *testing.T) {
	staleDate := time.Now().AddDate(-1, 0, 0).Format(gitDateFormat)
	repo := fstest.MapFS{}
	c := &DocumentationFreshnessChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return []byte(staleDate + "\n"), nil
		},
	}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (stale), got true; evidence: %s", r.Evidence)
	}
}

// ── SKILLS ────────────────────────────────────────────────────────────────────

func TestSkillsFound(t *testing.T) {
	repo := fstest.MapFS{
		".claude/skills/git-master.md": &fstest.MapFile{Data: []byte("# Git Master skill")},
	}
	c := &SkillsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestSkillsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &SkillsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

// ── AUTOMATED DOC GENERATION ──────────────────────────────────────────────────

func TestAutomatedDocGenTS(t *testing.T) {
	pkgJSON := fmt.Sprintf(`{"name":"app","devDependencies":{%q:"^0.23.0"}}`, "typedoc")
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(pkgJSON)},
	}
	c := &AutomatedDocGenerationChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (typedoc found), got false; evidence: %s", r.Evidence)
	}
}
