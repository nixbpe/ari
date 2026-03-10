package style

import (
	"bytes"
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type LargeFileDetectionChecker struct{}

func (c *LargeFileDetectionChecker) ID() checker.CheckerID  { return "large_file_detection" }
func (c *LargeFileDetectionChecker) Pillar() checker.Pillar { return checker.PillarConstraints }
func (c *LargeFileDetectionChecker) Level() checker.Level   { return checker.LevelAutonomous }
func (c *LargeFileDetectionChecker) Name() string           { return "Large File Detection" }
func (c *LargeFileDetectionChecker) Description() string {
	return "Checks for large file detection and prevention tools"
}
func (c *LargeFileDetectionChecker) Suggestion() string {
	return "Add Git LFS for large files: git lfs install && git lfs track '*.bin'"
}

func (c *LargeFileDetectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if content, err := fs.ReadFile(repo, ".gitattributes"); err == nil {
		if bytes.Contains(content, []byte("filter=lfs")) {
			result.Passed = true
			result.Evidence = "Found .gitattributes with Git LFS configuration"
			return result, nil
		}
	}

	if content, err := fs.ReadFile(repo, ".pre-commit-config.yaml"); err == nil {
		if bytes.Contains(content, []byte("check-added-large-files")) ||
			bytes.Contains(content, []byte("file-size")) {
			result.Passed = true
			result.Evidence = "Found .pre-commit-config.yaml with file size hook"
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No large file detection configured"
	return result, nil
}
