package reporter_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scanner"
	"github.com/nixbpe/ari/internal/scorer"
)

func makeRepoInfo() *scanner.RepoInfo {
	return &scanner.RepoInfo{
		RootPath:   "/test/repo",
		Language:   checker.LanguageGo,
		IsGitRepo:  true,
		CommitHash: "abc123",
		Branch:     "main",
	}
}

func makeScore() *scorer.Score {
	return &scorer.Score{
		Level:    checker.LevelFunctional,
		PassRate: 0.75,
	}
}

func makeResults() []*checker.Result {
	return []*checker.Result{
		{
			ID:         "go-001",
			Name:       "Go modules present",
			Pillar:     checker.PillarBuildSystem,
			Level:      checker.LevelFunctional,
			Passed:     true,
			Evidence:   "go.mod found",
			Mode:       "rule-based",
			Suggestion: "",
		},
		{
			ID:         "go-002",
			Name:       "CI configuration present",
			Pillar:     checker.PillarBuildSystem,
			Level:      checker.LevelDocumented,
			Passed:     false,
			Evidence:   "no .github/workflows found",
			Mode:       "rule-based",
			Suggestion: "Add a GitHub Actions workflow",
		},
		{
			ID:         "go-003",
			Name:       "Linter configured",
			Pillar:     checker.PillarStyleValidation,
			Level:      checker.LevelStandardized,
			Passed:     false,
			Skipped:    true,
			SkipReason: "No Go toolchain detected",
			Mode:       "rule-based",
			Suggestion: "Add golangci-lint",
		},
	}
}

func TestBuildReport(t *testing.T) {
	repoInfo := makeRepoInfo()
	score := makeScore()
	results := makeResults()

	r := reporter.BuildReport(repoInfo, score, results)

	if r.RepoPath != repoInfo.RootPath {
		t.Errorf("RepoPath: got %q, want %q", r.RepoPath, repoInfo.RootPath)
	}
	if r.Language != checker.LanguageGo.String() {
		t.Errorf("Language: got %q, want %q", r.Language, checker.LanguageGo.String())
	}
	if !r.IsGitRepo {
		t.Error("IsGitRepo: expected true")
	}
	if r.CommitHash != repoInfo.CommitHash {
		t.Errorf("CommitHash: got %q, want %q", r.CommitHash, repoInfo.CommitHash)
	}
	if r.Branch != repoInfo.Branch {
		t.Errorf("Branch: got %q, want %q", r.Branch, repoInfo.Branch)
	}
	if r.Score != score {
		t.Error("Score: expected same pointer")
	}
	if r.Level != score.Level.String() {
		t.Errorf("Level: got %q, want %q", r.Level, score.Level.String())
	}
	if r.PassRate != score.PassRate {
		t.Errorf("PassRate: got %v, want %v", r.PassRate, score.PassRate)
	}
	if len(r.CriteriaResults) != len(results) {
		t.Errorf("CriteriaResults: got %d, want %d", len(r.CriteriaResults), len(results))
	}
	if r.AriVersion == "" {
		t.Error("AriVersion: expected non-empty")
	}
	if r.GeneratedAt.IsZero() {
		t.Error("GeneratedAt: expected non-zero time")
	}

	first := r.CriteriaResults[0]
	if first.ID != "go-001" {
		t.Errorf("CriteriaResults[0].ID: got %q, want %q", first.ID, "go-001")
	}
	if first.Pillar != checker.PillarBuildSystem.String() {
		t.Errorf("CriteriaResults[0].Pillar: got %q, want %q", first.Pillar, checker.PillarBuildSystem.String())
	}
	if first.Level != int(checker.LevelFunctional) {
		t.Errorf("CriteriaResults[0].Level: got %d, want %d", first.Level, int(checker.LevelFunctional))
	}
}

