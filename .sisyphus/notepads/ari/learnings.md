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
