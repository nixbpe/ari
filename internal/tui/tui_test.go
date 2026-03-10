package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/scorer"
)

func TestModelInitialView(t *testing.T) {
	model := NewModel()
	if model.currentView != ProgressView {
		t.Fatalf("expected initial view %v, got %v", ProgressView, model.currentView)
	}

	view := model.View()
	if view.Content == "" {
		t.Fatal("expected non-empty initial view content")
	}
}

func TestProgressToReport(t *testing.T) {
	model := NewModel()
	next, _ := model.Update(ScanCompleteMsg{})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if updated.currentView != ReportView {
		t.Fatalf("expected view %v, got %v", ReportView, updated.currentView)
	}
}

func TestProgressUpdates(t *testing.T) {
	model := NewModel()
	result := &checker.Result{Name: "lint_config", Passed: true}

	next, _ := model.Update(CheckerCompleteMsg{Result: result, Done: 1, Total: 2})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}
	content := updated.View().Content

	if content == "" {
		t.Fatal("expected non-empty view content")
	}
	if !strings.Contains(content, "1/2") {
		t.Fatalf("expected progress to contain 1/2, got %q", content)
	}
	if !strings.Contains(content, "lint_config") {
		t.Fatalf("expected recent log to contain checker name, got %q", content)
	}
}

func TestQuitMsg(t *testing.T) {
	model := NewModel()
	msg := tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl})

	next, cmd := model.Update(msg)
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if !updated.quitting {
		t.Fatal("expected model to be in quitting state")
	}
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestReportViewPillars(t *testing.T) {
	score := &scorer.Score{
		Level:    checker.LevelDocumented,
		PassRate: 0.73,
		PillarScores: map[checker.Pillar]scorer.PillarScore{
			checker.PillarConstraints:   {Passed: 7, Total: 10, Rate: 0.7},
			checker.PillarEnvInfra:      {Passed: 8, Total: 10, Rate: 0.8},
			checker.PillarVerification:  {Passed: 6, Total: 10, Rate: 0.6},
			checker.PillarContextIntent: {Passed: 9, Total: 10, Rate: 0.9},
		},
	}
	model := NewModel()
	next, _ := model.Update(ScanCompleteMsg{Score: score})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if updated.currentView != ReportView {
		t.Fatalf("expected ReportView, got %v", updated.currentView)
	}

	content := updated.View().Content
	for _, name := range []string{"Context & Intent", "Environment & Infra", "Constraints & Governance", "Verification & Feedback"} {
		if !strings.Contains(content, name) {
			t.Errorf("View() missing pillar %q\ncontent:\n%s", name, content)
		}
	}
	if !strings.Contains(content, "73%") {
		t.Errorf("View() missing pass rate 73%%\ncontent:\n%s", content)
	}
}

func TestDetailViewCriteria(t *testing.T) {
	results := []*checker.Result{
		{
			ID:       "lint_config",
			Name:     "lint_config",
			Passed:   true,
			Evidence: "golangci.yml found",
			Level:    checker.LevelFunctional,
			Pillar:   checker.PillarConstraints,
		},
		{
			ID:         "unit_tests_exist",
			Name:       "unit_tests_exist",
			Passed:     false,
			Evidence:   "no test files found",
			Level:      checker.LevelFunctional,
			Pillar:     checker.PillarConstraints,
			Suggestion: "Add *_test.go files",
		},
	}

	model := NewModel()
	// Transition to report view with results
	next, _ := model.Update(ScanCompleteMsg{Results: results})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// Drill down into Style & Validation
	next2, _ := updated.Update(DrillDownMsg{Pillar: checker.PillarConstraints})
	drilled, ok := next2.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if drilled.currentView != DetailView {
		t.Fatalf("expected DetailView, got %v", drilled.currentView)
	}

	content := drilled.View().Content
	if !strings.Contains(content, "lint_config") {
		t.Errorf("detail view missing lint_config\ncontent:\n%s", content)
	}
	if !strings.Contains(content, "unit_tests_exist") {
		t.Errorf("detail view missing unit_tests_exist\ncontent:\n%s", content)
	}
	if !strings.Contains(content, "✓") {
		t.Errorf("detail view missing pass indicator ✓\ncontent:\n%s", content)
	}
	if !strings.Contains(content, "✗") {
		t.Errorf("detail view missing fail indicator ✗\ncontent:\n%s", content)
	}
}

func TestReportNavigation(t *testing.T) {
	score := &scorer.Score{
		Level:    checker.LevelFunctional,
		PassRate: 0.5,
		PillarScores: map[checker.Pillar]scorer.PillarScore{
			checker.PillarConstraints: {Passed: 5, Total: 10, Rate: 0.5},
		},
	}
	model := NewModel()
	next, _ := model.Update(ScanCompleteMsg{Score: score})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// Initial selected pillar should be 0
	if updated.report.SelectedPillar != 0 {
		t.Fatalf("expected SelectedPillar=0, got %d", updated.report.SelectedPillar)
	}

	// Press down
	next2, _ := updated.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	updated2, ok2 := next2.(Model)
	if !ok2 {
		t.Fatal("expected Model type")
	}
	if updated2.report.SelectedPillar != 1 {
		t.Fatalf("expected SelectedPillar=1 after down, got %d", updated2.report.SelectedPillar)
	}

	// Press up
	next3, _ := updated2.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyUp}))
	updated3, ok3 := next3.(Model)
	if !ok3 {
		t.Fatal("expected Model type")
	}
	if updated3.report.SelectedPillar != 0 {
		t.Fatalf("expected SelectedPillar=0 after up, got %d", updated3.report.SelectedPillar)
	}
}

func TestDetailBack(t *testing.T) {
	model := NewModel()
	next, _ := model.Update(ScanCompleteMsg{})
	updated, ok := next.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// Drill down
	next2, _ := updated.Update(DrillDownMsg{Pillar: checker.PillarVerification})
	drilled, ok2 := next2.(Model)
	if !ok2 {
		t.Fatal("expected Model type")
	}
	if drilled.currentView != DetailView {
		t.Fatalf("expected DetailView, got %v", drilled.currentView)
	}

	// Press Esc — should return BackMsg cmd
	_, cmd := drilled.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape}))
	if cmd == nil {
		t.Fatal("expected a cmd from Esc key in DetailView")
	}
	// Execute the cmd to get the BackMsg
	backMsg := cmd()
	if _, ok := backMsg.(BackMsg); !ok {
		t.Fatalf("expected BackMsg, got %T", backMsg)
	}
}
