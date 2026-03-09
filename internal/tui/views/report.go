package views

import (
	"github.com/bbik/ari/internal/reporter"
	"github.com/bbik/ari/internal/scorer"
)

type ReportModel struct {
	Score  *scorer.Score
	Report *reporter.Report
}

func (m ReportModel) View() string { return "" }

type DetailModel struct{}

func (m DetailModel) View() string { return "" }
