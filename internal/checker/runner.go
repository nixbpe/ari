package checker

import (
	"context"
	"fmt"
	"io/fs"
	"reflect"
)

type ProgressFunc func(done, total int, id CheckerID)

type Runner struct {
	Registry     *Registry
	Evaluator    interface{}
	ProgressFunc ProgressFunc

	OnStart func(id CheckerID, name string)
	OnDone  func(result *Result, done, total int)
}

func (r *Runner) Run(ctx context.Context, repo fs.FS, repoInfo interface{}) ([]*Result, error) {
	if r == nil || r.Registry == nil {
		return nil, fmt.Errorf("runner registry is required")
	}

	lang := languageFromRepoInfo(repoInfo)
	checkers := r.Registry.All()
	results := make([]*Result, 0, len(checkers))
	total := len(checkers)

	for i, ch := range checkers {
		if err := ctx.Err(); err != nil {
			return results, err
		}

		if r.OnStart != nil {
			r.OnStart(ch.ID(), ch.Name())
		}

		result := r.runOne(ctx, repo, lang, ch)
		results = append(results, result)

		if r.OnDone != nil {
			r.OnDone(result, i+1, total)
		}
		if r.ProgressFunc != nil {
			r.ProgressFunc(i+1, total, ch.ID())
		}
	}

	return results, nil
}

func (r *Runner) runOne(ctx context.Context, repo fs.FS, lang Language, ch Checker) (res *Result) {
	if applicable, reason := isApplicable(ch, lang); !applicable {
		return &Result{
			ID:         ch.ID(),
			Name:       ch.Name(),
			Pillar:     ch.Pillar(),
			Level:      ch.Level(),
			Mode:       "rule-based",
			Skipped:    true,
			SkipReason: reason,
			Passed:     true,
		}
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			res = &Result{
				ID:       ch.ID(),
				Name:     ch.Name(),
				Pillar:   ch.Pillar(),
				Level:    ch.Level(),
				Mode:     "rule-based",
				Passed:   false,
				Evidence: fmt.Sprintf("checker panicked: %v", recovered),
			}
		}
	}()

	out, err := ch.Check(ctx, repo, lang)
	if err != nil {
		return &Result{
			ID:       ch.ID(),
			Name:     ch.Name(),
			Pillar:   ch.Pillar(),
			Level:    ch.Level(),
			Mode:     "rule-based",
			Passed:   false,
			Evidence: err.Error(),
		}
	}

	if out == nil {
		return &Result{
			ID:       ch.ID(),
			Name:     ch.Name(),
			Pillar:   ch.Pillar(),
			Level:    ch.Level(),
			Mode:     "rule-based",
			Passed:   false,
			Evidence: "checker returned nil result",
		}
	}

	if out.ID == "" {
		out.ID = ch.ID()
	}
	if out.Name == "" {
		out.Name = ch.Name()
	}
	if out.Pillar < PillarStyleValidation || out.Pillar > PillarDocumentation {
		out.Pillar = ch.Pillar()
	}
	if out.Level == 0 {
		out.Level = ch.Level()
	}
	if out.Mode == "" {
		out.Mode = "rule-based"
	}

	if !out.Passed {
		if provider, ok := ch.(SuggestionProvider); ok && out.Suggestion == "" {
			out.Suggestion = provider.Suggestion()
		}
	}

	return out
}

func languageFromRepoInfo(repoInfo interface{}) Language {
	if repoInfo == nil {
		return LanguageUnknown
	}

	v := reflect.ValueOf(repoInfo)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return LanguageUnknown
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return LanguageUnknown
	}

	field := v.FieldByName("Language")
	if !field.IsValid() {
		return LanguageUnknown
	}

	langType := reflect.TypeOf(LanguageUnknown)
	if field.Type() != langType {
		return LanguageUnknown
	}

	lang, ok := field.Interface().(Language)
	if !ok {
		return LanguageUnknown
	}
	return lang
}

func isApplicable(ch Checker, lang Language) (bool, string) {
	type languageScopedChecker interface {
		SupportsLanguage(lang Language) bool
	}

	if scoped, ok := ch.(languageScopedChecker); ok {
		if !scoped.SupportsLanguage(lang) {
			return false, fmt.Sprintf("not applicable for %s", lang.String())
		}
	}

	return true, ""
}
