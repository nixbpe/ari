package checker

import (
	"io/fs"
	"strings"
)

// CIWorkflowContains scans all .github/workflows/*.yml and *.yaml files for keywords.
// Returns (true, "workflows/ci.yml contains gitleaks") on first match.
// Returns (false, "") if no match or no workflows dir.
func CIWorkflowContains(repo fs.FS, keywords []string) (bool, string) {
	workflowDir := ".github/workflows"
	entries, err := fs.ReadDir(repo, workflowDir)
	if err != nil {
		return false, ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") {
			continue
		}

		path := workflowDir + "/" + name
		data, readErr := fs.ReadFile(repo, path)
		if readErr != nil {
			continue
		}

		content := string(data)
		for _, keyword := range keywords {
			if strings.Contains(content, keyword) {
				return true, path + " contains " + keyword
			}
		}
	}

	return false, ""
}

// DepFileContains scans dependency files for package/module names.
// Checks go.mod for Go, package.json for TypeScript, pom.xml/build.gradle for Java.
// Falls back to checking all dep files if lang is LanguageUnknown.
// Returns (true, "go.mod contains go.uber.org/zap") on first match.
func DepFileContains(repo fs.FS, lang Language, packages []string) (bool, string) {
	var depFiles []string

	switch lang {
	case LanguageGo:
		depFiles = []string{"go.mod"}
	case LanguageTypeScript:
		depFiles = []string{"package.json"}
	case LanguageJava:
		depFiles = []string{"pom.xml", "build.gradle"}
	default:
		// LanguageUnknown: try all
		depFiles = []string{"go.mod", "package.json", "pom.xml", "build.gradle"}
	}

	for _, depFile := range depFiles {
		data, err := fs.ReadFile(repo, depFile)
		if err != nil {
			continue
		}

		content := string(data)
		for _, pkg := range packages {
			if strings.Contains(content, pkg) {
				return true, depFile + " contains " + pkg
			}
		}
	}

	return false, ""
}

// FileExistsAny checks if any of the candidate paths exist in the repo.
// Returns (true, "path/to/file") for the first path found.
// Returns (false, "") if none exist.
func FileExistsAny(repo fs.FS, paths []string) (bool, string) {
	for _, path := range paths {
		if _, err := fs.Stat(repo, path); err == nil {
			return true, path
		}
	}
	return false, ""
}

// FileContentContains reads a file and checks if any keyword appears in its content.
// Returns (true, "keyword") for the first keyword found.
// Returns (false, "") if file missing or no keyword found.
func FileContentContains(repo fs.FS, path string, keywords []string) (bool, string) {
	data, err := fs.ReadFile(repo, path)
	if err != nil {
		return false, ""
	}

	content := string(data)
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true, keyword
		}
	}

	return false, ""
}