func TestJSONReporterOutput(t *testing.T) {
	r := reporter.BuildReport(makeRepoInfo(), makeScore(), makeResults())
	jr := &reporter.JSONReporter{}

	var buf bytes.Buffer
	err := jr.Report(context.Background(), r, &buf)
	if err != nil {
		t.Fatalf("Report() error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("Report() produced empty output")
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v\nOutput:\n%s", err, buf.String())
	}
}

func TestJSONRequiredFields(t *testing.T) {
	r := reporter.BuildReport(makeRepoInfo(), makeScore(), makeResults())
	jr := &reporter.JSONReporter{}

	var buf bytes.Buffer
	if err := jr.Report(context.Background(), r, &buf); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	required := []string{"level", "passRate", "criteria", "suggestions"}
	for _, field := range required {
		if _, ok := out[field]; !ok {
			t.Errorf("JSON missing required field: %q", field)
		}
	}
}

func TestJSONRoundTrip(t *testing.T) {
	score := &scorer.Score{
		Level:    checker.LevelDocumented,
		PassRate: 0.5,
	}
	original := reporter.BuildReport(makeRepoInfo(), score, makeResults())
	origTime := original.GeneratedAt

	jr := &reporter.JSONReporter{}
	var buf bytes.Buffer
	if err := jr.Report(context.Background(), original, &buf); err != nil {
		t.Fatalf("Report() error: %v", err)
	}

	var rt reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &rt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if rt.RepoPath != original.RepoPath {
		t.Errorf("RepoPath: got %q, want %q", rt.RepoPath, original.RepoPath)
	}
	if rt.Level != original.Level {
		t.Errorf("Level: got %q, want %q", rt.Level, original.Level)
	}
	if rt.PassRate != original.PassRate {
		t.Errorf("PassRate: got %v, want %v", rt.PassRate, original.PassRate)
	}
	if rt.Language != original.Language {
		t.Errorf("Language: got %q, want %q", rt.Language, original.Language)
	}
	if len(rt.CriteriaResults) != len(original.CriteriaResults) {
		t.Errorf("CriteriaResults count: got %d, want %d", len(rt.CriteriaResults), len(original.CriteriaResults))
	}
	if len(rt.Suggestions) != len(original.Suggestions) {
		t.Errorf("Suggestions count: got %d, want %d", len(rt.Suggestions), len(original.Suggestions))
	}
	if !rt.GeneratedAt.Equal(origTime) {
		t.Errorf("GeneratedAt mismatch: got %v, want %v", rt.GeneratedAt, origTime)
	}
}

func TestSuggestionsOnlyForFailing(t *testing.T) {
	results := []*checker.Result{
		{
			ID:         "pass-001",
			Name:       "Passing criterion",
			Pillar:     checker.PillarBuildSystem,
			Level:      checker.LevelFunctional,
			Passed:     true,
			Suggestion: "Should not appear",
		},
		{
			ID:         "skip-001",
			Name:       "Skipped criterion",
			Pillar:     checker.PillarBuildSystem,
			Level:      checker.LevelFunctional,
			Passed:     false,
			Skipped:    true,
			SkipReason: "Not applicable",
			Suggestion: "Should not appear either",
		},
		{
			ID:         "fail-001",
			Name:       "Failing criterion",
			Pillar:     checker.PillarBuildSystem,
			Level:      checker.LevelFunctional,
			Passed:     false,
			Skipped:    false,
			Suggestion: "Fix this",
		},
		{
			ID:         "fail-002",
			Name:       "Another failing criterion",
			Pillar:     checker.PillarTesting,
			Level:      checker.LevelDocumented,
			Passed:     false,
			Skipped:    false,
			Suggestion: "Fix that too",
		},
	}

	r := reporter.BuildReport(makeRepoInfo(), makeScore(), results)

	if len(r.Suggestions) != 2 {
		t.Errorf("Suggestions count: got %d, want 2", len(r.Suggestions))
	}

	ids := make(map[string]bool)
	for _, s := range r.Suggestions {
		ids[s.CriterionID] = true
	}
	if ids["pass-001"] {
		t.Error("Unexpected suggestion for passing criterion pass-001")
	}
	if ids["skip-001"] {
		t.Error("Unexpected suggestion for skipped criterion skip-001")
	}
	if !ids["fail-001"] {
		t.Error("Missing suggestion for failing criterion fail-001")
	}
	if !ids["fail-002"] {
		t.Error("Missing suggestion for failing criterion fail-002")
	}
}

func TestBuildReportEmptyResults(t *testing.T) {
	repoInfo := makeRepoInfo()
	score := &scorer.Score{}
	results := []*checker.Result{}

	r := reporter.BuildReport(repoInfo, score, results)

	if r == nil {
		t.Fatal("BuildReport returned nil")
	}
	if len(r.CriteriaResults) != 0 {
		t.Errorf("CriteriaResults: expected empty, got %d", len(r.CriteriaResults))
	}
	if len(r.Suggestions) != 0 {
		t.Errorf("Suggestions: expected empty, got %d", len(r.Suggestions))
	}
	if r.Score != score {
		t.Error("Score: expected same pointer")
	}
	if r.AriVersion == "" {
		t.Error("AriVersion: expected non-empty")
	}

	jr := &reporter.JSONReporter{}
	var buf bytes.Buffer
	if err := jr.Report(context.Background(), r, &buf); err != nil {
		t.Fatalf("Report() error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON with empty results: %v\nOutput: %s", err, buf.String())
	}
}
