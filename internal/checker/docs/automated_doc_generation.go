package docs

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type AutomatedDocGenerationChecker struct{}

func (c *AutomatedDocGenerationChecker) ID() checker.CheckerID  { return "automated_doc_generation" }
func (c *AutomatedDocGenerationChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *AutomatedDocGenerationChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *AutomatedDocGenerationChecker) Name() string           { return "Automated Documentation Generation" }
func (c *AutomatedDocGenerationChecker) Description() string {
	return "Checks for automated documentation generation tooling (godoc, typedoc, javadoc, etc.)"
}
func (c *AutomatedDocGenerationChecker) Suggestion() string {
	return "Add automated doc generation. For Go: use godoc. For TS: npm install -D typedoc"
}

func (c *AutomatedDocGenerationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
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
		if found, evidence := checkGoDocGen(repo); found {
			result.Passed = true
			result.Evidence = evidence
			return result, nil
		}
	case checker.LanguageTypeScript:
		if found, evidence := checkTSDocGen(repo); found {
			result.Passed = true
			result.Evidence = evidence
			return result, nil
		}
	case checker.LanguageJava:
		if found, evidence := checkJavaDocGen(repo); found {
			result.Passed = true
			result.Evidence = evidence
			return result, nil
		}
	}

	if _, err := fs.Stat(repo, "docs"); err == nil {
		result.Passed = true
		result.Evidence = "Found docs/ directory — documentation present"
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No automated doc generation found"
	return result, nil
}

func checkGoDocGen(repo fs.FS) (bool, string) {
	if data, err := fs.ReadFile(repo, "go.mod"); err == nil {
		content := string(data)
		if strings.Contains(content, "pkgsite") {
			return true, "Found pkgsite in go.mod — automated doc generation configured"
		}
		if strings.Contains(content, "godoc") {
			return true, "Found godoc in go.mod — automated doc generation configured"
		}
	}

	// Walk for //go:generate in .go source files
	var generateCount int
	fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		data, readErr := fs.ReadFile(repo, path)
		if readErr != nil {
			return nil
		}
		if strings.Contains(string(data), "//go:generate") {
			generateCount++
		}
		return nil
	})

	if generateCount > 0 {
		return true, fmt.Sprintf("Found //go:generate directives in %d Go file(s) — automated generation configured", generateCount)
	}

	return false, ""
}

// checkTSDocGen checks package.json devDependencies for typedoc, storybook, or docusaurus.
func checkTSDocGen(repo fs.FS) (bool, string) {
	data, err := fs.ReadFile(repo, "package.json")
	if err != nil {
		return false, ""
	}
	content := string(data)

	tools := []string{"typedoc", "storybook", "docusaurus"}
	for _, tool := range tools {
		if strings.Contains(content, tool) {
			return true, fmt.Sprintf("Found %s in devDependencies — automated doc generation configured", tool)
		}
	}

	return false, ""
}

// checkJavaDocGen checks pom.xml for maven-javadoc-plugin or build.gradle for dokka.
func checkJavaDocGen(repo fs.FS) (bool, string) {
	if data, err := fs.ReadFile(repo, "pom.xml"); err == nil {
		if strings.Contains(string(data), "maven-javadoc-plugin") {
			return true, "Found maven-javadoc-plugin in pom.xml — automated doc generation configured"
		}
	}

	gradleFiles := []string{"build.gradle", "build.gradle.kts"}
	for _, gf := range gradleFiles {
		if data, err := fs.ReadFile(repo, gf); err == nil {
			if strings.Contains(string(data), "dokka") {
				return true, fmt.Sprintf("Found dokka in %s — automated doc generation configured", gf)
			}
		}
	}

	return false, ""
}
