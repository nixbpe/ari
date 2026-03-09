package tui

import (
	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/reporter"
	"github.com/bbik/ari/internal/scorer"
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
	Score  *scorer.Score
	Report *reporter.Report
}

type ErrorMsg struct{ Err error }

type OpenBrowserMsg struct{ Path string }
