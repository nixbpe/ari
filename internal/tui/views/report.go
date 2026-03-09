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
	sb.WriteString(CyberHeader)

	if m.Score != nil {
		lvlColor := LevelColor(m.Score.Level)
		passRate := m.Score.PassRate * 100

		rateColor := BrightGreen
		if passRate < 50 {
			rateColor = BrightRed
		} else if passRate < 80 {
			rateColor = BrightYellow
		}

		sb.WriteString(fmt.Sprintf("  %sLEVEL:%s %sL%d%s — %s%s%s\n",
			Dim, Reset,
			lvlColor+Bold, int(m.Score.Level), Reset,
			lvlColor, m.Score.Level.String(), Reset))

		sb.WriteString(fmt.Sprintf("  %sPASS RATE:%s %s%.0f%%%s\n\n",
			Dim, Reset, rateColor, passRate, Reset))
	}

	sb.WriteString(fmt.Sprintf("  %s>> PILLAR ANALYSIS%s\n", BrightMagenta, Reset))

	for i, pillar := range pillarOrder {
		prefix := "  "
		nameStyle := Dim + White
		if i == m.SelectedPillar {
			prefix = BrightCyan + "▶ " + Reset
			nameStyle = BrightCyan + Bold
		}

		rate := 0.0
		if m.Score != nil {
			if ps, ok := m.Score.PillarScores[pillar]; ok && ps.Total > 0 {
				rate = float64(ps.Passed) / float64(ps.Total)
			}
		}

		bar := progressBar(rate, 15)

		// Ensure exactly the same text but with formatting
		sb.WriteString(fmt.Sprintf("%s %s%-22s%s %s %s%.0f%%%s\n",
			prefix,
			nameStyle, pillar.String(), Reset,
			bar,
			Dim, rate*100, Reset))
	}

	sb.WriteString(fmt.Sprintf("\n  %s[%sARI%s]%s> %s↑↓%s navigate  %sEnter%s drill-down  %sh%s HTML  %sj%s JSON  %sq%s quit\n",
		Dim, BrightCyan, Dim, Reset,
		BrightMagenta, Dim,
		BrightMagenta, Dim,
		BrightMagenta, Dim,
		BrightMagenta, Dim,
		BrightMagenta, Reset))

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
	empty := width - filled

	barColor := BrightGreen
	if pct < 0.5 {
		barColor = BrightRed
	} else if pct < 0.8 {
		barColor = BrightYellow
	}

	return barColor + strings.Repeat("▓", filled) + Dim + Cyan + strings.Repeat("░", empty) + Reset
}
