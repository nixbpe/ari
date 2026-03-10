package all_test

import (
	"testing"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/all"
)

func TestRegisterAll(t *testing.T) {
	r := checker.NewDefaultRegistry()
	all.RegisterAll(r, nil)
}

func TestRegistryHas72Checkers(t *testing.T) {
	r := checker.NewDefaultRegistry()
	all.RegisterAll(r, nil)
	if len(r.All()) != 72 {
		t.Errorf("expected 72 checkers, got %d", len(r.All()))
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	r := checker.NewDefaultRegistry()
	all.RegisterAll(r, nil)
	seen := make(map[checker.CheckerID]bool)
	for _, ch := range r.All() {
		if seen[ch.ID()] {
			t.Errorf("duplicate checker ID: %s", ch.ID())
		}
		seen[ch.ID()] = true
	}
}
