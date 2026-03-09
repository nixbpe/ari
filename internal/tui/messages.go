package tui

import (
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scorer"
)

type ScanStartMsg struct{ Total int }

type CheckerStartMsg struct {
	Name string
	ID   checker.CheckerID
}

type CheckerCompleteMsg struct {
	Result *checker.Result
	Done   int
	Total  int
}

type ScanCompleteMsg struct {
	Score   *scorer.Score
	Report  *reporter.Report
	Results []*checker.Result
}

type ErrorMsg struct{ Err error }

type OpenBrowserMsg struct{ Path string }

// DrillDownMsg is sent when the user presses Enter on a pillar in the report view.
type DrillDownMsg struct{ Pillar checker.Pillar }

// BackMsg is sent when the user presses Esc in the detail view.
type BackMsg struct{}

// ExportJSONMsg is sent when the user presses 'j' in the report view.
type ExportJSONMsg struct{}
