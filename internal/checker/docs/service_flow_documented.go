package docs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/bbik/ari/internal/checker"
)

// ServiceFlowDocumentedChecker checks whether the repository contains architecture
// diagrams or documentation directories that describe service flows.
type ServiceFlowDocumentedChecker struct{}

func (c *ServiceFlowDocumentedChecker) ID() checker.CheckerID  { return "service_flow_documented" }
func (c *ServiceFlowDocumentedChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *ServiceFlowDocumentedChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *ServiceFlowDocumentedChecker) Name() string           { return "Service Flow Documentation" }
func (c *ServiceFlowDocumentedChecker) Description() string {
	return "Checks for architecture diagrams (.puml, .mmd, .drawio, .excalidraw) or docs/architecture directories"
}
func (c *ServiceFlowDocumentedChecker) Suggestion() string {
	return "Document service flows. Create architecture diagrams using PlantUML or Mermaid in docs/"
}

var diagramExtensions = []string{".puml", ".plantuml", ".mmd", ".drawio", ".excalidraw"}

var archDirs = []string{
	"docs/architecture",
	"docs/design",
	"docs/rfcs",
	"docs/adr",
	"adr",
}

var errFound = errors.New("found")

func (c *ServiceFlowDocumentedChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	// Check for known architecture directories first (cheap stat calls).
	for _, dir := range archDirs {
		if _, err := fs.Stat(repo, dir); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s/ — service flows documented", dir)
			return result, nil
		}
	}

	// Walk the entire repo looking for diagram files.
	var diagramCount int
	var firstDiagram string

	walkErr := fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		lp := strings.ToLower(path)
		for _, ext := range diagramExtensions {
			if strings.HasSuffix(lp, ext) {
				diagramCount++
				if firstDiagram == "" {
					firstDiagram = path
				}
				return nil
			}
		}
		return nil
	})
	if walkErr != nil && !errors.Is(walkErr, errFound) {
		return nil, walkErr
	}

	if diagramCount > 0 {
		result.Passed = true
		result.Evidence = fmt.Sprintf("Found %d diagram file(s) (e.g. %s) — service flows documented", diagramCount, firstDiagram)
		return result, nil
	}

	result.Passed = false
	result.Evidence = "No architecture documentation found"
	return result, nil
}
