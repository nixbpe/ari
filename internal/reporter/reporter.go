package reporter

import (
	"context"
	"io"
	"time"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/scanner"
	"github.com/bbik/ari/internal/scorer"
)

// Format represents the output format
type Format int

const (
	FormatJSON Format = iota
	FormatHTML
	FormatText
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatHTML:
		return "html"
	case FormatText:
		return "text"
	default:
		return "unknown"
	}
}

// CriterionReport holds the result of evaluating a single criterion.
type CriterionReport struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Pillar     string `json:"pillar"`
	Level      int    `json:"level"`
	Passed     bool   `json:"passed"`
	Skipped    bool   `json:"skipped"`
	SkipReason string `json:"skipReason,omitempty"`
	Evidence   string `json:"evidence,omitempty"`
	Mode       string `json:"mode"`
	Suggestion string `json:"suggestion,omitempty"`
}

// Suggestion holds a fix suggestion for a failing criterion.
type Suggestion struct {
	CriterionID string `json:"criterionId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Report is the central data model consumed by all output formatters.
type Report struct {
	RepoPath        string            `json:"repoPath"`
	Language        string            `json:"language"`
	Level           string            `json:"level"`
	PassRate        float64           `json:"passRate"`
	IsGitRepo       bool              `json:"isGitRepo"`
	CommitHash      string            `json:"commitHash"`
	Branch          string            `json:"branch,omitempty"`
	Score           *scorer.Score     `json:"score,omitempty"`
	CriteriaResults []CriterionReport `json:"criteria"`
	Suggestions     []Suggestion      `json:"suggestions"`
	GeneratedAt     time.Time         `json:"generatedAt"`
	AriVersion      string            `json:"ariVersion"`
}

// Reporter generates a report by writing to an io.Writer.
type Reporter interface {
	Report(ctx context.Context, report *Report, w io.Writer) error
}

// BuildReport assembles a Report from raw scanner, scorer, and checker data.
// Suggestions are generated only for failing, non-skipped criteria.
func BuildReport(repoInfo *scanner.RepoInfo, score *scorer.Score, results []*checker.Result) *Report {
	r := &Report{
		RepoPath:        repoInfo.RootPath,
		Language:        repoInfo.Language.String(),
		IsGitRepo:       repoInfo.IsGitRepo,
		CommitHash:      repoInfo.CommitHash,
		Branch:          repoInfo.Branch,
		Score:           score,
		GeneratedAt:     time.Now(),
		AriVersion:      "0.1.0",
		CriteriaResults: []CriterionReport{},
		Suggestions:     []Suggestion{},
	}

	if score != nil {
		r.Level = score.Level.String()
		r.PassRate = score.PassRate
	}

	for _, result := range results {
		cr := CriterionReport{
			ID:         string(result.ID),
			Name:       result.Name,
			Pillar:     result.Pillar.String(),
			Level:      int(result.Level),
			Passed:     result.Passed,
			Skipped:    result.Skipped,
			SkipReason: result.SkipReason,
			Evidence:   result.Evidence,
			Mode:       result.Mode,
			Suggestion: result.Suggestion,
		}
		r.CriteriaResults = append(r.CriteriaResults, cr)

		// Only generate suggestions for failing, non-skipped criteria.
		if !result.Passed && !result.Skipped {
			desc := result.Suggestion
			if desc == "" {
				desc = "No specific suggestion available."
			}
			r.Suggestions = append(r.Suggestions, Suggestion{
				CriterionID: string(result.ID),
				Title:       result.Name,
				Description: desc,
				Priority:    priorityForLevel(result.Level),
			})
		}
	}

	return r
}

func priorityForLevel(level checker.Level) string {
	switch level {
	case checker.LevelFunctional, checker.LevelDocumented:
		return "high"
	case checker.LevelStandardized, checker.LevelOptimized:
		return "medium"
	case checker.LevelAutonomous:
		return "low"
	default:
		return "medium"
	}
}
