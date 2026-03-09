package style

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type FormatterChecker struct{}

func (c *FormatterChecker) ID() checker.CheckerID  { return "formatter" }
func (c *FormatterChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *FormatterChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *FormatterChecker) Name() string           { return "Code Formatter" }
func (c *FormatterChecker) Description() string    { return "Checks for code formatter configuration" }
func (c *FormatterChecker) Suggestion() string {
	return "Add a code formatter. For TS: npm install -D prettier"
}

func (c *FormatterChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if lang == checker.LanguageGo {
		result.Passed = true
		result.Evidence = "Go uses built-in gofmt formatter"
		return result, nil
	}

	candidates := formatterCandidates(lang)
	for _, entry := range candidates {
		if _, err := fs.Stat(repo, entry.file); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — formatter %s configured", entry.file, entry.tool)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No formatter configuration found"
	return result, nil
}

type formatterEntry struct {
	file string
	tool string
}

func formatterCandidates(lang checker.Language) []formatterEntry {
	switch lang {
	case checker.LanguageTypeScript:
		return []formatterEntry{
			{".prettierrc", "Prettier"},
			{".prettierrc.js", "Prettier"},
			{".prettierrc.json", "Prettier"},
			{".prettierrc.yml", "Prettier"},
			{"prettier.config.js", "Prettier"},
			{"biome.json", "Biome"},
		}
	case checker.LanguageJava:
		return []formatterEntry{
			{".editorconfig", "EditorConfig"},
		}
	default:
		return nil
	}
}
