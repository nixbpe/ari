package views

import (
	"fmt"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
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
	sb.WriteString(CyberHeader)

	sb.WriteString(fmt.Sprintf("  %s>> PILLAR DETAILS:%s %s%s%s  %s(%sEsc%s to go back)%s\n\n",
		BrightMagenta, Reset,
		BrightCyan+Bold, m.Pillar.String(), Reset,
		Dim, White, Dim, Reset))

	if len(m.Results) == 0 {
		sb.WriteString(fmt.Sprintf("  %s(no criteria for this pillar)%s\n", Dim, Reset))
	} else {
		end := m.Offset + detailPageSize
		if end > len(m.Results) {
			end = len(m.Results)
		}

		for _, r := range m.Results[m.Offset:end] {
			if r == nil {
				continue
			}

			lvlColor := LevelColor(r.Level)

			if r.Skipped {
				reason := r.SkipReason
				if reason == "" {
					reason = "not applicable"
				}
				sb.WriteString(fmt.Sprintf("  %s↷%s [%sL%d%s] %s%s%s — %sskipped: %s%s\n",
					BrightYellow, Reset,
					lvlColor, int(r.Level), Reset,
					Dim, r.Name, Reset,
					Dim, reason, Reset))
				continue
			}

			status := BrightGreen + "✓" + Reset
			nameStyle := BrightCyan

			if !r.Passed {
				status = BrightRed + "✗" + Reset
				nameStyle = BrightRed
			}
			evidence := r.Evidence
			if evidence == "" {
				evidence = "(no evidence)"
			}

			sb.WriteString(fmt.Sprintf("  %s [%sL%d%s] %s%s%s — %s%s%s\n",
				status,
				lvlColor, int(r.Level), Reset,
				nameStyle, r.Name, Reset,
				Dim, evidence, Reset))

			if !r.Passed && r.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("      %s⚡ %s%s\n",
					BrightYellow, Dim+Yellow, r.Suggestion, Reset))
			}
		}

		if len(m.Results) > detailPageSize {
			sb.WriteString(fmt.Sprintf("\n  %sShowing %s%d–%d%s of %s%d%s\n",
				Dim,
				BrightCyan, m.Offset+1, end, Dim,
				BrightCyan, len(m.Results), Reset))
		}
	}

	sb.WriteString(fmt.Sprintf("\n  %s[%sARI%s]%s> %s↑↓%s scroll  %sEsc%s back\n",
		Dim, BrightCyan, Dim, Reset,
		BrightMagenta, Dim,
		BrightMagenta, Reset))

	return sb.String()
}
