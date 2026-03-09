package views

import (
	"fmt"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type ProgressModel struct {
	done      int
	total     int
	current   string
	completed []string
	spinner   int
}

func (m ProgressModel) View() string {
	frame := BrightCyan + spinnerFrames[m.spinner%len(spinnerFrames)] + Reset
	bar := m.progressBar(20)
	current := m.current
	if current == "" {
		current = "-"
	}

	recent := "  " + Dim + "(none yet)" + Reset
	if len(m.completed) > 0 {
		recent = "  " + strings.Join(m.completed, "\n  ")
	}

	return fmt.Sprintf(
		"%s\n%s %sScanning repository system...%s\n\n  %sProgress:%s [%s] %s%d/%d%s\n\n  %sCurrent:%s  %s%s%s\n\n  %sRecent:%s\n%s\n",
		CyberHeader,
		frame,
		BrightMagenta, Reset,
		Dim, Reset, bar, BrightCyan, m.done, m.total, Reset,
		Dim, Reset, BrightCyan, current, Reset,
		Dim, Reset,
		recent,
	)
}

func (m *ProgressModel) SetTotal(n int) {
	if n < 0 {
		n = 0
	}
	m.total = n
}

func (m *ProgressModel) SetCurrent(name string) {
	m.current = name
	m.spinner = (m.spinner + 1) % len(spinnerFrames)
}

func (m *ProgressModel) UpdateProgress(result *checker.Result, done, total int) {
	if total >= 0 {
		m.total = total
	}
	if done >= 0 {
		m.done = done
	}

	if result != nil {
		icon := BrightGreen + "✓" + Reset
		if result.Skipped {
			icon = BrightYellow + "↷" + Reset
		} else if !result.Passed {
			icon = BrightRed + "✗" + Reset
		}

		m.completed = append(m.completed, fmt.Sprintf("%s %s%s%s", icon, Dim, result.Name, Reset))
		if len(m.completed) > 5 {
			m.completed = m.completed[len(m.completed)-5:]
		}
		m.current = result.Name
	}

	m.spinner = (m.spinner + 1) % len(spinnerFrames)
}

func (m ProgressModel) progressBar(width int) string {
	if width <= 0 {
		width = 1
	}
	if m.total <= 0 {
		return Dim + strings.Repeat("░", width) + Reset
	}

	ratio := float64(m.done) / float64(m.total)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * float64(width))
	if m.done > 0 && filled == 0 {
		filled = 1
	}
	if filled > width {
		filled = width
	}
	empty := width - filled

	return BrightCyan + strings.Repeat("▓", filled) + Dim + Cyan + strings.Repeat("░", empty) + Reset
}
