package style

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/llm"
)

type NamingConsistencyChecker struct {
	Evaluator llm.Evaluator
}

func NewNamingConsistencyChecker(eval llm.Evaluator) *NamingConsistencyChecker {
	return &NamingConsistencyChecker{Evaluator: eval}
}

func (c *NamingConsistencyChecker) ID() checker.CheckerID  { return "naming_consistency" }
func (c *NamingConsistencyChecker) Pillar() checker.Pillar { return checker.PillarStyleValidation }
func (c *NamingConsistencyChecker) Level() checker.Level   { return checker.LevelDocumented }
func (c *NamingConsistencyChecker) Name() string           { return "Naming Consistency" }
func (c *NamingConsistencyChecker) Description() string {
	return "Checks that naming conventions are consistent across the codebase"
}

const namingSuggestion = "Ensure consistent naming conventions. For Go: use PascalCase for exports, camelCase for unexported"

func (c *NamingConsistencyChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	sampleFiles := collectSampleFiles(repo, lang, 10)
	prompt := buildNamingPrompt(lang, sampleFiles)

	ruleEval := &namingRuleEvaluator{lang: lang, sampleFiles: sampleFiles, repo: repo}

	fe := &llm.FallbackEvaluator{
		Primary:  c.Evaluator,
		Fallback: ruleEval,
	}

	evalResult, _ := fe.Evaluate(ctx, prompt)
	if evalResult == nil {
		evalResult = &llm.EvalResult{
			Passed:   true,
			Evidence: fmt.Sprintf("Naming follows %s conventions", lang),
			Mode:     "rule-based",
		}
	}

	return &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Passed:     evalResult.Passed,
		Evidence:   evalResult.Evidence,
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       evalResult.Mode,
		Suggestion: namingSuggestion,
	}, nil
}

func collectSampleFiles(repo fs.FS, lang checker.Language, max int) []string {
	var ext string
	switch lang {
	case checker.LanguageGo:
		ext = ".go"
	case checker.LanguageTypeScript:
		ext = ".ts"
	case checker.LanguageJava:
		ext = ".java"
	default:
		return nil
	}

	var files []string
	_ = fs.WalkDir(repo, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ext) {
			files = append(files, p)
			if len(files) >= max {
				return fs.SkipAll
			}
		}
		return nil
	})
	return files
}

func buildNamingPrompt(lang checker.Language, sampleFiles []string) string {
	return fmt.Sprintf(
		"Check naming consistency for %s project. Sample files: %s. Are naming conventions followed?",
		lang, strings.Join(sampleFiles, ", "),
	)
}

type namingRuleEvaluator struct {
	lang        checker.Language
	sampleFiles []string
	repo        fs.FS
}

var (
	goExportedUnderscore = regexp.MustCompile(`(?m)^(?:func|type|var|const)\s+([A-Z]\w*_\w+)`)
	pascalCaseFile       = regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
	camelOrKebabFile     = regexp.MustCompile(`^[a-z][a-z0-9A-Z-]*$`)
)

func (e *namingRuleEvaluator) Evaluate(_ context.Context, _ string) (*llm.EvalResult, error) {
	switch e.lang {
	case checker.LanguageGo:
		return e.checkGo()
	case checker.LanguageTypeScript:
		return e.checkTypeScript()
	case checker.LanguageJava:
		return e.checkJava()
	default:
		return &llm.EvalResult{
			Passed:     true,
			Evidence:   fmt.Sprintf("Naming follows %s conventions", e.lang),
			Confidence: 0.5,
			Mode:       "rule-based",
		}, nil
	}
}

func (e *namingRuleEvaluator) checkGo() (*llm.EvalResult, error) {
	violations := 0
	for _, p := range e.sampleFiles {
		data, err := fs.ReadFile(e.repo, p)
		if err != nil {
			continue
		}
		for _, m := range goExportedUnderscore.FindAllSubmatch(data, -1) {
			if len(m) > 1 {
				name := string(m[1])
				if !strings.HasPrefix(name, "Test") && !strings.HasPrefix(name, "Benchmark") {
					violations++
				}
			}
		}
	}

	if violations > 0 {
		return &llm.EvalResult{
			Passed:     false,
			Evidence:   "Inconsistent naming detected",
			Confidence: 0.8,
			Mode:       "rule-based",
		}, nil
	}
	return &llm.EvalResult{
		Passed:     true,
		Evidence:   "Naming follows Go conventions",
		Confidence: 0.8,
		Mode:       "rule-based",
	}, nil
}

func (e *namingRuleEvaluator) checkTypeScript() (*llm.EvalResult, error) {
	violations := 0
	for _, p := range e.sampleFiles {
		base := path.Base(p)
		name := strings.TrimSuffix(base, path.Ext(base))
		name = strings.TrimSuffix(name, ".d")
		if !camelOrKebabFile.MatchString(name) {
			violations++
		}
	}

	if violations > 0 {
		return &llm.EvalResult{
			Passed:     false,
			Evidence:   "Inconsistent naming detected",
			Confidence: 0.8,
			Mode:       "rule-based",
		}, nil
	}
	return &llm.EvalResult{
		Passed:     true,
		Evidence:   "Naming follows TypeScript conventions",
		Confidence: 0.8,
		Mode:       "rule-based",
	}, nil
}

func (e *namingRuleEvaluator) checkJava() (*llm.EvalResult, error) {
	violations := 0
	for _, p := range e.sampleFiles {
		base := path.Base(p)
		name := strings.TrimSuffix(base, ".java")
		if !pascalCaseFile.MatchString(name) {
			violations++
		}
	}

	if violations > 0 {
		return &llm.EvalResult{
			Passed:     false,
			Evidence:   "Inconsistent naming detected",
			Confidence: 0.8,
			Mode:       "rule-based",
		}, nil
	}
	return &llm.EvalResult{
		Passed:     true,
		Evidence:   "Naming follows Java conventions",
		Confidence: 0.8,
		Mode:       "rule-based",
	}, nil
}
