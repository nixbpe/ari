package reporter_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scorer"
)

// makeTestReport returns a fully-populated Report suitable for HTML rendering tests.
func makeTestReport() *reporter.Report {
	return &reporter.Report{
		RepoPath:    "/test/my-repo",
		Language:    "Go",
		Level:       checker.LevelStandardized.String(), // "Standardized"
		PassRate:    0.85,
		AriVersion:  "0.1.0",
		GeneratedAt: time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC),
		Score: &scorer.Score{
			Level:    checker.LevelStandardized,
			PassRate: 0.85,
			PillarScores: map[checker.Pillar]scorer.PillarScore{
				checker.PillarStyleValidation: {
					Pillar: checker.PillarStyleValidation,
					Passed: 5,
					Total:  6,
					Rate:   5.0 / 6.0,
				},
				checker.PillarBuildSystem: {
					Pillar: checker.PillarBuildSystem,
					Passed: 4,
					Total:  5,
					Rate:   0.8,
				},
				checker.PillarTesting: {
					Pillar: checker.PillarTesting,
					Passed: 3,
					Total:  4,
					Rate:   0.75,
				},
				checker.PillarDocumentation: {
					Pillar: checker.PillarDocumentation,
					Passed: 5,
					Total:  5,
					Rate:   1.0,
				},
			},
			LevelScores: map[checker.Level]scorer.LevelScore{
				checker.LevelFunctional: {
					Level:    checker.LevelFunctional,
					Passed:   7,
					Total:    7,
					Rate:     1.0,
					Achieved: true,
				},
				checker.LevelDocumented: {
					Level:    checker.LevelDocumented,
					Passed:   7,
					Total:    8,
					Rate:     0.875,
					Achieved: true,
				},
				checker.LevelStandardized: {
					Level:    checker.LevelStandardized,
					Passed:   3,
					Total:    9,
					Rate:     0.333,
					Achieved: false,
				},
				checker.LevelOptimized: {
					Level:    checker.LevelOptimized,
					Passed:   0,
					Total:    12,
					Rate:     0.0,
					Achieved: false,
				},
				checker.LevelAutonomous: {
					Level:    checker.LevelAutonomous,
					Passed:   0,
					Total:    4,
					Rate:     0.0,
					Achieved: false,
				},
			},
		},
		CriteriaResults: []reporter.CriterionReport{
			{
				ID:       "style-001",
				Name:     "Formatter configured",
				Pillar:   checker.PillarStyleValidation.String(),
				Level:    1,
				Passed:   true,
				Evidence: "prettier config found",
			},
			{
				ID:       "build-001",
				Name:     "Go modules present",
				Pillar:   checker.PillarBuildSystem.String(),
				Level:    1,
				Passed:   true,
				Evidence: "go.mod found",
			},
			{
				ID:       "test-001",
				Name:     "Unit tests exist",
				Pillar:   checker.PillarTesting.String(),
				Level:    1,
				Passed:   true,
				Evidence: "found *_test.go files",
			},
			{
				ID:       "docs-001",
				Name:     "README present",
				Pillar:   checker.PillarDocumentation.String(),
				Level:    1,
				Passed:   true,
				Evidence: "README.md found",
			},
			{
				ID:         "build-002",
				Name:       "CI configuration present",
				Pillar:     checker.PillarBuildSystem.String(),
				Level:      2,
				Passed:     false,
				Evidence:   "no .github/workflows",
				Suggestion: "Add GitHub Actions workflow",
			},
			{
				ID:         "test-002",
				Name:       "Test coverage tracked",
				Pillar:     checker.PillarTesting.String(),
				Level:      2,
				Passed:     false,
				Evidence:   "no coverage config",
				Suggestion: "Add coverage reporting",
			},
		},
		Suggestions: []reporter.Suggestion{
			{
				CriterionID: "build-002",
				Title:       "CI configuration present",
				Description: "Add GitHub Actions workflow",
				Priority:    "high",
			},
			{
				CriterionID: "test-002",
				Title:       "Test coverage tracked",
				Description: "Add coverage reporting",
				Priority:    "high",
			},
		},
	}
}

// renderHTML is a helper that runs the HTML reporter and returns the output string.
func renderHTML(t *testing.T, r *reporter.Report) string {
	t.Helper()
	hr := &reporter.HTMLReporter{}
	var buf bytes.Buffer
	if err := hr.Report(context.Background(), r, &buf); err != nil {
		t.Fatalf("HTMLReporter.Report() unexpected error: %v", err)
	}
	return buf.String()
}

