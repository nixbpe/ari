package style

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type CyclomaticComplexityChecker struct{}

func (c *CyclomaticComplexityChecker) ID() checker.CheckerID  { return "cyclomatic_complexity" }
func (c *CyclomaticComplexityChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *CyclomaticComplexityChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *CyclomaticComplexityChecker) Name() string           { return "Cyclomatic Complexity Analysis" }
func (c *CyclomaticComplexityChecker) Description() string {
	return "Checks for cyclomatic complexity analysis tools in project configuration"
}
func (c *CyclomaticComplexityChecker) Suggestion() string {
	return "Add complexity analysis. For Go: add gocyclo to golangci-lint. For TS: enable complexity ESLint rule"
}

func (c *CyclomaticComplexityChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		return c.checkGo(repo, result)
	case checker.LanguageTypeScript:
		return c.checkTypeScript(repo, result)
	case checker.LanguageJava:
		return c.checkJava(repo, result)
	default:
		result.Skipped = true
		result.SkipReason = "unsupported language"
		return result, nil
	}
}

func (c *CyclomaticComplexityChecker) checkGo(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{".golangci.yml", ".golangci.yaml"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		for _, tool := range []string{"gocyclo", "gocognit", "cyclop"} {
			if bytes.Contains(content, []byte(tool)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s in %s linter configuration", tool, configFile)
				return result, nil
			}
		}
	}
	result.Passed = false
	result.Evidence = "No complexity analysis tools configured"
	return result, nil
}

func (c *CyclomaticComplexityChecker) checkTypeScript(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{".eslintrc.json", "eslint.config.js"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		for _, pattern := range []string{"complexity", "sonarjs/cognitive-complexity"} {
			if bytes.Contains(content, []byte(pattern)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s rule in %s", pattern, configFile)
				return result, nil
			}
		}
	}
	result.Passed = false
	result.Evidence = "No complexity analysis tools configured"
	return result, nil
}

func (c *CyclomaticComplexityChecker) checkJava(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{"checkstyle.xml", "pmd.xml"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		for _, pattern := range []string{"CyclomaticComplexity", "NPathComplexity"} {
			if bytes.Contains(content, []byte(pattern)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s check in %s", pattern, configFile)
				return result, nil
			}
		}
	}
	result.Passed = false
	result.Evidence = "No complexity analysis tools configured"
	return result, nil
}
