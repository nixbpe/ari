package security

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type GitignoreComprehensiveChecker struct{}

func (c *GitignoreComprehensiveChecker) ID() checker.CheckerID  { return "gitignore_comprehensive" }
func (c *GitignoreComprehensiveChecker) Pillar() checker.Pillar { return checker.PillarSecurity }
func (c *GitignoreComprehensiveChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *GitignoreComprehensiveChecker) Name() string           { return "Comprehensive .gitignore" }
func (c *GitignoreComprehensiveChecker) Description() string {
	return "Checks that .gitignore exists and covers at least 3 security-sensitive patterns"
}
func (c *GitignoreComprehensiveChecker) Suggestion() string {
	return "Add security patterns to .gitignore: .env, *.pem, *.key, .idea/, .vscode/, node_modules, *.log, *.secret"
}

var securityPatterns = []string{".env", "*.pem", "*.key", ".idea/", ".vscode/", "node_modules", "*.log", "*.secret"}

func (c *GitignoreComprehensiveChecker) Check(_ context.Context, repo fs.FS, _ checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Suggestion: c.Suggestion(),
		Mode:       "rule-based",
	}

	data, err := fs.ReadFile(repo, ".gitignore")
	if err != nil {
		result.Passed = false
		result.Evidence = ".gitignore not found"
		return result, nil
	}

	content := string(data)
	matches := 0
	for _, pattern := range securityPatterns {
		if strings.Contains(content, pattern) {
			matches++
		}
	}

	if matches >= 3 {
		result.Passed = true
		result.Evidence = fmt.Sprintf(".gitignore covers %d security patterns", matches)
	} else {
		result.Passed = false
		result.Evidence = fmt.Sprintf(".gitignore only covers %d security patterns (need ≥3)", matches)
	}

	return result, nil
}
