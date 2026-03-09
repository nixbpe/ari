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

## [T3] Checker Registry + Runner
- Registry uses map[CheckerID]Checker with mutex for thread safety (if needed)
- Runner uses recover() for panic handling
- ProgressFunc callback pattern used for TUI progress updates

## [T5] LLM Interface
- FallbackEvaluator: Primary nil -> always use Fallback; Primary error -> use Fallback
- Mode field: "llm" when primary used, "rule-based" when fallback used
- MockProvider: CompleteFunc nil -> returns canned JSON response
- ConfigFromEnv reads: ARI_LLM_PROVIDER, ARI_API_KEY, ARI_LLM_MODEL, ARI_LLM_BASE_URL

## [T4] Scorer
- Gated progression: level N requires ≥80% at level N AND all previous levels pass
- Skipped results excluded from both numerator and denominator
- Level 0 returned when no level achieved (below L1 threshold)
- Use Result.Level field for grouping, not hardcoded criterion IDs

## [T6] Reporter
- Report struct is the central data model for all output formats
- JSONReporter writes to io.Writer (not file directly)
- Suggestions generated only for failing non-skipped criteria
- BuildReport assembles from RepoInfo + Score + []Result
- Level and PassRate duplicated as top-level Report fields (from Score) so JSON has "level" and "passRate" at root
- scorer.Score map fields (PillarScores/LevelScores) use int keys which Go JSON marshals as string keys — round-trips cleanly in Go 1.7+
- priorityForLevel: Functional/Documented→high, Standardized/Optimized→medium, Autonomous→low

## [T7] Style Checkers
- Checker files go in internal/checker/style/ package
- Each checker is a separate struct implementing Checker interface
- Go and Java always pass type_check and formatter (statically typed / built-in)
- Use fs.ReadFile(repo, filename) to check file existence — returns error if not found

## [T9] Style Checkers (cyclomatic, dead_code, duplicate_code)
- golangci-lint YAML parsing: use gopkg.in/yaml.v3 or simple string search in file content
- For simple config checks, string search (bytes.Contains) is more robust than full YAML parse
- package.json devDependencies: unmarshal JSON, check map key existence
- duplicate_code_detection is LevelOptimized (Level 4), the other two are LevelStandardized (Level 3)
- Language-agnostic checks (like .jscpd.json) should run before language-specific switch

## [T11] Build Checkers (build_cmd_doc, single_command_setup, deps_pinned)
- build_cmd_doc: search multiple doc files for language-specific build commands
- Use bytes.Contains(content, []byte("go build")) for simple content search
- deps_pinned: check lock file existence per language

## [T13] Build Checkers (vcs_cli_tools, agentic_development, automated_pr_review)
- agentic_development: check multiple AI agent doc files (AGENTS.md, CLAUDE.md, .cursor/rules)
- CODEOWNERS can be at root or .github/CODEOWNERS

## [T18 complete] Documentation Checkers
- documentation_freshness: GitRunner func field pattern for testability
- Parse git date format: "2006-01-02 15:04:05 -0700"
- skills: use fs.ReadDir(repo, ".claude/skills") — returns error if dir missing

## [T19] Documentation Checkers (service_flow, api_schema_docs)
- service_flow: walk all files for diagram extensions (.puml, .mmd, .drawio); also stat well-known arch dirs first
- api_schema_docs: skip heuristic — check if any .go file contains "net/http"; Java never skipped
- Use fs.WalkDir with early exit via sentinel errors.New("found") pattern
- t19_test.go created separately (docs_test.go already existed from T18); both in same package `docs`
- ctx variable already declared in docs_test.go — do NOT redeclare in t19_test.go

## [T21] Text Reporter
- Pure text, no ANSI codes (no `\x1b[`, no `\033[`)
- Group criteria by pillar using CriteriaResults[i].Pillar field
- ✓ (U+2713), ✗ (U+2717), ↷ (U+21B7) are fine UTF-8 chars
- Output format: header with repo/language/level/pass-rate, pillars with criteria, suggestions, footer with version/timestamp
- TextReporter implements Reporter interface: Report(ctx context.Context, report *Report, w io.Writer) error

## [T22] TUI (Bubbletea v2)
- Import: charm.land/bubbletea/v2
- View() on root Model returns tea.View (use tea.NewView(content))
- Sub-models (ProgressModel, ReportModel) return string from View()
- Init() returns tea.Cmd in bubbletea v2 interface; root update still follows Elm-style state transitions
- Message types defined in internal/tui/messages.go
- views package: internal/tui/views/

## [T20] HTML Reporter
- Template embedded with `//go:embed templates/report.html` (file must be in a subdirectory of the package)
- No external deps — all CSS inline in `<style>` tag, no http:// or https:// anywhere
- Use `html/template` (not `text/template`) for XSS-safe output
- CSS width in style attributes must use `template.CSS` type (not plain string) to avoid `ZgotmplZ` sanitization
- `barStyle(f float64) template.CSS` returns `template.CSS(fmt.Sprintf("width:%d%%", pct))`
- Template functions for maps → sorted slices: `sortedLevelScores`, `sortedPillarScores`
- "Style & Validation" is escaped to "Style &amp; Validation" by html/template in text context — test for the HTML-encoded form
- Test HTML validity with `strings.Contains(output, "<html")` not `golang.org/x/net/html`
- Go 1.21+ has builtin `min` — do NOT define custom `min` in test files (would shadow, may cause confusion)

## [T24] LLM Providers
- Use net/http directly (no SDK deps) — httptest.NewServer for all tests
- OpenAI: POST /v1/chat/completions, Authorization: Bearer header
- Anthropic: POST /v1/messages, x-api-key header, anthropic-version: 2023-06-01
- Ollama: POST /api/chat, no auth, stream: false
- Rate limiting: exponential backoff (1s, 2s, 4s), max 3 retries
- Default models: OpenAI=gpt-4o-mini, Anthropic=claude-sonnet-4-20250514, Ollama=llama3
