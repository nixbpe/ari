package checker

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	checkers map[CheckerID]Checker
}

func NewDefaultRegistry() *Registry {
	return &Registry{checkers: make(map[CheckerID]Checker)}
}

func (r *Registry) Register(ch Checker) error {
	if ch == nil {
		return fmt.Errorf("checker cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := ch.ID()
	if _, exists := r.checkers[id]; exists {
		return fmt.Errorf("checker %q already registered", id)
	}

	r.checkers[id] = ch
	return nil
}

func (r *Registry) Get(id CheckerID) (Checker, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ch, ok := r.checkers[id]
	return ch, ok
}

func (r *Registry) GetByPillar(pillar Pillar) []Checker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Checker, 0)
	for _, ch := range r.checkers {
		if ch.Pillar() == pillar {
			out = append(out, ch)
		}
	}
	sortCheckersByID(out)
	return out
}

func (r *Registry) GetByLevel(level Level) []Checker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Checker, 0)
	for _, ch := range r.checkers {
		if ch.Level() == level {
			out = append(out, ch)
		}
	}
	sortCheckersByID(out)
	return out
}

func (r *Registry) All() []Checker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Checker, 0, len(r.checkers))
	for _, ch := range r.checkers {
		out = append(out, ch)
	}
	sortCheckersByID(out)
	return out
}

func sortCheckersByID(checkers []Checker) {
	sort.Slice(checkers, func(i, j int) bool {
		return checkers[i].ID() < checkers[j].ID()
	})
}
