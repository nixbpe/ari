package checktesting

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type TestNamingConventionsChecker struct{}

func (c *TestNamingConventionsChecker) ID() checker.CheckerID  { return "test_naming_conventions" }
func (c *TestNamingConventionsChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *TestNamingConventionsChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *TestNamingConventionsChecker) Name() string           { return "Test Naming Conventions" }
func (c *TestNamingConventionsChecker) Description() string {
	return "Checks that test files follow language-specific naming conventions"
}
func (c *TestNamingConventionsChecker) Suggestion() string {
	return "Follow test naming conventions. For Go: use *_test.go files with Test* functions"
}

func (c *TestNamingConventionsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	total, conforming := countTestFileConformance(repo, lang)

	if total == 0 {
		result.Skipped = true
		result.SkipReason = "No test files found"
		result.Evidence = "No test files found — skipping naming convention check"
		return result, nil
	}

	if conforming == total {
		result.Passed = true
		result.Evidence = fmt.Sprintf("All %d test files follow %s convention", total, conventionName(lang))
	} else {
		result.Passed = false
		result.Evidence = fmt.Sprintf("Test files don't follow naming convention: %d/%d conform to %s", conforming, total, conventionName(lang))
	}

	return result, nil
}

func countTestFileConformance(repo fs.FS, lang checker.Language) (total, conforming int) {
	switch lang {
	case checker.LanguageGo:
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			name := d.Name()
			if strings.HasSuffix(name, "_test.go") {
				total++
				conforming++
			}
			return nil
		})
	case checker.LanguageTypeScript:
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if d.Name() == "__tests__" {
					total++
					conforming++
				}
				return nil
			}
			name := d.Name()
			if strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".js") {
				if strings.HasSuffix(name, ".test.ts") || strings.HasSuffix(name, ".spec.ts") ||
					strings.HasSuffix(name, ".test.js") || strings.HasSuffix(name, ".spec.js") {
					total++
					conforming++
				}
			}
			return nil
		})
	case checker.LanguageJava:
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			name := d.Name()
			if strings.HasSuffix(name, ".java") && strings.Contains(strings.ToLower(path), "test") {
				total++
				if strings.HasSuffix(name, "Test.java") || strings.HasSuffix(name, "Spec.java") {
					conforming++
				}
			}
			return nil
		})
	}
	return total, conforming
}

func conventionName(lang checker.Language) string {
	switch lang {
	case checker.LanguageGo:
		return "Go *_test.go"
	case checker.LanguageTypeScript:
		return "TypeScript *.test.ts / *.spec.ts"
	case checker.LanguageJava:
		return "Java *Test.java / *Spec.java"
	default:
		return "standard"
	}
}
