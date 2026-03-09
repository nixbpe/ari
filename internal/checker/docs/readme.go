package docs

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type ReadmeChecker struct{}

func (c *ReadmeChecker) ID() checker.CheckerID  { return "readme" }
func (c *ReadmeChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *ReadmeChecker) Level() checker.Level   { return checker.LevelFunctional }
func (c *ReadmeChecker) Name() string           { return "README Exists" }
func (c *ReadmeChecker) Description() string {
	return "Checks that a README file exists and has meaningful content (50+ chars with headings)"
}
func (c *ReadmeChecker) Suggestion() string {
	return "Create a README.md with: project description, installation, usage, contributing guide"
}

var readmeCandidates = []string{"README.md", "README", "README.rst", "README.txt"}

func (c *ReadmeChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	for _, name := range readmeCandidates {
		data, err := fs.ReadFile(repo, name)
		if err != nil {
			continue
		}

		if len(data) < 50 {
			result.Passed = false
			result.Evidence = fmt.Sprintf("%s exists but has insufficient content", name)
			return result, nil
		}

		hasHeading := bytes.Contains(data, []byte("#")) || bytes.Contains(data, []byte("====="))
		if !hasHeading {
			result.Passed = false
			result.Evidence = fmt.Sprintf("%s exists but has insufficient content", name)
			return result, nil
		}

		result.Passed = true
		result.Evidence = fmt.Sprintf("Found %s (%d bytes) with headings", name, len(data))
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No README found"
	return result, nil
}
