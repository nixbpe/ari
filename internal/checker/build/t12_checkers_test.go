package build_test

import (
	"context"
	"fmt"
	"testing"
	"testing/fstest"
	"time"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/checker/build"
)

func TestFastCIFeedback(t *testing.T) {
	repo := fstest.MapFS{
		".github/workflows/ci.yml": &fstest.MapFile{Data: []byte("on: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n")},
	}
	c := &build.FastCIFeedbackChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFastCIFeedbackGitLab(t *testing.T) {
	repo := fstest.MapFS{
		".gitlab-ci.yml": &fstest.MapFile{Data: []byte("stages:\n  - test\n")},
	}
	c := &build.FastCIFeedbackChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestFastCIFeedbackMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &build.FastCIFeedbackChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestReleaseAutomation(t *testing.T) {
	repo := fstest.MapFS{
		".goreleaser.yml": &fstest.MapFile{Data: []byte("project_name: myapp\n")},
	}
	c := &build.ReleaseAutomationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestReleaseAutomationMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo")},
	}
	c := &build.ReleaseAutomationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestDeployFreqNoGit(t *testing.T) {
	c := &build.DeploymentFrequencyChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return nil, fmt.Errorf("not a git repository")
		},
	}
	res, err := c.Check(context.Background(), fstest.MapFS{}, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Errorf("expected Skipped=true, got false; evidence: %s", res.Evidence)
	}
}

func TestDeployFreqRecentTag(t *testing.T) {
	recent := time.Now().AddDate(0, 0, -10).Format("2006-01-02 15:04:05 -0700")
	c := &build.DeploymentFrequencyChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return []byte("v1.2.3|" + recent + "\n"), nil
		},
	}
	res, err := c.Check(context.Background(), fstest.MapFS{}, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestDeployFreqOldTag(t *testing.T) {
	old := time.Now().AddDate(0, 0, -120).Format("2006-01-02 15:04:05 -0700")
	c := &build.DeploymentFrequencyChecker{
		GitRunner: func(args ...string) ([]byte, error) {
			return []byte("v0.1.0|" + old + "\n"), nil
		},
	}
	res, err := c.Check(context.Background(), fstest.MapFS{}, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}
