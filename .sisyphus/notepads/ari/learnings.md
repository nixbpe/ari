# Learnings — ari Project

## Environment
- Go 1.26.0 installed (satisfies Go 1.24+ requirement)
- Working directory: /Users/khakhana.t/Code/BBIK/lab/agent-readiness (this IS the module root)
- Git initialized, main branch

## Architecture Decisions
- Module name: github.com/bbik/ari
- Bubbletea v2 import: charm.land/bubbletea/v2 (NOT github.com/charmbracelet/bubbletea)
- View() returns tea.View (not string) in Bubbletea v2
- ALL checker code must use fs.FS, never os.* directly
- Checker registration: explicit NewDefaultRegistry() — no init() magic

## Key Patterns
- FallbackEvaluator: LLM → rule-based → never fails
- fs.FS for all file access → MapFS in tests
- Each criterion: ID, Pillar, Level, pass/fail, evidence, skip logic

## Conventions
- Package names: internal/checker, internal/scanner, internal/scorer, internal/llm, internal/reporter, internal/tui
- Test files: *_test.go (Go convention)
- Evidence files: .sisyphus/evidence/task-{N}-{scenario-slug}.txt

## T2 Scanner Notes
- `DefaultScanner` should walk with `fs.WalkDir` and skip common heavy directories by directory name (`.git`, `node_modules`, `vendor`, etc.)
- For `MapFS` tests, `.git/HEAD` implies `.git` is visible as a directory to `fs.Stat`, so git detection can stay filesystem-interface based
- Language detection is most stable with signature-file priority (`go.mod`, `pom.xml`/`build.gradle*`, then `package.json`) before extension counting
- File limits are easiest to enforce with a simple counter checked before appending to `RepoInfo.Files`

## [T5] LLM Interface
- FallbackEvaluator: Primary nil -> always use Fallback; Primary error -> use Fallback
- Mode field: "llm" when primary used, "rule-based" when fallback used
- MockProvider: CompleteFunc nil -> returns canned JSON response
- ConfigFromEnv reads: ARI_LLM_PROVIDER, ARI_API_KEY, ARI_LLM_MODEL, ARI_LLM_BASE_URL
