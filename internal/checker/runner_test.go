package checker_test

import (
	"context"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
)

type runMockChecker struct {
	id          checker.CheckerID
	pillar      checker.Pillar
	level       checker.Level
	name        string
	description string
	checkFn     func(context.Context, fs.FS, checker.Language) (*checker.Result, error)
	called      int
}

func (m *runMockChecker) ID() checker.CheckerID { return m.id }

func (m *runMockChecker) Pillar() checker.Pillar { return m.pillar }

func (m *runMockChecker) Level() checker.Level { return m.level }

func (m *runMockChecker) Name() string { return m.name }

func (m *runMockChecker) Description() string { return m.description }

func (m *runMockChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	m.called++
	if m.checkFn != nil {
		return m.checkFn(ctx, repo, lang)
	}
	return &checker.Result{ID: m.id, Passed: true}, nil
}

type langScopedMockChecker struct {
	*runMockChecker
	supported map[checker.Language]bool
}

func (m *langScopedMockChecker) SupportsLanguage(lang checker.Language) bool {
	return m.supported[lang]
}

func TestRunnerExecutesAllCheckers(t *testing.T) {
	reg := checker.NewDefaultRegistry()
	first := &runMockChecker{id: "a", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "A"}
	second := &runMockChecker{id: "b", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "B"}
	third := &runMockChecker{id: "c", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "C"}

	if err := reg.Register(first); err != nil {
		t.Fatalf("register first: %v", err)
	}
	if err := reg.Register(second); err != nil {
		t.Fatalf("register second: %v", err)
	}
	if err := reg.Register(third); err != nil {
		t.Fatalf("register third: %v", err)
	}

	progressCalls := 0
	runner := &checker.Runner{
		Registry: reg,
		ProgressFunc: func(done, total int, id checker.CheckerID) {
			progressCalls++
			if total != 3 {
				t.Fatalf("Progress total = %d, want 3", total)
			}
			if done < 1 || done > 3 {
				t.Fatalf("Progress done = %d, want 1..3", done)
			}
			if id == "" {
				t.Fatal("Progress id should not be empty")
			}
		},
	}

	results, err := runner.Run(context.Background(), fstest.MapFS{}, struct{ Language checker.Language }{Language: checker.LanguageGo})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("Run() len(results) = %d, want 3", len(results))
	}
	if progressCalls != 3 {
		t.Fatalf("Progress calls = %d, want 3", progressCalls)
	}
	if first.called != 1 || second.called != 1 || third.called != 1 {
		t.Fatalf("checker called counts = [%d %d %d], want [1 1 1]", first.called, second.called, third.called)
	}
}

func TestRunnerSkipsInapplicableCheckers(t *testing.T) {
	reg := checker.NewDefaultRegistry()
	ch := &langScopedMockChecker{
		runMockChecker: &runMockChecker{
			id:     "go-only",
			pillar: checker.PillarBuildSystem,
			level:  checker.LevelFunctional,
			name:   "Go Only",
		},
		supported: map[checker.Language]bool{checker.LanguageGo: true},
	}

	if err := reg.Register(ch); err != nil {
		t.Fatalf("register checker: %v", err)
	}

	runner := &checker.Runner{Registry: reg}
	results, err := runner.Run(context.Background(), fstest.MapFS{}, struct{ Language checker.Language }{Language: checker.LanguageJava})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Run() len(results) = %d, want 1", len(results))
	}
	if ch.called != 0 {
		t.Fatalf("inapplicable checker called %d times, want 0", ch.called)
	}
	if !results[0].Skipped {
		t.Fatal("result should be skipped")
	}
	if !strings.Contains(results[0].SkipReason, "not applicable") {
		t.Fatalf("SkipReason = %q, want contains %q", results[0].SkipReason, "not applicable")
	}
}

func TestRunnerPanicRecovery(t *testing.T) {
	reg := checker.NewDefaultRegistry()
	ch := &runMockChecker{
		id:     "panic-checker",
		pillar: checker.PillarTesting,
		level:  checker.LevelFunctional,
		name:   "Panic Checker",
		checkFn: func(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
			panic("boom")
		},
	}
	if err := reg.Register(ch); err != nil {
		t.Fatalf("register checker: %v", err)
	}

	runner := &checker.Runner{Registry: reg}
	results, err := runner.Run(context.Background(), fstest.MapFS{}, struct{ Language checker.Language }{Language: checker.LanguageGo})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Run() len(results) = %d, want 1", len(results))
	}
	if results[0].Passed {
		t.Fatal("panic result should be failed")
	}
	if !strings.Contains(results[0].Evidence, "panicked") {
		t.Fatalf("Evidence = %q, want contains %q", results[0].Evidence, "panicked")
	}
}

func TestRunnerContextCancellation(t *testing.T) {
	reg := checker.NewDefaultRegistry()
	ctx, cancel := context.WithCancel(context.Background())

	first := &runMockChecker{
		id:     "a",
		pillar: checker.PillarTesting,
		level:  checker.LevelFunctional,
		name:   "First",
		checkFn: func(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
			cancel()
			return &checker.Result{ID: "a", Passed: true}, nil
		},
	}
	second := &runMockChecker{
		id:     "b",
		pillar: checker.PillarTesting,
		level:  checker.LevelFunctional,
		name:   "Second",
	}

	if err := reg.Register(first); err != nil {
		t.Fatalf("register first: %v", err)
	}
	if err := reg.Register(second); err != nil {
		t.Fatalf("register second: %v", err)
	}

	runner := &checker.Runner{Registry: reg}
	results, err := runner.Run(ctx, fstest.MapFS{}, struct{ Language checker.Language }{Language: checker.LanguageGo})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
	if len(results) != 1 {
		t.Fatalf("Run() len(results) = %d, want 1", len(results))
	}
	if first.called != 1 {
		t.Fatalf("first called %d times, want 1", first.called)
	}
	if second.called != 0 {
		t.Fatalf("second called %d times, want 0", second.called)
	}
}
