package style

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type TypeCheckChecker struct{}

func (c *TypeCheckChecker) ID() checker.CheckerID  { return "type_check" }
func (c *TypeCheckChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *TypeCheckChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *TypeCheckChecker) Name() string           { return "Type Checking" }
func (c *TypeCheckChecker) Description() string {
	return "Checks for static type checking configuration"
}
func (c *TypeCheckChecker) Suggestion() string {
	return "Add TypeScript. Create tsconfig.json and install typescript"
}

func (c *TypeCheckChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	switch lang {
	case checker.LanguageGo:
		result.Passed = true
		result.Evidence = "Go is statically typed"
		return result, nil

	case checker.LanguageJava:
		result.Passed = true
		result.Evidence = "Java is statically typed"
		return result, nil

	case checker.LanguageTypeScript:
		if _, err := fs.Stat(repo, "tsconfig.json"); err == nil {
			result.Passed = true
			result.Evidence = "Found tsconfig.json — TypeScript configured"
			return result, nil
		}
		result.Passed = false
		result.Evidence = "JavaScript without TypeScript — no type checking"
		return result, nil

	default:
		result.Passed = false
		result.Evidence = "JavaScript without TypeScript — no type checking"
		return result, nil
	}
}
