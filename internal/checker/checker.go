package checker

import (
	"context"
	"io/fs"
)

// Language represents the primary programming language of a repository
type Language int

const (
	LanguageUnknown Language = iota
	LanguageGo
	LanguageTypeScript
	LanguageJava
)

func (l Language) String() string {
	switch l {
	case LanguageGo:
		return "Go"
	case LanguageTypeScript:
		return "TypeScript"
	case LanguageJava:
		return "Java"
	default:
		return "Unknown"
	}
}

// Level represents a maturity level (1-5)
type Level int

const (
	LevelFunctional   Level = 1
	LevelDocumented   Level = 2
	LevelStandardized Level = 3
	LevelOptimized    Level = 4
	LevelAutonomous   Level = 5
)

func (l Level) String() string {
	switch l {
	case LevelFunctional:
		return "Functional"
	case LevelDocumented:
		return "Documented"
	case LevelStandardized:
		return "Standardized"
	case LevelOptimized:
		return "Optimized"
	case LevelAutonomous:
		return "Autonomous"
	default:
		return "Unknown"
	}
}

// Pillar represents one of the four evaluation pillars
type Pillar int

const (
	PillarStyleValidation Pillar = iota
	PillarBuildSystem
	PillarTesting
	PillarDocumentation
)

func (p Pillar) String() string {
	switch p {
	case PillarStyleValidation:
		return "Style & Validation"
	case PillarBuildSystem:
		return "Build System"
	case PillarTesting:
		return "Testing"
	case PillarDocumentation:
		return "Documentation"
	default:
		return "Unknown"
	}
}

// CheckerID is a unique identifier for a criterion
type CheckerID string

// Result holds the outcome of a single criterion evaluation
type Result struct {
	ID         CheckerID
	Name       string
	Passed     bool
	Evidence   string
	Level      Level
	Pillar     Pillar
	Skipped    bool
	SkipReason string
	Mode       string // "rule-based" or "llm"
	Suggestion string // what to do if criterion fails
}

// Checker evaluates a single criterion against a repository
type Checker interface {
	ID() CheckerID
	Pillar() Pillar
	Level() Level
	Name() string
	Description() string
	Check(ctx context.Context, repo fs.FS, lang Language) (*Result, error)
}

// SuggestionProvider provides a fix suggestion for a failing criterion
type SuggestionProvider interface {
	Suggestion() string
}
