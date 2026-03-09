package checktesting

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type UnitTestsExistChecker struct{}

func (c *UnitTestsExistChecker) ID() checker.CheckerID  { return "unit_tests_exist" }
func (c *UnitTestsExistChecker) Pillar() checker.Pillar { return checker.PillarTesting }
func (c *UnitTestsExistChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *UnitTestsExistChecker) Name() string           { return "Unit Tests Exist" }
func (c *UnitTestsExistChecker) Description() string {
	return "Checks that unit test files exist in the repository"
}
func (c *UnitTestsExistChecker) Suggestion() string {
	return "Add unit tests. For Go: create *_test.go files with Test* functions"
}

func (c *UnitTestsExistChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	count, pattern := countTestFiles(repo, lang)

	if count > 0 {
		result.Passed = true
		result.Evidence = fmt.Sprintf("Found %d test files (%s) across the repository", count, pattern)
	} else {
		result.Passed = false
		result.Evidence = "No test files found"
	}

	return result, nil
}

func countTestFiles(repo fs.FS, lang checker.Language) (int, string) {
	var count int
	var pattern string

	switch lang {
	case checker.LanguageGo:
		pattern = "*_test.go"
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && strings.HasSuffix(path, "_test.go") {
				count++
			}
			return nil
		})
	case checker.LanguageTypeScript:
		pattern = "*.test.ts, *.spec.ts, *.test.js, *.spec.js, __tests__/"
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if d.Name() == "__tests__" {
					count++
				}
				return nil
			}
			name := d.Name()
			if strings.HasSuffix(name, ".test.ts") || strings.HasSuffix(name, ".spec.ts") ||
				strings.HasSuffix(name, ".test.js") || strings.HasSuffix(name, ".spec.js") {
				count++
			}
			return nil
		})
	case checker.LanguageJava:
		pattern = "*Test.java, *Spec.java, src/test/"
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				name := d.Name()
				if strings.HasSuffix(name, "Test.java") || strings.HasSuffix(name, "Spec.java") {
					count++
				}
			} else if strings.Contains(path, "src/test") {
				count++
			}
			return nil
		})
	default:
		pattern = "*_test.go, *.test.ts, *.spec.ts, *Test.java"
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				name := d.Name()
				if strings.HasSuffix(name, "_test.go") ||
					strings.HasSuffix(name, ".test.ts") || strings.HasSuffix(name, ".spec.ts") ||
					strings.HasSuffix(name, ".test.js") || strings.HasSuffix(name, ".spec.js") ||
					strings.HasSuffix(name, "Test.java") || strings.HasSuffix(name, "Spec.java") {
					count++
				}
			}
			return nil
		})
	}

	return count, pattern
}
