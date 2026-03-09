package views

import (
	"fmt"
	"strings"

	"github.com/bbik/ari/internal/checker"
)

const detailPageSize = 15

// DetailModel renders a scrollable list of criteria for a single pillar.
type DetailModel struct {
	Pillar  checker.Pillar
	Results []*checker.Result
	Offset  int
}

// ScrollUp moves the scroll position up by one.
func (m *DetailModel) ScrollUp() {
	if m.Offset > 0 {
		m.Offset--
	}
}

// ScrollDown moves the scroll position down by one.
func (m *DetailModel) ScrollDown() {
	if m.Offset < len(m.Results)-1 {
		m.Offset++
	}
}

// View renders the detail view as a plain string.
func (m DetailModel) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Pillar: %s  (Esc to go back)\n\n", m.Pillar.String()))

	if len(m.Results) == 0 {
		sb.WriteString("  (no criteria for this pillar)\n")
	} else {
		end := m.Offset + detailPageSize
		if end > len(m.Results) {
			end = len(m.Results)
		}

		for _, r := range m.Results[m.Offset:end] {
			if r == nil {
				continue
			}
			if r.Skipped {
				reason := r.SkipReason
				if reason == "" {
					reason = "not applicable"
				}
				sb.WriteString(fmt.Sprintf("  ↷ [L%d] %s — skipped: %s\n", int(r.Level), r.Name, reason))
				continue
			}

			status := "✓"
			if !r.Passed {
				status = "✗"
			}
			evidence := r.Evidence
			if evidence == "" {
				evidence = "(no evidence)"
			}
			sb.WriteString(fmt.Sprintf("  %s [L%d] %s — %s\n", status, int(r.Level), r.Name, evidence))
			if !r.Passed && r.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("      → %s\n", r.Suggestion))
			}
		}

		if len(m.Results) > detailPageSize {
			sb.WriteString(fmt.Sprintf("\n  Showing %d–%d of %d\n", m.Offset+1, end, len(m.Results)))
		}
	}

	sb.WriteString("\n  ↑↓ scroll  Esc back\n")
	return sb.String()
}
