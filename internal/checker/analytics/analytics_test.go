package analytics

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
)

var ctx = context.Background()

// ── ANALYTICS SDK ─────────────────────────────────────────────────────────────

func TestAnalyticsSdkFound(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"posthog-js":"^1.0.0"}}`)},
	}
	c := &AnalyticsSdkChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
	if r.Pillar != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", r.Pillar)
	}
}

func TestAnalyticsSdkSegment(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"@segment/analytics-next":"^1.0.0"}}`)},
	}
	c := &AnalyticsSdkChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for @segment/analytics-next; evidence: %s", r.Evidence)
	}
}

func TestAnalyticsSdkMissing(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"lodash":"^4.0.0"}}`)},
	}
	c := &AnalyticsSdkChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false; evidence: %s", r.Evidence)
	}
}

// ── TRACKING PLAN DOCS ────────────────────────────────────────────────────────

func TestTrackingPlanDocsMd(t *testing.T) {
	repo := fstest.MapFS{
		"docs/tracking-plan.md": &fstest.MapFile{Data: []byte("# Tracking Plan\n...")},
	}
	c := &TrackingPlanDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true; evidence: %s", r.Evidence)
	}
	if r.Pillar != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", r.Pillar)
	}
}

func TestTrackingPlanDocsAvoJson(t *testing.T) {
	repo := fstest.MapFS{
		"avo.json": &fstest.MapFile{Data: []byte(`{"events":[]}`)},
	}
	c := &TrackingPlanDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for avo.json; evidence: %s", r.Evidence)
	}
}

func TestTrackingPlanDocsAvoDir(t *testing.T) {
	repo := fstest.MapFS{
		".avo/tracking-plan.json": &fstest.MapFile{Data: []byte(`{}`)},
	}
	c := &TrackingPlanDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for .avo/ dir; evidence: %s", r.Evidence)
	}
}

func TestTrackingPlanDocsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# App")},
	}
	c := &TrackingPlanDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false; evidence: %s", r.Evidence)
	}
}

// ── EXPERIMENT INFRASTRUCTURE ─────────────────────────────────────────────────

func TestExperimentInfrastructureFound(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"@growthbook/growthbook":"^0.20.0"}}`)},
	}
	c := &ExperimentInfrastructureChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true; evidence: %s", r.Evidence)
	}
	if r.Pillar != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", r.Pillar)
	}
}

func TestExperimentInfrastructureMissing(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"react":"^18.0.0"}}`)},
	}
	c := &ExperimentInfrastructureChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false; evidence: %s", r.Evidence)
	}
}

// ── PRODUCT METRICS DOCS ──────────────────────────────────────────────────────

func TestProductMetricsDocsFound(t *testing.T) {
	repo := fstest.MapFS{
		"docs/metrics.md": &fstest.MapFile{Data: []byte("# Metrics\n\n## North Star\nour north star metric is DAU/MAU ratio\n\n## KPIs\n- conversion rate\n- retention")},
	}
	c := &ProductMetricsDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true; evidence: %s", r.Evidence)
	}
	if r.Pillar != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", r.Pillar)
	}
}

func TestProductMetricsDocsMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# App")},
	}
	c := &ProductMetricsDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false; evidence: %s", r.Evidence)
	}
}

func TestProductMetricsDocsFileExistsButNoKeywords(t *testing.T) {
	repo := fstest.MapFS{
		"docs/kpis.md": &fstest.MapFile{Data: []byte("This document is a placeholder.")},
	}
	c := &ProductMetricsDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false (no keywords); evidence: %s", r.Evidence)
	}
}

// ── ERROR-TO-INSIGHT PIPELINE ─────────────────────────────────────────────────

func TestErrorToInsightSentryConfig(t *testing.T) {
	repo := fstest.MapFS{
		".sentry.properties": &fstest.MapFile{Data: []byte("defaults.project=my-project")},
	}
	c := &ErrorToInsightPipelineChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true; evidence: %s", r.Evidence)
	}
	if r.Pillar != checker.PillarContextIntent {
		t.Errorf("expected PillarContextIntent, got %v", r.Pillar)
	}
}

func TestErrorToInsightSentryCIWorkflow(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/release.yml": &fstest.MapFile{Data: []byte("steps:\n  - run: sentry-cli releases finalize $VERSION\n")},
	}
	c := &ErrorToInsightPipelineChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true; evidence: %s", r.Evidence)
	}
}

func TestErrorToInsightGHIssueCIWorkflow(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/monitor.yml": &fstest.MapFile{Data: []byte("steps:\n  - run: gh issue create --title \"Alert\" --body \"$MSG\"\n")},
	}
	c := &ErrorToInsightPipelineChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true for gh issue create; evidence: %s", r.Evidence)
	}
}

func TestErrorToInsightMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# App")},
	}
	c := &ErrorToInsightPipelineChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false; evidence: %s", r.Evidence)
	}
}
