package style

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type StrictTypingChecker struct{}

func NewStrictTypingChecker() *StrictTypingChecker {
	return &StrictTypingChecker{}
}

func (c *StrictTypingChecker) ID() checker.CheckerID  { return "strict_typing" }
func (c *StrictTypingChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *StrictTypingChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *StrictTypingChecker) Name() string           { return "Strict Type Enforcement" }
func (c *StrictTypingChecker) Description() string {
	return "Ensures strict type checking is enabled for the project language"
}

const strictTypingSuggestion = "Enable strict mode. For TypeScript: set strict: true in tsconfig.json compilerOptions"

func (c *StrictTypingChecker) Check(_ context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	switch lang {
	case checker.LanguageGo:
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     true,
			Evidence:   "Go enforces strict typing by default",
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil

	case checker.LanguageJava:
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     true,
			Evidence:   "Java enforces strict typing by default",
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil

	case checker.LanguageTypeScript:
		return c.checkTypeScript(repo)

	default:
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     false,
			Evidence:   fmt.Sprintf("strict typing check not supported for language %s", lang),
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Skipped:    true,
			SkipReason: "unsupported language",
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil
	}
}

type tsconfigJSON struct {
	CompilerOptions struct {
		Strict bool `json:"strict"`
	} `json:"compilerOptions"`
}

func (c *StrictTypingChecker) checkTypeScript(repo fs.FS) (*checker.Result, error) {
	data, err := fs.ReadFile(repo, "tsconfig.json")
	if err != nil {
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     false,
			Evidence:   "tsconfig.json not found",
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil
	}

	var tsconfig tsconfigJSON
	if err := json.Unmarshal(data, &tsconfig); err != nil {
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     false,
			Evidence:   "tsconfig.json could not be parsed",
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil
	}

	if tsconfig.CompilerOptions.Strict {
		return &checker.Result{
			ID:         c.ID(),
			Name:       c.Name(),
			Passed:     true,
			Evidence:   "tsconfig.json has strict: true",
			Level:      c.Level(),
			Pillar:     c.Pillar(),
			Mode:       "rule-based",
			Suggestion: strictTypingSuggestion,
		}, nil
	}

	return &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Passed:     false,
		Evidence:   "tsconfig.json missing strict: true in compilerOptions",
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: strictTypingSuggestion,
	}, nil
}
