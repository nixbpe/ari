package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

type BuildCmdDocChecker struct{}

func (c *BuildCmdDocChecker) ID() checker.CheckerID  { return "build_cmd_doc" }
func (c *BuildCmdDocChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *BuildCmdDocChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *BuildCmdDocChecker) Name() string           { return "Build Command Documentation" }
func (c *BuildCmdDocChecker) Description() string {
	return "Checks that build commands are documented in project files"
}
func (c *BuildCmdDocChecker) Suggestion() string {
	return "Document build commands in README.md. Add a Getting Started section with build instructions"
}

func (c *BuildCmdDocChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	docFiles := []string{"README.md", "CONTRIBUTING.md", "Makefile", "CLAUDE.md", "AGENTS.md"}
	buildCmds := buildCommandsForLang(lang)

	for _, docFile := range docFiles {
		content, err := fs.ReadFile(repo, docFile)
		if err != nil {
			continue
		}
		for _, cmd := range buildCmds {
			if bytes.Contains(content, []byte(cmd)) {
				result.Passed = true
				result.Evidence = fmt.Sprintf("Build command documented in %s: %s", docFile, cmd)
				return result, nil
			}
		}
	}

	result.Passed = false
	result.Evidence = "No build command documentation found"
	return result, nil
}

func buildCommandsForLang(lang checker.Language) []string {
	switch lang {
	case checker.LanguageGo:
		return []string{"go build", "make build", "./dev build"}
	case checker.LanguageTypeScript:
		return []string{"npm run build", "yarn build", "pnpm build"}
	case checker.LanguageJava:
		return []string{"mvn", "gradle", "./gradlew"}
	default:
		return []string{"go build", "make build", "npm run build", "yarn build", "./gradlew", "mvn"}
	}
}
