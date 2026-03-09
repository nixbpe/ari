package scanner

import (
	"path/filepath"

	"github.com/nixbpe/ari/internal/checker"
)

func DetectLanguage(files []FileInfo) checker.Language {
	for _, f := range files {
		switch filepath.Base(f.Path) {
		case "go.mod":
			return checker.LanguageGo
		case "pom.xml", "build.gradle", "build.gradle.kts":
			return checker.LanguageJava
		}
	}

	for _, f := range files {
		if filepath.Base(f.Path) == "package.json" {
			return checker.LanguageTypeScript
		}
	}

	counts := map[checker.Language]int{}
	for _, f := range files {
		switch f.Extension {
		case ".go":
			counts[checker.LanguageGo]++
		case ".ts", ".tsx":
			counts[checker.LanguageTypeScript]++
		case ".java", ".kt", ".kts":
			counts[checker.LanguageJava]++
		}
	}

	best := checker.LanguageUnknown
	bestCount := 0
	for lang, count := range counts {
		if count > bestCount {
			best = lang
			bestCount = count
		}
	}

	return best
}
