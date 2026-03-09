package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/all"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scanner"
	"github.com/nixbpe/ari/internal/scorer"
)

func runPipeline(t *testing.T, repoPath string) (*reporter.Report, *scorer.Score) {
	t.Helper()

	ctx := context.Background()
	repoFS := os.DirFS(repoPath)

	sc := scanner.NewScanner()
	repoInfo, err := sc.Scan(ctx, repoFS)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	repoInfo.RootPath = repoPath

	registry := checker.NewDefaultRegistry()
	all.RegisterAll(registry, nil)

	rnr := &checker.Runner{Registry: registry}
	results, err := rnr.Run(ctx, repoFS, repoInfo)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	score := scorer.New().Calculate(results)
	report := reporter.BuildReport(repoInfo, score, results)
	return report, score
}

func testdataPath(name string) string {
	return filepath.Join("..", "testdata", name)
}

func TestFullPipelineGoRepo(t *testing.T) {
	report, score := runPipeline(t, testdataPath("sample-go-repo"))

	if score == nil {
		t.Fatal("score is nil")
	}
	if len(report.CriteriaResults) == 0 {
		t.Fatal("no criteria results")
	}
	if len(report.CriteriaResults) != 40 {
		t.Errorf("expected 40 criteria, got %d", len(report.CriteriaResults))
	}
	if int(score.Level) < 1 {
		t.Errorf("expected level >= 1 for sample-go-repo, got %d", int(score.Level))
	}
}

func TestFullPipelineEmptyRepo(t *testing.T) {
	_, score := runPipeline(t, testdataPath("empty-repo"))

	if score == nil {
		t.Fatal("score is nil")
	}
	if int(score.Level) > 1 {
		t.Errorf("expected level 0 or 1 for empty-repo, got %d", int(score.Level))
	}
}

func TestFullPipelineNoGitRepo(t *testing.T) {
	report, _ := runPipeline(t, testdataPath("no-git-repo"))
	if report.IsGitRepo {
		t.Error("expected IsGitRepo=false for no-git-repo")
	}
}

func TestFullPipelineWellConfiguredRepo(t *testing.T) {
	_, score := runPipeline(t, testdataPath("well-configured-repo"))
	if score == nil {
		t.Fatal("score is nil")
	}
	if int(score.Level) < 1 {
		t.Errorf("expected level >= 1 for well-configured-repo, got %d", int(score.Level))
	}
}

func TestJSONOutputValid(t *testing.T) {
	report, _ := runPipeline(t, testdataPath("sample-go-repo"))

	if report.Language == "" {
		t.Error("report missing language")
	}
	if report.AriVersion == "" {
		t.Error("report missing ari version")
	}
	if report.GeneratedAt.IsZero() {
		t.Error("report missing generated_at")
	}
}

func TestAllCriteriaEvaluated(t *testing.T) {
	report, _ := runPipeline(t, testdataPath("sample-go-repo"))

	seen := make(map[string]bool)
	for _, cr := range report.CriteriaResults {
		seen[cr.ID] = true
	}

	required := []string{"lint_config", "unit_tests_exist", "readme", "agents_md"}
	for _, id := range required {
		if !seen[id] {
			t.Errorf("criterion %q not found in results", id)
		}
	}
}
