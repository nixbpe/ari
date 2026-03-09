package build_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/checker/build"
)

func TestFeatureFlagGo(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/myapp\n\nrequire github.com/open-feature/go-sdk-contrib/providers/openfeature v0.1.0\n")},
	}
	c := &build.FeatureFlagInfrastructureChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFeatureFlagTS(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"@openfeature/js-sdk":"^1.0.0"}}`)},
	}
	c := &build.FeatureFlagInfrastructureChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFeatureFlagMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/myapp\n\nrequire github.com/some/other v0.1.0\n")},
	}
	c := &build.FeatureFlagInfrastructureChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestReleaseNotesChangelog(t *testing.T) {
	repo := fstest.MapFS{
		"CHANGELOG.md": &fstest.MapFile{Data: []byte("# Changelog\n\n## v1.0.0\n- Initial release\n")},
	}
	c := &build.ReleaseNotesAutomationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestReleaseNotesMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/myapp\n")},
	}
	c := &build.ReleaseNotesAutomationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestBuildPerfTurbo(t *testing.T) {
	repo := fstest.MapFS{
		"turbo.json": &fstest.MapFile{Data: []byte(`{"pipeline":{"build":{}}}`)},
	}
	c := &build.BuildPerformanceTrackingChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestBuildPerfMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/myapp\n")},
	}
	c := &build.BuildPerformanceTrackingChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestUnusedDepsGo(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{Data: []byte("name: CI\njobs:\n  tidy:\n    steps:\n      - run: go mod tidy\n")},
	}
	c := &build.UnusedDependenciesDetectionChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestBuildAdvancedMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod":  &fstest.MapFile{Data: []byte("module example.com/myapp\n")},
		"main.go": &fstest.MapFile{Data: []byte("package main\nfunc main() {}\n")},
	}
	ctx := context.Background()

	buildPerf := &build.BuildPerformanceTrackingChecker{}
	res, err := buildPerf.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("build_performance_tracking error: %v", err)
	}
	if res.Passed {
		t.Errorf("build_performance_tracking: expected Passed=false, got true")
	}

	featureFlag := &build.FeatureFlagInfrastructureChecker{}
	res, err = featureFlag.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("feature_flag_infrastructure error: %v", err)
	}
	if res.Passed {
		t.Errorf("feature_flag_infrastructure: expected Passed=false, got true")
	}

	releaseNotes := &build.ReleaseNotesAutomationChecker{}
	res, err = releaseNotes.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("release_notes_automation error: %v", err)
	}
	if res.Passed {
		t.Errorf("release_notes_automation: expected Passed=false, got true")
	}

	unusedDeps := &build.UnusedDependenciesDetectionChecker{}
	res, err = unusedDeps.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unused_dependencies_detection error: %v", err)
	}
	if res.Passed {
		t.Errorf("unused_dependencies_detection: expected Passed=false, got true")
	}
}
