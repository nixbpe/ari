package taskdiscovery

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
)

var ctx = context.Background()

// ── CONTRIBUTING GUIDE ────────────────────────────────────────────────────────

func TestContributingGuideFound(t *testing.T) {
	repo := fstest.MapFS{
		"CONTRIBUTING.md": &fstest.MapFile{Data: []byte("# Contributing\n\nHow to contribute...")},
	}
	c := &ContributingGuideChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestContributingGuideFoundInGitHub(t *testing.T) {
	repo := fstest.MapFS{
		".github/CONTRIBUTING.md": &fstest.MapFile{Data: []byte("# Contributing")},
	}
	c := &ContributingGuideChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for .github/CONTRIBUTING.md, got false; evidence: %s", r.Evidence)
	}
}

func TestContributingGuideMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project")},
	}
	c := &ContributingGuideChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no contributing guide), got true; evidence: %s", r.Evidence)
	}
}

func TestContributingGuidePillar(t *testing.T) {
	c := &ContributingGuideChecker{}
	if c.Pillar() != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", c.Pillar())
	}
	if c.Level() != checker.LevelFunctional {
		t.Errorf("expected LevelFunctional, got %v", c.Level())
	}
}

// ── ISSUE TEMPLATES ───────────────────────────────────────────────────────────

func TestIssueTemplatesFound(t *testing.T) {
	repo := fstest.MapFS{
		".github/ISSUE_TEMPLATE/bug.md": &fstest.MapFile{Data: []byte("## Bug Report")},
	}
	c := &IssueTemplatesChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestIssueTemplatesMultiple(t *testing.T) {
	repo := fstest.MapFS{
		".github/ISSUE_TEMPLATE/bug.md":     &fstest.MapFile{Data: []byte("## Bug Report")},
		".github/ISSUE_TEMPLATE/feature.md": &fstest.MapFile{Data: []byte("## Feature Request")},
	}
	c := &IssueTemplatesChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (2 templates), got false; evidence: %s", r.Evidence)
	}
}

func TestIssueTemplatesMissing(t *testing.T) {
	repo := fstest.MapFS{}
	c := &IssueTemplatesChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no dir), got true; evidence: %s", r.Evidence)
	}
}

func TestIssueTemplatesPillar(t *testing.T) {
	c := &IssueTemplatesChecker{}
	if c.Pillar() != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", c.Pillar())
	}
	if c.Level() != checker.LevelDocumented {
		t.Errorf("expected LevelDocumented, got %v", c.Level())
	}
}

// ── PR TEMPLATE ───────────────────────────────────────────────────────────────

func TestPRTemplateFound(t *testing.T) {
	repo := fstest.MapFS{
		".github/PULL_REQUEST_TEMPLATE.md": &fstest.MapFile{Data: []byte("## PR Description")},
	}
	c := &PRTemplateChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestPRTemplateFoundLowercase(t *testing.T) {
	repo := fstest.MapFS{
		".github/pull_request_template.md": &fstest.MapFile{Data: []byte("## PR Description")},
	}
	c := &PRTemplateChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for lowercase template, got false; evidence: %s", r.Evidence)
	}
}

func TestPRTemplateMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project")},
	}
	c := &PRTemplateChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no PR template), got true; evidence: %s", r.Evidence)
	}
}

func TestPRTemplatePillar(t *testing.T) {
	c := &PRTemplateChecker{}
	if c.Pillar() != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", c.Pillar())
	}
	if c.Level() != checker.LevelDocumented {
		t.Errorf("expected LevelDocumented, got %v", c.Level())
	}
}

// ── ISSUE LABELING SYSTEM ─────────────────────────────────────────────────────

func TestIssueLabelingSystemFound(t *testing.T) {
	repo := fstest.MapFS{
		".github/labels.yml": &fstest.MapFile{Data: []byte("- name: bug\n  color: d73a4a")},
	}
	c := &IssueLabelingSystemChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestIssueLabelingSystemFoundJson(t *testing.T) {
	repo := fstest.MapFS{
		".github/labels.json": &fstest.MapFile{Data: []byte(`[{"name":"bug","color":"d73a4a"}]`)},
	}
	c := &IssueLabelingSystemChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for labels.json, got false; evidence: %s", r.Evidence)
	}
}

func TestIssueLabelingSystemMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project")},
	}
	c := &IssueLabelingSystemChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no labels config), got true; evidence: %s", r.Evidence)
	}
}

func TestIssueLabelingSystemPillar(t *testing.T) {
	c := &IssueLabelingSystemChecker{}
	if c.Pillar() != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", c.Pillar())
	}
	if c.Level() != checker.LevelStandardized {
		t.Errorf("expected LevelStandardized, got %v", c.Level())
	}
}

// ── BACKLOG STRUCTURE DOCS ────────────────────────────────────────────────────

func TestBacklogStructureDocsFoundInAgentsMd(t *testing.T) {
	repo := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# AGENTS\n\n## Backlog\nWe use priority labels to triage issues.")},
	}
	c := &BacklogStructureDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (AGENTS.md has 'priority'), got false; evidence: %s", r.Evidence)
	}
}

func TestBacklogStructureDocsFoundInProcess(t *testing.T) {
	repo := fstest.MapFS{
		"docs/process.md": &fstest.MapFile{Data: []byte("# Process\n\nWe use milestones and sprint planning to manage our backlog.")},
	}
	c := &BacklogStructureDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (docs/process.md has 'backlog'), got false; evidence: %s", r.Evidence)
	}
}

func TestBacklogStructureDocsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project\n\nJust a readme.")},
	}
	c := &BacklogStructureDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no process docs), got true; evidence: %s", r.Evidence)
	}
}

func TestBacklogStructureDocsNoKeywords(t *testing.T) {
	repo := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# AGENTS\n\nBuild commands and architecture overview only.")},
	}
	c := &BacklogStructureDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (AGENTS.md has no backlog keywords), got true; evidence: %s", r.Evidence)
	}
}

func TestBacklogStructureDocsPillar(t *testing.T) {
	c := &BacklogStructureDocsChecker{}
	if c.Pillar() != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", c.Pillar())
	}
	if c.Level() != checker.LevelOptimized {
		t.Errorf("expected LevelOptimized, got %v", c.Level())
	}
}
