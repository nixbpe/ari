package views

import (
	"fmt"
	"strings"

	"github.com/bbik/ari/internal/checker"
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
	frame := spinnerFrames[m.spinner%len(spinnerFrames)]
	bar := m.progressBar(16)
	current := m.current
	if current == "" {
		current = "-"
	}

	recent := "  (none yet)"
	if len(m.completed) > 0 {
		recent = "  " + strings.Join(m.completed, "\n  ")
	}

	return fmt.Sprintf(
		"%s Scanning repository...\n\nProgress: [%s] %d/%d\n\nCurrent: %s\n\nRecent:\n%s\n",
		frame,
		bar,
		m.done,
		m.total,
		current,
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
		icon := "✓"
		if result.Skipped {
			icon = "↷"
		} else if !result.Passed {
			icon = "✗"
		}

		m.completed = append(m.completed, fmt.Sprintf("%s %s", icon, result.Name))
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
		return strings.Repeat("░", width)
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

	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}
