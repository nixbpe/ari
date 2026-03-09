package reporter

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/scorer"
)

//go:embed templates/report.html
var reportTemplate string

// HTMLReporter renders a report as a self-contained HTML file.
// The template is embedded at compile time; no external CSS, JS, or fonts
// are referenced — the output is fully self-contained.
type HTMLReporter struct{}

// Report implements Reporter.
func (r *HTMLReporter) Report(ctx context.Context, report *Report, w io.Writer) error {
	funcMap := template.FuncMap{
		// fmtDate formats a time.Time for display.
		"fmtDate": func(t time.Time) string {
			return t.UTC().Format("2006-01-02 15:04:05 UTC")
		},
		// fmtPercent formats a 0..1 float as "75%" for text display.
		"fmtPercent": func(f float64) string {
			return fmt.Sprintf("%.0f%%", f*100)
		},
		// barStyle returns a CSS width declaration safe for use in style attributes.
		// Returns template.CSS so html/template does not sanitize it.
		"barStyle": func(f float64) template.CSS {
			pct := int(f * 100)
			if pct < 0 {
				pct = 0
			}
			if pct > 100 {
				pct = 100
			}
			return template.CSS(fmt.Sprintf("width:%d%%", pct))
		},
		// levelClass maps a level string to a CSS modifier class (l0–l5).
		"levelClass": func(level string) string {
			switch level {
			case "Functional":
				return "l1"
			case "Documented":
				return "l2"
			case "Standardized":
				return "l3"
			case "Optimized":
				return "l4"
			case "Autonomous":
				return "l5"
			default:
				return "l0"
			}
		},
		// levelLabel returns "L3 — Standardized" style labels for the score card.
		"levelLabel": func(level string) string {
			switch level {
			case "Functional":
				return "L1 \u2014 Functional"
			case "Documented":
				return "L2 \u2014 Documented"
			case "Standardized":
				return "L3 \u2014 Standardized"
			case "Optimized":
				return "L4 \u2014 Optimized"
			case "Autonomous":
				return "L5 \u2014 Autonomous"
			default:
				return "L0 \u2014 None"
			}
		},
		// levelName returns the human-readable name of a checker.Level.
		"levelName": func(l checker.Level) string {
			return l.String()
		},
		// sortedLevelScores returns level scores sorted L1→L5.
		"sortedLevelScores": func(m map[checker.Level]scorer.LevelScore) []scorer.LevelScore {
			if m == nil {
				return nil
			}
			ordered := []checker.Level{
				checker.LevelFunctional,
				checker.LevelDocumented,
				checker.LevelStandardized,
				checker.LevelOptimized,
				checker.LevelAutonomous,
			}
			result := make([]scorer.LevelScore, 0, len(ordered))
			for _, l := range ordered {
				if s, ok := m[l]; ok {
					result = append(result, s)
				}
			}
			return result
		},
		// sortedPillarScores returns pillar scores in canonical pillar order.
		"sortedPillarScores": func(m map[checker.Pillar]scorer.PillarScore) []scorer.PillarScore {
			if m == nil {
				return nil
			}
			ordered := []checker.Pillar{
				checker.PillarStyleValidation,
				checker.PillarBuildSystem,
				checker.PillarTesting,
				checker.PillarDocumentation,
			}
			result := make([]scorer.PillarScore, 0, len(ordered))
			for _, p := range ordered {
				if s, ok := m[p]; ok {
					result = append(result, s)
				}
			}
			return result
		},
		// sortedCriteria sorts criteria by pillar name then level number.
		"sortedCriteria": func(criteria []CriterionReport) []CriterionReport {
			if len(criteria) == 0 {
				return criteria
			}
			sorted := make([]CriterionReport, len(criteria))
			copy(sorted, criteria)
			sort.Slice(sorted, func(i, j int) bool {
				if sorted[i].Pillar != sorted[j].Pillar {
					return sorted[i].Pillar < sorted[j].Pillar
				}
				return sorted[i].Level < sorted[j].Level
			})
			return sorted
		},
		// upper converts a string to upper case (used for priority badges).
		"upper": strings.ToUpper,
		// add returns a+b (used for 1-based suggestion numbering).
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("html reporter: parse template: %w", err)
	}

	if err := tmpl.Execute(w, report); err != nil {
		return fmt.Errorf("html reporter: execute template: %w", err)
	}
	return nil
}
