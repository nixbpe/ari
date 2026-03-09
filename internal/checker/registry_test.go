package checker_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/bbik/ari/internal/checker"
)

type mockChecker struct {
	id          checker.CheckerID
	pillar      checker.Pillar
	level       checker.Level
	name        string
	description string
	checkFn     func(context.Context, fs.FS, checker.Language) (*checker.Result, error)
}

func (m *mockChecker) ID() checker.CheckerID { return m.id }

func (m *mockChecker) Pillar() checker.Pillar { return m.pillar }

func (m *mockChecker) Level() checker.Level { return m.level }

func (m *mockChecker) Name() string { return m.name }

func (m *mockChecker) Description() string { return m.description }

func (m *mockChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	if m.checkFn != nil {
		return m.checkFn(ctx, repo, lang)
	}
	return &checker.Result{ID: m.id, Passed: true}, nil
}

func TestRegistryRegisterAndGet(t *testing.T) {
	reg := checker.NewDefaultRegistry()

	checkers := []*mockChecker{
		{id: "a", pillar: checker.PillarStyleValidation, level: checker.LevelFunctional, name: "A"},
		{id: "b", pillar: checker.PillarBuildSystem, level: checker.LevelDocumented, name: "B"},
		{id: "c", pillar: checker.PillarTesting, level: checker.LevelStandardized, name: "C"},
	}

	for _, ch := range checkers {
		if err := reg.Register(ch); err != nil {
			t.Fatalf("Register(%s) error = %v", ch.ID(), err)
		}
	}

	for _, ch := range checkers {
		got, ok := reg.Get(ch.ID())
		if !ok {
			t.Fatalf("Get(%s) not found", ch.ID())
		}
		if got.ID() != ch.ID() {
			t.Fatalf("Get(%s).ID() = %s, want %s", ch.ID(), got.ID(), ch.ID())
		}
	}
}

func TestRegistryDuplicateError(t *testing.T) {
	reg := checker.NewDefaultRegistry()
	ch := &mockChecker{id: "dup", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "Dup"}

	if err := reg.Register(ch); err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	if err := reg.Register(ch); err == nil {
		t.Fatal("second Register() error = nil, want duplicate error")
	}
}

func TestRegistryGetByPillar(t *testing.T) {
	reg := checker.NewDefaultRegistry()

	_ = reg.Register(&mockChecker{id: "style-1", pillar: checker.PillarStyleValidation, level: checker.LevelFunctional, name: "style-1"})
	_ = reg.Register(&mockChecker{id: "build-1", pillar: checker.PillarBuildSystem, level: checker.LevelFunctional, name: "build-1"})
	_ = reg.Register(&mockChecker{id: "style-2", pillar: checker.PillarStyleValidation, level: checker.LevelDocumented, name: "style-2"})

	got := reg.GetByPillar(checker.PillarStyleValidation)
	if len(got) != 2 {
		t.Fatalf("GetByPillar() len = %d, want 2", len(got))
	}
	if got[0].ID() != "style-1" || got[1].ID() != "style-2" {
		t.Fatalf("GetByPillar() IDs = [%s, %s], want [style-1, style-2]", got[0].ID(), got[1].ID())
	}
}

func TestRegistryGetByLevel(t *testing.T) {
	reg := checker.NewDefaultRegistry()

	_ = reg.Register(&mockChecker{id: "f1", pillar: checker.PillarStyleValidation, level: checker.LevelFunctional, name: "f1"})
	_ = reg.Register(&mockChecker{id: "d1", pillar: checker.PillarTesting, level: checker.LevelDocumented, name: "d1"})
	_ = reg.Register(&mockChecker{id: "f2", pillar: checker.PillarDocumentation, level: checker.LevelFunctional, name: "f2"})

	got := reg.GetByLevel(checker.LevelFunctional)
	if len(got) != 2 {
		t.Fatalf("GetByLevel() len = %d, want 2", len(got))
	}
	if got[0].ID() != "f1" || got[1].ID() != "f2" {
		t.Fatalf("GetByLevel() IDs = [%s, %s], want [f1, f2]", got[0].ID(), got[1].ID())
	}
}

func TestRegistryAll(t *testing.T) {
	reg := checker.NewDefaultRegistry()

	_ = reg.Register(&mockChecker{id: "c", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "c"})
	_ = reg.Register(&mockChecker{id: "a", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "a"})
	_ = reg.Register(&mockChecker{id: "b", pillar: checker.PillarTesting, level: checker.LevelFunctional, name: "b"})

	all := reg.All()
	if len(all) != 3 {
		t.Fatalf("All() len = %d, want 3", len(all))
	}
	if all[0].ID() != "a" || all[1].ID() != "b" || all[2].ID() != "c" {
		t.Fatalf("All() IDs = [%s, %s, %s], want [a, b, c]", all[0].ID(), all[1].ID(), all[2].ID())
	}
}
