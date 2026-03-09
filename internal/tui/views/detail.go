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

		// Table Header
		sb.WriteString("  " + BrightMagenta + "STS" + Dim + Cyan + " │ " +
			BrightMagenta + "LVL" + Dim + Cyan + " │ " +
			BrightMagenta + fmt.Sprintf("%-28s", "CRITERION") + Dim + Cyan + " │ " +
			BrightMagenta + "DETAILS\n" + Reset)

		// Table Divider
		sb.WriteString("  " + Dim + Cyan + "────┼─────┼──────────────────────────────┼────────────────────────────────────────────────────────\n" + Reset)

		for _, r := range m.Results[m.Offset:end] {
			if r == nil {
				continue
			}

			lvlColor := LevelColor(r.Level)

			statusRaw := " ✓ "
			statusColor := BrightGreen
			nameColor := BrightCyan
			evColor := Dim + White

			if r.Skipped {
				statusRaw = " ↷ "
				statusColor = BrightYellow
			} else if !r.Passed {
				statusRaw = " ✗ "
				statusColor = BrightRed
				nameColor = BrightRed
				evColor = White
			}

			lvlRaw := fmt.Sprintf("L%d ", r.Level)

			nameRaw := r.Name
			if len(nameRaw) > 28 {
				nameRaw = nameRaw[:27] + "…"
			}

			evidence := r.Evidence
			if r.Skipped {
				reason := r.SkipReason
				if reason == "" {
					reason = "not applicable"
				}
				evidence = "skipped: " + reason
			} else if evidence == "" {
				evidence = "(no evidence)"
			}
			evidence = strings.ReplaceAll(evidence, "\n", " ")

			// Main Row
			sb.WriteString(fmt.Sprintf("  %s%s%s │ %s%s%s │ %s%-28s%s │ %s%s%s\n",
				statusColor, statusRaw, Dim+Cyan,
				lvlColor, lvlRaw, Dim+Cyan,
				nameColor, nameRaw, Dim+Cyan,
				evColor, evidence, Reset,
			))

			// Suggestion Row (if failed)
			if !r.Passed && r.Suggestion != "" {
				sugg := strings.ReplaceAll(r.Suggestion, "\n", " ")
				sb.WriteString(fmt.Sprintf("  %s   %s │ %s   %s │ %s%-28s%s │ %s⚡ %s%s%s\n",
					Dim+Cyan, Dim+Cyan,
					Dim+Cyan, Dim+Cyan,
					Dim+Cyan, "", Dim+Cyan,
					BrightYellow, Dim+Yellow, sugg, Reset,
				))
			}
		}

		// Table Footer / Pagination
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