// TestHTMLReporterSelfContained verifies the output is valid HTML with no
// external CSS, JS or font references (href/src must not contain http(s)://).
func TestHTMLReporterSelfContained(t *testing.T) {
	html := renderHTML(t, makeTestReport())

	// Basic HTML structure check — no external parser needed.
	if !strings.Contains(html, "<html") {
		t.Error("output does not contain <html — not valid HTML output")
	}
	if !strings.Contains(html, "</html>") {
		t.Error("output does not contain </html>")
	}
	if !strings.Contains(html, "<head") {
		t.Error("output does not contain <head")
	}
	if !strings.Contains(html, "<body") {
		t.Error("output does not contain <body")
	}

	// No external links in href or src attributes.
	for _, ext := range []string{`href="http://`, `href="https://`, `src="http://`, `src="https://`} {
		if strings.Contains(html, ext) {
			t.Errorf("HTML contains external reference %q — report must be self-contained", ext)
		}
	}

	// No external stylesheets or scripts via link/script tags with URLs.
	if strings.Contains(html, "http://") || strings.Contains(html, "https://") {
		t.Error("HTML contains http:// or https:// — should be fully self-contained")
	}
}

// TestHTMLReporterLevel verifies that a Level 3 report shows "Standardized" in the HTML.
func TestHTMLReporterLevel(t *testing.T) {
	r := makeTestReport()
	// makeTestReport already sets Level to checker.LevelStandardized.String() == "Standardized"
	if r.Level != "Standardized" {
		t.Fatalf("test setup error: expected Level='Standardized', got %q", r.Level)
	}

	html := renderHTML(t, r)

	if !strings.Contains(html, "Standardized") {
		t.Error("HTML does not contain 'Standardized' — expected level label in output")
	}
	// The level badge should show the formatted label.
	if !strings.Contains(html, "L3") {
		t.Error("HTML does not contain 'L3' — expected level number in badge")
	}
}

// TestHTMLReporterPillars verifies all 4 pillar names appear in the HTML output.
func TestHTMLReporterPillars(t *testing.T) {
	html := renderHTML(t, makeTestReport())

	// "Style & Validation" is HTML-escaped to "Style &amp; Validation" by html/template.
	pillars := []string{
		"Style &amp; Validation",
		"Build System",
		"Testing",
		"Documentation",
	}
	for _, p := range pillars {
		if !strings.Contains(html, p) {
			t.Errorf("HTML does not contain pillar name %q", p)
		}
	}
}

// TestHTMLReporterEmptyCriteria verifies no crash and valid HTML when criteria list is empty.
func TestHTMLReporterEmptyCriteria(t *testing.T) {
	r := &reporter.Report{
		RepoPath:        "/empty/repo",
		Language:        "Go",
		Level:           "",
		PassRate:        0.0,
		AriVersion:      "0.1.0",
		GeneratedAt:     time.Now(),
		Score:           nil, // no score data
		CriteriaResults: []reporter.CriterionReport{},
		Suggestions:     []reporter.Suggestion{},
	}

	html := renderHTML(t, r)

	if !strings.Contains(html, "<html") {
		t.Error("output with empty criteria does not contain <html")
	}
	// Empty-state messages should appear.
	if !strings.Contains(html, "No criteria evaluated") {
		t.Error("HTML does not show empty criteria message")
	}
	if !strings.Contains(html, "No suggestions") {
		t.Error("HTML does not show empty suggestions message")
	}
}

// TestHTMLReporterSuggestions verifies that a report with 2 failing criteria
// produces HTML containing exactly 2 suggestion entries.
func TestHTMLReporterSuggestions(t *testing.T) {
	r := &reporter.Report{
		RepoPath:    "/test/repo",
		Language:    "Go",
		Level:       "Functional",
		PassRate:    0.0,
		AriVersion:  "0.1.0",
		GeneratedAt: time.Now(),
		CriteriaResults: []reporter.CriterionReport{
			{
				ID:       "fail-001",
				Name:     "Missing CI workflow",
				Pillar:   checker.PillarBuildSystem.String(),
				Level:    1,
				Passed:   false,
				Evidence: "no .github/workflows directory",
			},
			{
				ID:       "fail-002",
				Name:     "Missing unit tests",
				Pillar:   checker.PillarTesting.String(),
				Level:    1,
				Passed:   false,
				Evidence: "no *_test.go files found",
			},
		},
		Suggestions: []reporter.Suggestion{
			{
				CriterionID: "fail-001",
				Title:       "Missing CI workflow",
				Description: "Add a .github/workflows/ci.yml file",
				Priority:    "high",
			},
			{
				CriterionID: "fail-002",
				Title:       "Missing unit tests",
				Description: "Create *_test.go files with unit tests",
				Priority:    "high",
			},
		},
	}

	html := renderHTML(t, r)

	// Count suggestion div elements by their class attribute prefix.
	count := strings.Count(html, `class="suggestion `)
	if count != 2 {
		t.Errorf("expected 2 suggestion divs in HTML, got %d\nHTML snippet:\n%s",
			count, html[strings.Index(html, "Suggestions"):min(len(html), strings.Index(html, "Suggestions")+2000)])
	}

	// Specific suggestion content should appear.
	if !strings.Contains(html, "Missing CI workflow") {
		t.Error("HTML does not contain first suggestion title")
	}
	if !strings.Contains(html, "Missing unit tests") {
		t.Error("HTML does not contain second suggestion title")
	}
}
