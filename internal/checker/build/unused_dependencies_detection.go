package build

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type UnusedDependenciesDetectionChecker struct{}

func (c *UnusedDependenciesDetectionChecker) ID() checker.CheckerID {
	return "unused_dependencies_detection"
}
func (c *UnusedDependenciesDetectionChecker) Pillar() checker.Pillar {
	return checker.PillarBuildSystem
}
func (c *UnusedDependenciesDetectionChecker) Level() checker.Level { return checker.LevelAutonomous }
func (c *UnusedDependenciesDetectionChecker) Name() string {
	return "Unused Dependencies Detection"
}
func (c *UnusedDependenciesDetectionChecker) Description() string {
	return "Checks that unused dependency detection is configured"
}
func (c *UnusedDependenciesDetectionChecker) Suggestion() string {
	return "Add unused dependency detection. For Go: add go mod tidy to CI. For TS: npm install -D depcheck"
}

func (c *UnusedDependenciesDetectionChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	workflows, _ := fs.Glob(repo, ".github/workflows/*.yml")
	yamlFiles, _ := fs.Glob(repo, ".github/workflows/*.yaml")
	workflows = append(workflows, yamlFiles...)

	for _, wf := range workflows {
		content, err := fs.ReadFile(repo, wf)
		if err != nil {
			continue
		}
		if bytes.Contains(content, []byte("go mod tidy")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found go mod tidy in CI workflow %s — unused dependency detection configured", wf)
			return result, nil
		}
		if bytes.Contains(content, []byte("depcheck")) || bytes.Contains(content, []byte("knip")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found unused dependency tool in CI workflow %s", wf)
			return result, nil
		}
		if bytes.Contains(content, []byte("dependency:analyze")) {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found dependency:analyze in CI workflow %s — Maven unused dependency detection configured", wf)
			return result, nil
		}
	}

	if makefile, err := fs.ReadFile(repo, "Makefile"); err == nil {
		if bytes.Contains(makefile, []byte("mod-tidy")) || bytes.Contains(makefile, []byte("go mod tidy")) {
			result.Passed = true
			result.Evidence = "Found mod-tidy target in Makefile — unused dependency detection configured"
			return result, nil
		}
	}

	if pkgJSON, err := fs.ReadFile(repo, "package.json"); err == nil {
		if bytes.Contains(pkgJSON, []byte("depcheck")) || bytes.Contains(pkgJSON, []byte("knip")) {
			result.Passed = true
			result.Evidence = "Found depcheck/knip in package.json scripts — unused dependency detection configured"
			return result, nil
		}
	}

	if pomXML, err := fs.ReadFile(repo, "pom.xml"); err == nil {
		if bytes.Contains(pomXML, []byte("dependency:analyze")) {
			result.Passed = true
			result.Evidence = "Found dependency:analyze in pom.xml — Maven unused dependency detection configured"
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No unused dependency detection found"
	return result, nil
}
