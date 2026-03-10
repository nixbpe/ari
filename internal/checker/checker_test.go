package checker_test

import (
	"testing"

	"github.com/nixbpe/ari/internal/checker"
)

func TestLanguageString(t *testing.T) {
	tests := []struct {
		lang checker.Language
		want string
	}{
		{checker.LanguageGo, "Go"},
		{checker.LanguageTypeScript, "TypeScript"},
		{checker.LanguageJava, "Java"},
		{checker.LanguageUnknown, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.lang.String(); got != tt.want {
				t.Errorf("Language.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level checker.Level
		want  string
	}{
		{checker.LevelFunctional, "Functional"},
		{checker.LevelDocumented, "Documented"},
		{checker.LevelStandardized, "Standardized"},
		{checker.LevelOptimized, "Optimized"},
		{checker.LevelAutonomous, "Autonomous"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPillarString(t *testing.T) {
	tests := []struct {
		pillar checker.Pillar
		want   string
	}{
		{checker.PillarContextIntent, "Context & Intent"},
		{checker.PillarEnvInfra, "Environment & Infra"},
		{checker.PillarConstraints, "Constraints & Governance"},
		{checker.PillarVerification, "Verification & Feedback"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.pillar.String(); got != tt.want {
				t.Errorf("Pillar.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResultFields(t *testing.T) {
	r := checker.Result{
		ID:         "lint_config",
		Name:       "Linter Configuration",
		Passed:     true,
		Evidence:   "Found .golangci.yml",
		Level:      checker.LevelFunctional,
		Pillar:     checker.PillarConstraints,
		Mode:       "rule-based",
		Suggestion: "Add a linter config",
	}
	if r.ID != "lint_config" {
		t.Errorf("Result.ID = %q, want %q", r.ID, "lint_config")
	}
	if !r.Passed {
		t.Error("Result.Passed should be true")
	}
}
