package observability

import (
	"context"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

type ProfilingInstrumentationChecker struct{}

func (c *ProfilingInstrumentationChecker) ID() checker.CheckerID  { return "profiling_instrumentation" }
func (c *ProfilingInstrumentationChecker) Pillar() checker.Pillar { return checker.PillarObservability }
func (c *ProfilingInstrumentationChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *ProfilingInstrumentationChecker) Name() string           { return "Profiling Instrumentation" }
func (c *ProfilingInstrumentationChecker) Description() string {
	return "Checks for profiling instrumentation in the codebase"
}
func (c *ProfilingInstrumentationChecker) Suggestion() string {
	return "Add profiling support (e.g., import net/http/pprof in Go; clinic or pyroscope for Node.js)"
}

func (c *ProfilingInstrumentationChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	if lang == checker.LanguageGo || lang == checker.LanguageUnknown {
		goModFound, goModEvidence := checker.DepFileContains(repo, checker.LanguageGo, []string{"fgprof"})
		if goModFound {
			result.Passed = true
			result.Evidence = "Found profiling: " + goModEvidence
			return result, nil
		}

		pprofFound := false
		pprofEvidence := ""
		_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			data, readErr := fs.ReadFile(repo, path)
			if readErr != nil {
				return nil
			}
			if strings.Contains(string(data), "net/http/pprof") {
				pprofFound = true
				pprofEvidence = path + " imports net/http/pprof"
				return fs.SkipAll
			}
			return nil
		})

		if pprofFound {
			result.Passed = true
			result.Evidence = "Found profiling: " + pprofEvidence
			return result, nil
		}
	}

	if lang == checker.LanguageTypeScript || lang == checker.LanguageUnknown {
		tsFound, tsEvidence := checker.DepFileContains(repo, checker.LanguageTypeScript, []string{"clinic", "0x", "pyroscope"})
		if tsFound {
			result.Passed = true
			result.Evidence = "Found profiling: " + tsEvidence
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No profiling instrumentation found (checked net/http/pprof, fgprof, clinic, 0x, pyroscope)"
	return result, nil
}
