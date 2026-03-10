package style

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DuplicateCodeDetectionChecker struct{}

func (c *DuplicateCodeDetectionChecker) ID() checker.CheckerID  { return "duplicate_code_detection" }
func (c *DuplicateCodeDetectionChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *DuplicateCodeDetectionChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *DuplicateCodeDetectionChecker) Name() string           { return "Duplicate Code Detection" }
func (c *DuplicateCodeDetectionChecker) Description() string {
	return "Checks for duplicate code detection tools in project configuration"
}
func (c *DuplicateCodeDetectionChecker) Suggestion() string {
	return "Add duplicate code detection. For Go: add dupl to golangci-lint. For all: use jscpd"
}

func (c *DuplicateCodeDetectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if _, err := fs.Stat(repo, ".jscpd.json"); err == nil {
		result.Passed = true
		result.Evidence = "Found .jscpd.json — jscpd duplicate detection configured"
		return result, nil
	}

	if content, err := fs.ReadFile(repo, "package.json"); err == nil {
		var pkg struct {
			Scripts map[string]string `json:"scripts"`
		}
		if err := json.Unmarshal(content, &pkg); err == nil {
			for _, script := range pkg.Scripts {
				if bytes.Contains([]byte(script), []byte("jscpd")) {
					result.Passed = true
					result.Evidence = "Found jscpd in package.json scripts — duplicate code detection configured"
					return result, nil
				}
			}
		}
	}

	switch lang {
	case checker.LanguageGo:
		return c.checkGo(repo, result)
	case checker.LanguageJava:
		return c.checkJava(repo, result)
	}

	result.Passed = false
	result.Evidence = "No duplicate code detection tools found"
	return result, nil
}

func (c *DuplicateCodeDetectionChecker) checkGo(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{".golangci.yml", ".golangci.yaml"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("dupl")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found dupl linter in %s", configFile)
			return result, nil
		}
	}
	result.Passed = false
	result.Evidence = "No duplicate code detection tools found"
	return result, nil
}

func (c *DuplicateCodeDetectionChecker) checkJava(repo fs.FS, result *checker.Result) (*checker.Result, error) {
	for _, configFile := range []string{"pom.xml", "build.gradle"} {
		content, err := fs.ReadFile(repo, configFile)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("cpd")) || bytes.Contains(content, []byte("CPD")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found PMD CPD in %s — duplicate code detection configured", configFile)
			return result, nil
		}
	}
	result.Passed = false
	result.Evidence = "No duplicate code detection tools found"
	return result, nil
}
