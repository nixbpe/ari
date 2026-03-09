package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/bbik/ari/internal/checker"
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
	updated := next.(Model)

	if updated.currentView != ReportView {
		t.Fatalf("expected view %v, got %v", ReportView, updated.currentView)
	}
}

func TestProgressUpdates(t *testing.T) {
	model := NewModel()
	result := &checker.Result{Name: "lint_config", Passed: true}

	next, _ := model.Update(CheckerCompleteMsg{Result: result, Done: 1, Total: 2})
	updated := next.(Model)
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
	msg := tea.KeyPressMsg(tea.Key{Text: "c", Code: 'c', Mod: tea.ModCtrl})

	next, cmd := model.Update(msg)
	updated := next.(Model)

	if !updated.quitting {
		t.Fatal("expected model to be in quitting state")
	}
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}
