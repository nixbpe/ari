package style

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DeadCodeDetectionChecker struct{}

func (c *DeadCodeDetectionChecker) ID() checker.CheckerID  { return "dead_code_detection" }
func (c *DeadCodeDetectionChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *DeadCodeDetectionChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *DeadCodeDetectionChecker) Name() string           { return "Dead Code Detection" }
func (c *DeadCodeDetectionChecker) Description() string {
	return "Checks for dead code detection tools in project configuration"
}
func (c *DeadCodeDetectionChecker) Suggestion() string {
	return "Add dead code detection. For Go: add deadcode to golangci-lint. For TS: npm install -D knip"
}

func (c *DeadCodeDetectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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

func (c *DeadCodeDetectionChecker) checkGo(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{".golangci.yml", ".golangci.yaml"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		for _, tool := range []string{"deadcode", "unused", "unparam"} {
			if bytes.Contains(content, []byte(tool)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s — dead code detection configured", tool)
				return result, nil
			}
		}
	}
	result.Passed = false
	result.Evidence = "No dead code detection tools found"
	return result, nil
}

func (c *DeadCodeDetectionChecker) checkTypeScript(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	content, err := fs.ReadFile(repo, "package.json")
	if err != nil {
		result.Passed = false
		result.Evidence = "No dead code detection tools found"
		return result, nil
	}

	var pkg struct {
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(content, &pkg); err == nil {
		for _, tool := range []string{"knip", "ts-prune"} {
			if _, ok := pkg.DevDependencies[tool]; ok {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s — dead code detection configured", tool)
				return result, nil
			}
		}
	}

	result.Passed = false
	result.Evidence = "No dead code detection tools found"
	return result, nil
}

func (c *DeadCodeDetectionChecker) checkJava(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{"build.gradle", "pom.xml", ".spotbugs.xml"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		for _, pattern := range []string{"spotbugs", "SpotBugs", "pmd", "PMD"} {
			if bytes.Contains(content, []byte(pattern)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Found %s — dead code detection configured", pattern)
				return result, nil
			}
		}
	}
	result.Passed = false
	result.Evidence = "No dead code detection tools found"
	return result, nil
}
