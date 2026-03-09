package views

import (
	"fmt"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scorer"
)

var pillarOrder = []checker.Pillar{
	checker.PillarStyleValidation,
	checker.PillarBuildSystem,
	checker.PillarTesting,
	checker.PillarDocumentation,
}

type ReportModel struct {
	Score          *scorer.Score
	Report         *reporter.Report
	SelectedPillar int
}

func (m *ReportModel) MoveUp() {
	if m.SelectedPillar > 0 {
		m.SelectedPillar--
	}
}

func (m *ReportModel) MoveDown() {
	if m.SelectedPillar < len(pillarOrder)-1 {
		m.SelectedPillar++
	}
}

func (m ReportModel) CurrentPillar() checker.Pillar {
	if m.SelectedPillar < 0 || m.SelectedPillar >= len(pillarOrder) {
		return pillarOrder[0]
	}

	return pillarOrder[m.SelectedPillar]
}

func (m ReportModel) View() string {
	var sb strings.Builder
	sb.WriteString("╔══════════════════════════════════════════╗\n")
	sb.WriteString("║     ARI — Agent Readiness Index          ║\n")
	sb.WriteString("╚══════════════════════════════════════════╝\n\n")

	if m.Score != nil {
		sb.WriteString(fmt.Sprintf("  Level: L%d — %s\n", int(m.Score.Level), m.Score.Level.String()))
		sb.WriteString(fmt.Sprintf("  Pass Rate: %.0f%%\n\n", m.Score.PassRate*100))
	}

	sb.WriteString("  Pillars:\n")
	for i, pillar := range pillarOrder {
		prefix := "  "
		if i == m.SelectedPillar {
			prefix = "> "
		}

		rate := 0.0
		if m.Score != nil {
			if ps, ok := m.Score.PillarScores[pillar]; ok && ps.Total > 0 {
				rate = float64(ps.Passed) / float64(ps.Total)
			}
		}

		bar := progressBar(rate, 10)
		sb.WriteString(fmt.Sprintf("  %s%-22s %s %.0f%%\n", prefix, pillar.String(), bar, rate*100))
	}

	sb.WriteString("\n  ↑↓ navigate  Enter drill-down  h HTML  j JSON  q quit\n")
	return sb.String()
}

func progressBar(pct float64, width int) string {
	if width <= 0 {
		width = 1
	}

	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}

	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}
