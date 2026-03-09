package build

import (
	"context"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type ReleaseAutomationChecker struct{}

func (c *ReleaseAutomationChecker) ID() checker.CheckerID  { return "release_automation" }
func (c *ReleaseAutomationChecker) Pillar() checker.Pillar { return checker.PillarBuildSystem }
func (c *ReleaseAutomationChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *ReleaseAutomationChecker) Name() string           { return "Release Automation" }
func (c *ReleaseAutomationChecker) Description() string {
	return "Checks for release automation configuration (goreleaser, release-please, semantic-release)"
}
func (c *ReleaseAutomationChecker) Suggestion() string {
	return "Automate releases. For Go: use goreleaser. For JS: use release-please or semantic-release"
}

func (c *ReleaseAutomationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	releaseCandidates := []string{
		".goreleaser.yml",
		".goreleaser.yaml",
		"release-please-config.json",
		".releaserc",
		".releaserc.json",
		"cliff.toml",
	}

	for _, path := range releaseCandidates {
		if _, err := fs.Stat(repo, path); err == nil {
			result.Passed = true
			result.Evidence = "Found " + path + " — release automation configured"
			return result, nil
		}
	}

	if _, err := fs.Stat(repo, "Makefile"); err == nil {
		data, readErr := fs.ReadFile(repo, "Makefile")
		if readErr == nil && strings.Contains(string(data), "release") {
			result.Passed = true
			result.Evidence = "Found release target in Makefile — release automation configured"
			return result, nil
		}
	}

	if found, name := checkCIWorkflowForRelease(repo); found {
		result.Passed = true
		result.Evidence = "Found release job in " + name + " — release automation configured"
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No release automation found"
	return result, nil
}

func checkCIWorkflowForRelease(repo fs.FS) (bool, string) {
	workflowDirs := []string{".github/workflows", ".circleci"}
	for _, dir := range workflowDirs {
		entries, err := fs.ReadDir(repo, dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			path := dir + "/" + entry.Name()
			data, readErr := fs.ReadFile(repo, path)
			if readErr != nil {
				continue
			}
			content := string(data)
			if strings.Contains(content, "release") {
				return true, path
			}
		}
	}
	return false, ""
}
