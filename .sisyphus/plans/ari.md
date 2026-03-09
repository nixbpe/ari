# ARI — Agent Readiness Index

## TL;DR

> **Quick Summary**: Build "ari" — a Go CLI tool with Bubbletea TUI that evaluates how ready a codebase is for AI coding agents. Inspired by Factory.ai's Agent Readiness framework, ari scans a local repository, checks 40 criteria across 4 pillars, assigns a maturity level (1–5), and generates an HTML report with actionable suggestions.
> 
> **Deliverables**:
> - `ari` CLI binary (Go, single binary distribution)
> - Interactive TUI with progress → report → detail drill-down views
> - 40 checkers across 4 pillars (Style & Validation, Build System, Testing, Documentation)
> - 5-level maturity scoring with gated progression (80% threshold)
> - HTML report output (standalone, self-contained)
> - JSON output for programmatic use
> - Multi-provider LLM integration (OpenAI, Anthropic, Ollama) with rule-based fallback
> - Support for Go, TypeScript/JavaScript, Java/Kotlin repositories
> 
> **Estimated Effort**: Large
> **Parallel Execution**: YES — 7 waves
> **Critical Path**: T1 → T2/T3 → T7-T19 → T25 → T27 → FINAL

---

## Context

### Original Request
User studied Factory.ai's Agent Readiness framework (https://factory.ai/news/agent-readiness and https://factory.ai/agent-readiness/cockroachdb_cockroach) and wants to build a similar system as a Go CLI tool to evaluate how ready a project is for AI coding agents.

### Interview Summary
**Key Discussions**:
- **Form Factor**: CLI tool with TUI using Bubbletea v2 (charmbracelet/bubbletea, Elm Architecture)
- **Tech Stack**: Go (chosen for single binary distribution, performance)
- **Evaluation Method**: Hybrid — Rule-based for clear checks (file existence, config parsing), LLM for criteria needing interpretation (doc quality, naming consistency)
- **LLM Provider**: Multi-provider via interface (user chooses OpenAI, Anthropic, or Ollama)
- **Target Languages**: Go, Java/Kotlin, TypeScript/JavaScript repos
- **MVP Scope**: Core 4 Pillars (Style & Validation, Build System, Testing, Documentation)
- **Scoring**: 5 maturity levels with 80% gated progression (following Factory.ai)
- **Output**: HTML report (standalone file, viewable in browser from TUI) + JSON
- **Scan Source**: Local directory path only (no GitHub API)
- **Remediation**: Report + Suggestions only (no auto-fix)
- **Criteria Source**: Follow Factory.ai criteria as baseline
- **Project Name**: "ari" (Agent Readiness Index)

**Research Findings**:
- Factory.ai has 9 pillars (not 8 as blog states), 60+ criteria, binary pass/fail
- Bubbletea v2 released Feb 2026 — `View()` returns `tea.View` not `string`, import path is `charm.land/bubbletea/v2`
- Factory.ai solved LLM evaluation variance (7% → 0.6%) by grounding on previous reports
- CockroachDB at Level 4 (74%), FastAPI at Level 3 (53%), Express at Level 2 (28%)
- Each criterion includes: pass/fail, evidence string, level assignment, skip logic
- Scanner MUST use `fs.FS` interface for testability (`testing/fstest.MapFS`)

### Metis Review
**Identified Gaps** (addressed):
- **Criteria-to-Level mapping was missing**: Applied Factory.ai-like mapping (see Scoring section)
- **Monorepo handling unclear**: Explicitly excluded from MVP — single-app repos only
- **LLM fallback behavior undefined**: Applied FallbackEvaluator pattern (LLM → rule-based)
- **Configuration strategy missing**: Applied env vars + CLI flags (no config file in MVP)
- **Git dependency unstated**: Applied "warn + skip git-dependent criteria" for non-git dirs
- **40+ criteria is massive scope**: Grouped criteria into batches of 2–4 per task for manageability
- **Bubbletea v2 is brand new**: Keep TUI as thin presentation layer, all logic is pure Go
- **HTML report over-engineering risk**: Locked to `html/template` + embedded CSS, no JavaScript
- **Multi-language detection complexity**: Detect primary language by file count ratio

---

## Work Objectives

### Core Objective
Build a Go CLI tool called "ari" that scans a local repository and evaluates its readiness for AI coding agents across 40 criteria in 4 pillars, presenting results via an interactive TUI and HTML report.

### Concrete Deliverables
- `cmd/ari/main.go` — CLI entry point with flag parsing
- `internal/checker/` — 40 checker implementations across 4 pillars
- `internal/scanner/` — Repository scanner with `fs.FS` and language detection
- `internal/scorer/` — 5-level maturity scoring engine
- `internal/llm/` — Multi-provider LLM interface with fallback
- `internal/reporter/` — HTML + JSON + Text reporters
- `internal/tui/` — Bubbletea v2 interactive views
- `testdata/` — Integration test fixtures
- Standalone HTML report template (embedded via `//go:embed`)

### Definition of Done
- [ ] `go test ./...` — all tests pass, zero failures
- [ ] `go build ./cmd/ari` — produces working binary
- [ ] `./ari --path ./testdata/sample-go-repo` — launches TUI, shows report
- [ ] `./ari --path ./testdata/sample-go-repo --output json` — outputs valid JSON with level 1–5
- [ ] `./ari --path ./testdata/sample-go-repo --output html --out /tmp/report.html` — generates self-contained HTML
- [ ] LLM evaluation works with `ARI_LLM_PROVIDER=openai` and falls back gracefully without it
- [ ] All 40 criteria have pass/fail tests per supported language

### Must Have
- All 40 criteria across 4 pillars implemented with evidence
- 5-level gated progression scoring (80% threshold per level)
- TUI with progress → report → detail views
- HTML report (single self-contained file, no external deps)
- JSON output for programmatic use
- Multi-provider LLM via interface with rule-based fallback
- Go, TypeScript/JavaScript, Java/Kotlin language detection
- `fs.FS` interface for all file system access (testability)
- Criteria skip logic for inapplicable checks
- Suggestion text for each failing criterion

### Must NOT Have (Guardrails)
- **No monorepo support** — single-app repos only, no sub-application discovery
- **No automated remediation** — report + suggestions only, never create/modify repo files
- **No GitHub API calls** — all evaluation is local filesystem + git CLI only
- **No criteria beyond the defined 40** — no "while we're at it" additions
- **No interactive charts or JavaScript in HTML report** — static HTML + inline CSS only
- **No configuration file system** — env vars + CLI flags only, no `.ari.yaml`
- **No languages beyond Go, Java/Kotlin, TypeScript/JavaScript**
- **No report comparison, history tracking, or trending**
- **No TUI animations beyond functional feedback** (spinner, progress bar)
- **No `as any`, `@ts-ignore` equivalent (`//nolint` without reason)**
- **No global mutable state** — all checkers are pure functions on `fs.FS`
- **No direct `os` file system calls in checker code** — always `fs.FS`

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: NO (greenfield — will be set up in T1)
- **Automated tests**: TDD (RED → GREEN → REFACTOR)
- **Framework**: Go built-in `testing` package + `testing/fstest.MapFS`
- **TUI testing**: `teatest` from `charmbracelet/x/exp/teatest`
- **HTML testing**: Golden file comparison with `-update` flag

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Checkers**: Use `go test` with `MapFS` fixtures — verify pass/fail/skip per language
- **Scorer**: Use `go test` with hardcoded `Result` slices — verify level calculation
- **Reporters**: Use `go test` + golden file comparison — verify output format
- **TUI**: Use `teatest` — verify state transitions and view rendering
- **CLI**: Use `go build` + run binary against testdata — verify end-to-end

### Criteria-to-Level Mapping (Factory.ai-based)

**Level 1 — Functional** (7 criteria):
`readme`, `lint_config`, `formatter`, `type_check`, `unit_tests_exist`, `build_cmd_doc`, `deps_pinned`

**Level 2 — Documented** (8 criteria):
`agents_md`, `documentation_freshness`, `strict_typing`, `pre_commit_hooks`, `naming_consistency`, `single_command_setup`, `unit_tests_runnable`, `test_naming_conventions`

**Level 3 — Standardized** (9 criteria):
`integration_tests_exist`, `test_coverage_thresholds`, `service_flow_documented`, `api_schema_docs`, `fast_ci_feedback`, `release_automation`, `agentic_development`, `cyclomatic_complexity`, `dead_code_detection`

**Level 4 — Optimized** (12 criteria):
`automated_doc_generation`, `skills`, `automated_pr_review`, `build_performance_tracking`, `feature_flag_infrastructure`, `release_notes_automation`, `flaky_test_detection`, `test_performance_tracking`, `test_isolation`, `duplicate_code_detection`, `code_modularization`, `deployment_frequency`

**Level 5 — Autonomous** (4 criteria):
`tech_debt_tracking`, `large_file_detection`, `unused_dependencies_detection`, `vcs_cli_tools`

**Progression Rule**: To unlock Level N, must pass ≥80% of criteria at Level N AND all previous levels.

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Foundation — 1 task, must be first):
└── T1: Project scaffolding + core types & interfaces [quick]

Wave 2 (Core Infrastructure — 5 parallel tasks):
├── T2: Scanner + Language Detection [deep]
├── T3: Checker Registry + Runner Engine [deep]
├── T4: Scorer (level calculation + gated progression) [deep]
├── T5: LLM Provider Interface + Mock + Fallback [deep]
└── T6: Reporter Interface + JSON Reporter [unspecified-high]

Wave 3 (ALL 40 Criteria — 13 parallel tasks, MAX PARALLEL):
├── T7: lint_config + formatter + type_check [unspecified-high]
├── T8: strict_typing + pre_commit_hooks + naming_consistency [unspecified-high]
├── T9: cyclomatic_complexity + dead_code_detection + duplicate_code_detection [unspecified-high]
├── T10: code_modularization + large_file_detection + tech_debt_tracking [unspecified-high]
├── T11: build_cmd_doc + single_command_setup + deps_pinned [unspecified-high]
├── T12: fast_ci_feedback + release_automation + deployment_frequency [unspecified-high]
├── T13: vcs_cli_tools + agentic_development + automated_pr_review [unspecified-high]
├── T14: build_perf_tracking + feature_flag_infra + release_notes + unused_deps [unspecified-high]
├── T15: unit_tests_exist + unit_tests_runnable + test_naming + test_isolation [unspecified-high]
├── T16: integration_tests_exist + test_coverage_thresholds [unspecified-high]
├── T17: flaky_test_detection + test_performance_tracking [unspecified-high]
├── T18: readme + agents_md + documentation_freshness + skills [unspecified-high]
└── T19: automated_doc_gen + service_flow_documented + api_schema_docs [unspecified-high]

Wave 4 (Presentation + LLM — 6 parallel tasks):
├── T20: HTML Reporter with embedded template [unspecified-high]
├── T21: Text Reporter (non-TTY fallback) [quick]
├── T22: TUI Root Model + Progress View [deep]
├── T23: TUI Report View + Detail View + Browser Opening [visual-engineering]
├── T24: LLM Provider Implementations (OpenAI, Anthropic, Ollama) [deep]
└── T25: CLI Entry Point (cmd/ari, flag parsing, wiring) [unspecified-high]

Wave 5 (Integration + Polish — 3 parallel tasks):
├── T26: Wire LLM into applicable criteria + FallbackEvaluator [deep]
├── T27: End-to-end integration tests + testdata fixtures [deep]
└── T28: README + AGENTS.md + goreleaser setup [writing]

Wave FINAL (Verification — 4 parallel tasks):
├── F1: Plan compliance audit [oracle]
├── F2: Code quality review [unspecified-high]
├── F3: Real manual QA [unspecified-high]
└── F4: Scope fidelity check [deep]
```

**Critical Path**: T1 → T2/T3 → T7-T19 → T25 → T27 → FINAL
**Parallel Speedup**: ~75% faster than sequential
**Max Concurrent**: 13 (Wave 3)

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| T1 | — | T2-T6 | 1 |
| T2 | T1 | T7-T19, T25 | 2 |
| T3 | T1 | T7-T19, T25 | 2 |
| T4 | T1 | T20, T21, T22, T23, T25 | 2 |
| T5 | T1 | T24, T26 | 2 |
| T6 | T1 | T20, T21, T25 | 2 |
| T7-T19 | T2, T3 | T25, T26, T27 | 3 |
| T20 | T4, T6 | T25 | 4 |
| T21 | T4, T6 | T25 | 4 |
| T22 | T4 | T25 | 4 |
| T23 | T4, T22 | T25 | 4 |
| T24 | T5 | T26 | 4 |
| T25 | T2-T4, T6, T7-T23 | T27 | 4 |
| T26 | T5, T7-T19, T24 | T27 | 5 |
| T27 | T25, T26 | F1-F4 | 5 |
| T28 | T25 | F1-F4 | 5 |
| F1-F4 | T27, T28 | — | FINAL |

### Agent Dispatch Summary

| Wave | Tasks | Dispatch |
|------|-------|----------|
| 1 | 1 | T1 → `quick` |
| 2 | 5 | T2-T4 → `deep`, T5 → `deep`, T6 → `unspecified-high` |
| 3 | 13 | T7-T19 → `unspecified-high` |
| 4 | 6 | T20 → `unspecified-high`, T21 → `quick`, T22 → `deep`, T23 → `visual-engineering`, T24 → `deep`, T25 → `unspecified-high` |
| 5 | 3 | T26 → `deep`, T27 → `deep`, T28 → `writing` |
| FINAL | 4 | F1 → `oracle`, F2-F3 → `unspecified-high`, F4 → `deep` |

---

## TODOs

- [ ] 1. Project Scaffolding + Core Types & Interfaces

  **What to do**:
  - Initialize Go module: `go mod init github.com/bbik/ari`
  - Create directory structure at the REPO ROOT (working directory, NOT a nested `ari/` folder):
    ```
    .                              # repo root = go module root
    ├── cmd/ari/main.go            # CLI entry, flag parsing, wiring
    ├── internal/
    │   ├── checker/checker.go
    │   ├── scanner/scanner.go
    │   ├── scorer/scorer.go
    │   ├── llm/llm.go
    │   ├── reporter/reporter.go
    │   └── tui/model.go
    ├── testdata/
    ├── go.mod
    └── go.sum
    ```
    All commands (e.g. `go build ./cmd/ari`, `go test ./...`) run from this repo root.
  - Define core types in `internal/checker/checker.go`:
    - `Language` enum: `Go`, `TypeScript`, `Java`, `Unknown`
    - `Level` constants: `Functional(1)`, `Documented(2)`, `Standardized(3)`, `Optimized(4)`, `Autonomous(5)`
    - `Pillar` enum: `StyleValidation`, `BuildSystem`, `Testing`, `Documentation`
    - `CheckerID` string type for unique criterion identifiers
    - `Result` struct: `Passed bool`, `Evidence string`, `Level Level`, `Skipped bool`, `SkipReason string`, `Mode string` (rule-based/llm)
    - `Checker` interface: `ID() CheckerID`, `Pillar() Pillar`, `Level() Level`, `Name() string`, `Description() string`, `Check(ctx context.Context, repo fs.FS, lang Language) (*Result, error)`
    - `SuggestionProvider` interface: `Suggestion() string` (what to do if criterion fails)
  - Define scanner types in `internal/scanner/scanner.go`:
    - `FileInfo` struct: `Path string`, `Size int64`, `IsDir bool`, `Extension string`
    - `RepoInfo` struct: `Files []FileInfo`, `Language Language`, `IsGitRepo bool`, `RootPath string`
    - `Scanner` interface: `Scan(ctx context.Context, repo fs.FS) (*RepoInfo, error)`
  - Define scorer types in `internal/scorer/scorer.go`:
    - `Score` struct: `Level Level`, `PassRate float64`, `PillarScores map[Pillar]PillarScore`, `Results []Result`
    - `PillarScore` struct: `Pillar Pillar`, `Passed int`, `Total int`, `Rate float64`
  - Define LLM types in `internal/llm/llm.go`:
    - `Provider` interface: `Complete(ctx context.Context, prompt string, opts ...Option) (string, error)`, `Name() string`
    - `Option` functional options for temperature, max tokens, etc.
    - `Evaluator` interface: `Evaluate(ctx context.Context, prompt string) (*EvalResult, error)`
    - `EvalResult` struct: `Passed bool`, `Evidence string`, `Confidence float64`
  - Define reporter types in `internal/reporter/reporter.go`:
    - `Reporter` interface: `Report(ctx context.Context, score *scorer.Score) error`
    - `Format` enum: `JSON`, `HTML`, `Text`
  - `cmd/ari/main.go`: Minimal main with `--help` flag, prints "ari - Agent Readiness Index" and exits
  - Install Bubbletea v2: `go get charm.land/bubbletea/v2`
  - Write tests for type construction and validation (e.g., `Level.String()`, `Language.String()`)

  **Must NOT do**:
  - Do not implement any business logic — types and interfaces only
  - Do not add any third-party dependencies beyond bubbletea v2

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Scaffolding task with no complex logic — just directory creation, type definitions, interface declarations
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `git-master`: Not needed — no git operations beyond initial setup

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1 (solo)
  - **Blocks**: T2, T3, T4, T5, T6 (all Wave 2 tasks)
  - **Blocked By**: None (can start immediately)

  **References**:

  **Pattern References**:
  - Bubbletea v2 import: `charm.land/bubbletea/v2` (NOT `github.com/charmbracelet/bubbletea`)
  - Go module naming convention: standard `github.com/{org}/{repo}` pattern
  - Factory.ai criteria structure: each criterion has ID, pillar, level, pass/fail, evidence

  **External References**:
  - Bubbletea v2 API: https://pkg.go.dev/charm.land/bubbletea/v2 — View() returns tea.View, not string
  - Go fs.FS: https://pkg.go.dev/io/fs — Scanner must accept this interface
  - Factory.ai CockroachDB report: https://factory.ai/agent-readiness/cockroachdb_cockroach — criterion structure reference

  **Acceptance Criteria**:
  - [ ] `go build ./...` compiles without errors
  - [ ] `go test ./...` passes (type construction tests)
  - [ ] `./ari --help` prints usage text
  - [ ] All interface definitions compile and are well-documented
  - [ ] Directory structure matches the specified layout

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Project builds successfully
    Tool: Bash
    Preconditions: Go 1.24+ installed, module initialized
    Steps:
      1. Run `go build ./cmd/ari`
      2. Run `./ari --help`
      3. Verify output contains "ari" and "Agent Readiness Index"
    Expected Result: Binary builds, help text displayed
    Failure Indicators: Compilation errors, missing imports, no output
    Evidence: .sisyphus/evidence/task-1-build.txt

  Scenario: Core types are usable
    Tool: Bash (go test)
    Preconditions: Types defined in internal/ packages
    Steps:
      1. Run `go test ./internal/checker/ -v -run TestTypes`
      2. Verify Level.String() returns "Functional", "Documented", etc.
      3. Verify Language.String() returns "Go", "TypeScript", "Java"
    Expected Result: All type tests pass
    Failure Indicators: Test failures, missing String() methods
    Evidence: .sisyphus/evidence/task-1-types-test.txt
  ```

  **Commit**: YES
  - Message: `feat(core): scaffold ari project with core types and interfaces`
  - Files: `cmd/ari/main.go`, `internal/checker/checker.go`, `internal/scanner/scanner.go`, `internal/scorer/scorer.go`, `internal/llm/llm.go`, `internal/reporter/reporter.go`, `go.mod`, `go.sum`
  - Pre-commit: `go build ./... && go test ./...`

- [ ] 2. Scanner + Language Detection

  **What to do**:
  - Implement `DefaultScanner` in `internal/scanner/scanner.go`:
    - `Scan(ctx context.Context, repo fs.FS) (*RepoInfo, error)` — walks repo using `fs.WalkDir`
    - Collects `FileInfo` for each file (path, size, extension)
    - Filters out common ignore patterns (`.git/`, `node_modules/`, `vendor/`, `.idea/`, `build/`, `dist/`)
    - Sets configurable file limit (default 5000) to handle large repos
    - Detects symlink cycles (log warning, skip)
    - Handles permission errors (log warning, continue scan)
  - Implement `DetectLanguage` in `internal/scanner/language.go`:
    - Detects primary language by presence of signature files:
      - Go: `go.mod` exists → `Language.Go`
      - TypeScript/JavaScript: `package.json` exists → `Language.TypeScript`
      - Java/Kotlin: `pom.xml` OR `build.gradle` OR `build.gradle.kts` exists → `Language.Java`
    - If multiple detected, rank by source file count ratio (`.go` vs `.ts` vs `.java` files)
    - Return `Language.Unknown` if none detected
  - Implement `DetectGitRepo` in `internal/scanner/git.go`:
    - Check for `.git/` directory existence
    - If git repo: extract commit hash, branch name, has local changes (via `git` CLI)
    - If not git repo: set `IsGitRepo = false`, log warning
  - Write comprehensive TDD tests in `internal/scanner/scanner_test.go`:
    - Use `testing/fstest.MapFS` for ALL test fixtures
    - Test: Go repo detection (go.mod present)
    - Test: TypeScript repo detection (package.json with typescript dep)
    - Test: Java repo detection (pom.xml present)
    - Test: Unknown language (no signature files)
    - Test: Mixed language repo (go.mod + package.json → picks primary by file count)
    - Test: Empty directory (no files) → returns empty FileInfo, Language.Unknown
    - Test: Large repo (>5000 files) → respects file limit
    - Test: Ignore patterns (.git, node_modules excluded from results)

  **Must NOT do**:
  - Do not use direct `os.Open` or `os.ReadDir` — always `fs.FS`
  - Do not call GitHub API for any metadata

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core infrastructure with complex logic (file walking, language detection, edge cases) requiring careful TDD
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with T3, T4, T5, T6)
  - **Blocks**: T7-T19 (all checker tasks), T25 (CLI entry point)
  - **Blocked By**: T1

  **References**:

  **Pattern References**:
  - `io/fs.WalkDir` — standard Go pattern for walking filesystems
  - `testing/fstest.MapFS` — in-memory filesystem for tests

  **External References**:
  - Go fs.FS docs: https://pkg.go.dev/io/fs
  - Go fstest docs: https://pkg.go.dev/testing/fstest
  - Factory.ai language detection: detects Go, Python, JavaScript/TypeScript, Java, Rust, Ruby

  **Acceptance Criteria**:
  - [ ] `go test ./internal/scanner/ -v` — all tests pass
  - [ ] Go repo detection works (go.mod → Language.Go)
  - [ ] TypeScript repo detection works (package.json → Language.TypeScript)
  - [ ] Java repo detection works (pom.xml/build.gradle → Language.Java)
  - [ ] Empty directory handled gracefully (no crash)
  - [ ] .git/, node_modules/ excluded from scan results
  - [ ] File limit (5000) respected for large repos

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Scanner correctly detects Go repository
    Tool: Bash (go test)
    Preconditions: MapFS fixture with go.mod, main.go, main_test.go
    Steps:
      1. Run `go test ./internal/scanner/ -v -run TestScanGoRepo`
      2. Verify RepoInfo.Language == Language.Go
      3. Verify RepoInfo.Files contains "main.go" but not ".git/HEAD"
    Expected Result: Language detected as Go, .git excluded
    Failure Indicators: Wrong language, .git files in results
    Evidence: .sisyphus/evidence/task-2-go-detection.txt

  Scenario: Scanner handles empty directory gracefully
    Tool: Bash (go test)
    Preconditions: MapFS fixture with zero files
    Steps:
      1. Run `go test ./internal/scanner/ -v -run TestScanEmptyRepo`
      2. Verify RepoInfo.Language == Language.Unknown
      3. Verify len(RepoInfo.Files) == 0
      4. Verify no panic or error
    Expected Result: Returns empty RepoInfo with Unknown language
    Failure Indicators: Panic, error returned, non-zero file count
    Evidence: .sisyphus/evidence/task-2-empty-repo.txt

  Scenario: Scanner respects file limit
    Tool: Bash (go test)
    Preconditions: MapFS fixture with 6000+ files
    Steps:
      1. Run `go test ./internal/scanner/ -v -run TestScanFileLimitRespected`
      2. Verify len(RepoInfo.Files) <= 5000
    Expected Result: File count capped at limit
    Failure Indicators: All 6000 files returned, timeout
    Evidence: .sisyphus/evidence/task-2-file-limit.txt
  ```

  **Commit**: YES
  - Message: `feat(scanner): add repo scanner with fs.FS and language detection`
  - Files: `internal/scanner/scanner.go`, `internal/scanner/language.go`, `internal/scanner/git.go`, `internal/scanner/scanner_test.go`
  - Pre-commit: `go test ./internal/scanner/...`

- [ ] 3. Checker Registry + Runner Engine

  **What to do**:
  - Implement `Registry` in `internal/checker/registry.go`:
    - `Registry` struct holding `map[CheckerID]Checker`
    - `Register(checker Checker) error` — adds checker, error if duplicate ID
    - `Get(id CheckerID) (Checker, bool)` — lookup by ID
    - `GetByPillar(pillar Pillar) []Checker` — all checkers for a pillar
    - `GetByLevel(level Level) []Checker` — all checkers for a level
    - `All() []Checker` — all registered checkers
    - `NewDefaultRegistry() *Registry` — returns registry with all 40 MVP checkers (initially empty, filled as checkers are implemented)
  - Implement `Runner` in `internal/checker/runner.go`:
    - `Runner` struct with `Registry`, optional `llm.Evaluator`
    - `Run(ctx context.Context, repo fs.FS, repoInfo *scanner.RepoInfo) ([]Result, error)`
    - For each checker in registry:
      1. Check if applicable for detected language (skip if not)
      2. Execute `checker.Check(ctx, repo, lang)`
      3. Handle panics gracefully (recover, mark as failed with evidence)
      4. Collect all results
    - Support context cancellation for long-running scans
    - Log progress (checker N of M) via callback function
  - Write TDD tests in `internal/checker/registry_test.go`:
    - Test: Register checker → retrieve by ID
    - Test: Duplicate registration → returns error
    - Test: GetByPillar returns only matching checkers
    - Test: GetByLevel returns only matching checkers
  - Write TDD tests in `internal/checker/runner_test.go`:
    - Test: Runner executes all checkers and collects results
    - Test: Runner skips inapplicable checkers (marks as skipped with reason)
    - Test: Runner handles checker panic without crashing
    - Test: Runner respects context cancellation
    - Use mock checkers (implementing Checker interface with canned responses)

  **Must NOT do**:
  - Do not use `init()` for checker registration — explicit `NewDefaultRegistry()` only
  - Do not implement actual checkers here — only the infrastructure

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core orchestration logic with complex error handling, skip logic, and panic recovery
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with T2, T4, T5, T6)
  - **Blocks**: T7-T19 (all checker tasks), T25 (CLI entry point)
  - **Blocked By**: T1

  **References**:

  **Pattern References**:
  - pingcap/tidb precheck pattern: `Checker` interface with `Check(ctx) (*CheckResult, error)` + `GetCheckItemID()`

  **External References**:
  - Factory.ai evaluation: each criterion returns binary pass/fail + evidence string

  **Acceptance Criteria**:
  - [ ] `go test ./internal/checker/ -v` — all tests pass
  - [ ] Registry correctly stores and retrieves checkers by ID, pillar, level
  - [ ] Runner executes all applicable checkers and returns results
  - [ ] Inapplicable checkers are skipped with reason
  - [ ] Checker panics don't crash the runner

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Registry stores and retrieves checkers
    Tool: Bash (go test)
    Preconditions: Mock checkers registered
    Steps:
      1. Run `go test ./internal/checker/ -v -run TestRegistryRegisterAndGet`
      2. Register 3 mock checkers with different pillars
      3. Verify Get returns each by ID
      4. Verify GetByPillar returns only matching
    Expected Result: All lookups return correct checkers
    Failure Indicators: Wrong checker returned, nil result
    Evidence: .sisyphus/evidence/task-3-registry.txt

  Scenario: Runner handles checker panic gracefully
    Tool: Bash (go test)
    Preconditions: Mock checker that panics
    Steps:
      1. Run `go test ./internal/checker/ -v -run TestRunnerPanicRecovery`
      2. Register a mock checker that panics in Check()
      3. Verify runner doesn't crash
      4. Verify result shows failed with panic evidence
    Expected Result: Runner recovers, result marked as failed
    Failure Indicators: Test crashes, unrecovered panic
    Evidence: .sisyphus/evidence/task-3-panic-recovery.txt
  ```

  **Commit**: YES
  - Message: `feat(checker): add checker registry and runner engine`
  - Files: `internal/checker/registry.go`, `internal/checker/runner.go`, `internal/checker/registry_test.go`, `internal/checker/runner_test.go`
  - Pre-commit: `go test ./internal/checker/...`

- [ ] 4. Scorer — Level Calculation + Gated Progression

  **What to do**:
  - Implement `Scorer` in `internal/scorer/scorer.go`:
    - `Calculate(results []*checker.Result) *Score` — main scoring function
    - Groups results by pillar → calculates per-pillar pass rate
    - Groups results by level → calculates per-level pass rate
    - Applies gated progression: To unlock Level N, must pass ≥80% at Level N AND all previous levels
    - Skip results don't count toward total (numerator OR denominator)
    - Returns `Score` with: `Level`, `PassRate`, `PillarScores`, `LevelScores`
  - Level score mapping (from interview + Factory.ai baseline):
    - Level 1 (7 criteria): readme, lint_config, formatter, type_check, unit_tests_exist, build_cmd_doc, deps_pinned
    - Level 2 (8 criteria): agents_md, documentation_freshness, strict_typing, pre_commit_hooks, naming_consistency, single_command_setup, unit_tests_runnable, test_naming_conventions
    - Level 3 (9 criteria): integration_tests_exist, test_coverage_thresholds, service_flow_documented, api_schema_docs, fast_ci_feedback, release_automation, agentic_development, cyclomatic_complexity, dead_code_detection
    - Level 4 (12 criteria): automated_doc_generation, skills, automated_pr_review, build_performance_tracking, feature_flag_infrastructure, release_notes_automation, flaky_test_detection, test_performance_tracking, test_isolation, duplicate_code_detection, code_modularization, deployment_frequency
    - Level 5 (4 criteria): tech_debt_tracking, large_file_detection, unused_dependencies_detection, vcs_cli_tools
  - Write comprehensive TDD tests:
    - Test: 100% pass at L1 → Level 1 achieved
    - Test: 80% pass at L1 → Level 1 achieved (threshold)
    - Test: 79% pass at L1 → Level 0 (below threshold)
    - Test: 100% L1 + 80% L2 → Level 2 achieved
    - Test: 100% L1 + 0% L2 → Level 1 (gated by L2 failure)
    - Test: 100% L1 + 100% L2 + 100% L3 + 80% L4 → Level 4
    - Test: Skipped results excluded from calculation
    - Test: Empty results → Level 0

  **Must NOT do**:
  - Do not hardcode criteria IDs in scorer — use `Result.Level` field
  - Do not implement weighted scoring — binary pass/fail with 80% threshold only

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core algorithm with precise threshold logic requiring careful edge case testing
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with T2, T3, T5, T6)
  - **Blocks**: T20-T23, T25 (reporters and TUI need scores)
  - **Blocked By**: T1

  **References**:

  **External References**:
  - Factory.ai scoring: "To unlock a level, you must pass 80% of criteria from that level and all previous levels"
  - Factory.ai CockroachDB: L1 100%, L2 100%, L3 100%, L4 70% → Level 3 (L4 below 80%)

  **Acceptance Criteria**:
  - [ ] `go test ./internal/scorer/ -v` — all tests pass
  - [ ] 80% threshold correctly gates progression
  - [ ] Gated progression works (can't skip levels)
  - [ ] Skipped results excluded from denominator
  - [ ] Per-pillar scores calculated correctly

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Scorer applies 80% gated progression correctly
    Tool: Bash (go test)
    Preconditions: Hardcoded result slices with known pass/fail counts
    Steps:
      1. Run `go test ./internal/scorer/ -v -run TestGatedProgression`
      2. Input: 6/7 L1 pass (85%), 5/8 L2 pass (62%)
      3. Verify Score.Level == 1 (L1 passes, L2 fails threshold)
    Expected Result: Level 1 — L2 blocked by 62% < 80%
    Failure Indicators: Level 2 reported, wrong percentage
    Evidence: .sisyphus/evidence/task-4-gated-progression.txt

  Scenario: Scorer handles all-skip results
    Tool: Bash (go test)
    Preconditions: Results where all criteria are skipped
    Steps:
      1. Run `go test ./internal/scorer/ -v -run TestAllSkipped`
      2. Input: all results have Skipped=true
      3. Verify Score.Level == 0, Score.PassRate == 0
    Expected Result: Level 0 with 0% pass rate, no division by zero
    Failure Indicators: Panic (div/0), non-zero level
    Evidence: .sisyphus/evidence/task-4-all-skipped.txt
  ```

  **Commit**: YES
  - Message: `feat(scorer): add 5-level maturity scoring with gated progression`
  - Files: `internal/scorer/scorer.go`, `internal/scorer/scorer_test.go`
  - Pre-commit: `go test ./internal/scorer/...`

- [ ] 5. LLM Provider Interface + Mock + Fallback Evaluator

  **What to do**:
  - Implement provider interface in `internal/llm/llm.go` (types already defined in T1):
    - `Option` type with functional options: `WithTemperature(float64)`, `WithMaxTokens(int)`, `WithJSONOutput(bool)`
    - `Config` struct: `Provider string`, `APIKey string`, `Model string`, `BaseURL string`
    - `NewProviderFromConfig(cfg Config) (Provider, error)` — factory function
  - Implement `MockProvider` in `internal/llm/mock.go`:
    - Configurable responses via function field: `CompleteFunc func(ctx, prompt, opts) (string, error)`
    - Default: returns canned JSON response `{"passed": true, "evidence": "mock", "confidence": 0.9}`
    - For use in all checker tests
  - Implement `FallbackEvaluator` in `internal/llm/fallback.go`:
    - Takes `primary Evaluator` (LLM-based) and `fallback Evaluator` (rule-based)
    - Try primary → if error or timeout → use fallback
    - Record which mode was used in result (`Mode: "llm"` or `Mode: "rule-based"`)
    - If no LLM configured (`primary == nil`), always use fallback
  - Implement `RuleBasedEvaluator` in `internal/llm/rule_evaluator.go`:
    - Simple keyword/pattern matching evaluator for criteria that need interpretation
    - Less accurate than LLM but always works without API key
  - Configuration via environment variables:
    - `ARI_LLM_PROVIDER` — "openai", "anthropic", "ollama", or empty (rule-based only)
    - `ARI_API_KEY` — API key for the chosen provider
    - `ARI_LLM_MODEL` — model name (default per provider)
    - `ARI_LLM_BASE_URL` — custom base URL (for Ollama or proxies)
  - Write TDD tests:
    - Test: MockProvider returns configured responses
    - Test: FallbackEvaluator uses primary when available
    - Test: FallbackEvaluator falls back when primary errors
    - Test: FallbackEvaluator falls back when primary is nil
    - Test: RuleBasedEvaluator provides basic evaluation
    - Test: Config from env vars is correctly parsed

  **Must NOT do**:
  - Do not implement actual OpenAI/Anthropic/Ollama clients here (that's T24)
  - Do not create config file parsing (`.ari.yaml`) — env vars only
  - Do not implement streaming — simple request/response only

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Interface design + fallback pattern requires careful error handling and mode tracking
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with T2, T3, T4, T6)
  - **Blocks**: T24 (LLM provider implementations), T26 (LLM wiring)
  - **Blocked By**: T1

  **References**:

  **Pattern References**:
  - smhanov/llmhub: Provider interface with `Complete(ctx, prompt, opts) (string, error)`
  - noperator/siftrank: `LLMProvider` interface pattern
  - Functional options pattern: `type Option func(*config)`

  **Acceptance Criteria**:
  - [ ] `go test ./internal/llm/ -v` — all tests pass
  - [ ] MockProvider works with configurable responses
  - [ ] FallbackEvaluator correctly falls back on error
  - [ ] FallbackEvaluator records mode ("llm" vs "rule-based") in results
  - [ ] Config correctly reads from env vars

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: FallbackEvaluator uses primary when available
    Tool: Bash (go test)
    Preconditions: MockProvider configured as primary
    Steps:
      1. Run `go test ./internal/llm/ -v -run TestFallbackUsesPrimary`
      2. Create FallbackEvaluator with mock primary + rule fallback
      3. Call Evaluate() — verify primary was used
      4. Verify result.Mode == "llm"
    Expected Result: Primary evaluator used, mode is "llm"
    Failure Indicators: Fallback used instead, wrong mode
    Evidence: .sisyphus/evidence/task-5-fallback-primary.txt

  Scenario: FallbackEvaluator falls back on error
    Tool: Bash (go test)
    Preconditions: MockProvider that returns error
    Steps:
      1. Run `go test ./internal/llm/ -v -run TestFallbackOnError`
      2. Create FallbackEvaluator with error-returning primary
      3. Call Evaluate() — verify fallback was used
      4. Verify result.Mode == "rule-based"
    Expected Result: Fallback evaluator used, mode is "rule-based"
    Failure Indicators: Error propagated, primary mode reported
    Evidence: .sisyphus/evidence/task-5-fallback-error.txt
  ```

  **Commit**: YES
  - Message: `feat(llm): add multi-provider LLM interface with mock and fallback`
  - Files: `internal/llm/llm.go`, `internal/llm/mock.go`, `internal/llm/fallback.go`, `internal/llm/rule_evaluator.go`, `internal/llm/llm_test.go`
  - Pre-commit: `go test ./internal/llm/...`

- [ ] 6. Reporter Interface + JSON Reporter

  **What to do**:
  - Implement reporter interface in `internal/reporter/reporter.go` (types from T1):
    - `Report` struct: full report data model
      - `RepoPath string`, `Language string`, `IsGitRepo bool`
      - `CommitHash string`, `Branch string` (if git repo)
      - `Score *scorer.Score`
      - `CriteriaResults []CriterionReport` (per-criterion detail)
      - `Suggestions []Suggestion` (actionable recommendations)
      - `GeneratedAt time.Time`, `AriVersion string`
    - `CriterionReport` struct: `ID string`, `Name string`, `Pillar string`, `Level int`, `Passed bool`, `Skipped bool`, `SkipReason string`, `Evidence string`, `Mode string`, `Suggestion string`
    - `Suggestion` struct: `CriterionID string`, `Title string`, `Description string`, `Priority string`
    - `BuildReport(repoInfo, score, results) *Report` — assembles report from raw data
  - Implement `JSONReporter` in `internal/reporter/json.go`:
    - `Report(ctx context.Context, report *Report, w io.Writer) error`
    - Pretty-prints JSON with indentation
    - Includes all fields: score, level, criteria results, suggestions
  - Write TDD tests:
    - Test: BuildReport correctly assembles from components
    - Test: JSONReporter outputs valid JSON
    - Test: JSON output includes all required fields (level, passRate, criteria, suggestions)
    - Test: JSON output can be unmarshalled back into Report struct
    - Test: Suggestions are generated for failing criteria only

  **Must NOT do**:
  - Do not implement HTML or Text reporters here (T20, T21)
  - Do not add file writing logic — reporters write to `io.Writer`

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Data modeling + serialization with clear structure
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with T2, T3, T4, T5)
  - **Blocks**: T20 (HTML reporter), T21 (Text reporter), T25 (CLI entry point)
  - **Blocked By**: T1

  **References**:

  **External References**:
  - Factory.ai report structure: level, passRate, pillar scores, per-criterion results with evidence

  **Acceptance Criteria**:
  - [ ] `go test ./internal/reporter/ -v` — all tests pass
  - [ ] JSON output is valid and pretty-printed
  - [ ] All required fields present in JSON output
  - [ ] JSON can be round-tripped (marshal → unmarshal → compare)
  - [ ] Suggestions only generated for failing criteria

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: JSON reporter produces valid output
    Tool: Bash (go test)
    Preconditions: Mock report with known data
    Steps:
      1. Run `go test ./internal/reporter/ -v -run TestJSONReporterOutput`
      2. Build report with 2 passed, 1 failed criteria
      3. Render to JSON, parse with json.Unmarshal
      4. Verify level, passRate, criteria count, suggestion count
    Expected Result: Valid JSON with 3 criteria, 1 suggestion
    Failure Indicators: Invalid JSON, missing fields, wrong counts
    Evidence: .sisyphus/evidence/task-6-json-output.txt

  Scenario: JSON reporter handles empty results
    Tool: Bash (go test)
    Preconditions: Report with zero criteria results
    Steps:
      1. Run `go test ./internal/reporter/ -v -run TestJSONReporterEmpty`
      2. Build report with no results
      3. Verify JSON output has empty arrays, level 0
    Expected Result: Valid JSON with empty criteria array, level 0
    Failure Indicators: Null instead of empty array, error
    Evidence: .sisyphus/evidence/task-6-json-empty.txt
  ```

  **Commit**: YES
  - Message: `feat(reporter): add reporter interface and JSON reporter`
  - Files: `internal/reporter/reporter.go`, `internal/reporter/json.go`, `internal/reporter/reporter_test.go`
  - Pre-commit: `go test ./internal/reporter/...`

- [ ] 7. Style & Validation — lint_config + formatter + type_check

  **What to do**:
  - Implement 3 checkers in `internal/checker/style/`:
  - **lint_config** (Level 1): Check for linter configuration files
    - Go: `.golangci.yml`, `.golangci.yaml`, `.golangci.json`, `.golangci.toml`
    - TypeScript: `.eslintrc`, `.eslintrc.js`, `.eslintrc.json`, `.eslintrc.yml`, `eslint.config.js`, `eslint.config.mjs`, `biome.json`, `biome.jsonc`
    - Java: `checkstyle.xml`, `pmd.xml`, `.spotbugs.xml`, `spotless` in build.gradle
    - Evidence: "Found [file] — [tool] configured"
    - Suggestion: "Add a linter configuration. For Go: `golangci-lint init`. For TS: `npm init @eslint/config`"
  - **formatter** (Level 1): Check for formatter configuration
    - Go: `gofmt` is built-in (always pass for Go repos), check for `gofumpt` or `crlfmt` config
    - TypeScript: `.prettierrc`, `.prettierrc.js`, `.prettierrc.json`, `prettier.config.js`, `biome.json` with formatter section
    - Java: `spotless` in build.gradle, `.editorconfig` with java settings
    - Evidence: "Found [file] — formatter [name] configured"
    - Suggestion: "Add a code formatter. For Go: built-in gofmt. For TS: `npm install -D prettier`"
  - **type_check** (Level 1): Check that language has type checking
    - Go: Always pass (statically typed by nature)
    - TypeScript: `tsconfig.json` exists
    - Java: Always pass (statically typed by nature)
    - JavaScript-only (no TS): Fail — no type checking
    - Evidence: "Go is statically typed" / "Found tsconfig.json" / "JavaScript without TypeScript — no type checking"
  - All checkers implement `Checker` interface from T1
  - All use `fs.FS` for file existence checks (never `os.Stat`)
  - Write TDD tests per checker per language (pass + fail + skip cases)
  - Use `testing/fstest.MapFS` for all test fixtures

  **Must NOT do**:
  - Do not run the linter/formatter — only check configuration exists
  - Do not read file contents deeply — existence check + minimal parsing

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: File existence pattern matching across 3 languages, well-defined but needs thoroughness
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T8-T19)
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:

  **Pattern References**:
  - Factory.ai CockroachDB: `lint_config` 1/1 — "Go custom linter via ./dev lint command and ESLint configured for UI workspaces"
  - Factory.ai CockroachDB: `formatter` 1/1 — "crlfmt custom formatter for Go code, Prettier configured for UI workspaces"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented, each with Checker interface
  - [ ] Tests pass for Go, TypeScript, Java repos (pass and fail cases)
  - [ ] Evidence strings are descriptive
  - [ ] Suggestion text provided for each failing criterion

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: lint_config detects Go linter
    Tool: Bash (go test)
    Preconditions: MapFS with go.mod + .golangci.yml
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestLintConfigGo`
      2. Verify Result.Passed == true
      3. Verify Result.Evidence contains "golangci"
    Expected Result: Pass — golangci-lint detected
    Failure Indicators: Failed when config exists, wrong evidence
    Evidence: .sisyphus/evidence/task-7-lint-config-go.txt

  Scenario: type_check fails for JS-only repo
    Tool: Bash (go test)
    Preconditions: MapFS with package.json (no typescript dep), no tsconfig.json
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestTypeCheckJSOnly`
      2. Verify Result.Passed == false
      3. Verify Result.Evidence contains "JavaScript without TypeScript"
    Expected Result: Fail — no type checking
    Failure Indicators: Passes when no type checking exists
    Evidence: .sisyphus/evidence/task-7-type-check-js.txt
  ```

  **Commit**: YES (groups with T8, T9, T10)
  - Message: `feat(style): add style & validation checkers (lint, formatter, type_check)`
  - Files: `internal/checker/style/lint_config.go`, `internal/checker/style/formatter.go`, `internal/checker/style/type_check.go`, `internal/checker/style/*_test.go`
  - Pre-commit: `go test ./internal/checker/style/...`

- [ ] 8. Style & Validation — strict_typing + pre_commit_hooks + naming_consistency

  **What to do**:
  - Implement 3 checkers in `internal/checker/style/`:
  - **strict_typing** (Level 2): Check for strict type enforcement
    - Go: Always pass (strict by nature)
    - TypeScript: Parse `tsconfig.json`, check `"strict": true` in `compilerOptions`
    - Java: Check for `-Xlint` flags in build config, or `@NonNull` annotation usage
    - Evidence: "tsconfig.json has strict: true" / "Go enforces strict typing by default"
    - Suggestion: "Enable strict mode. For TS: set `strict: true` in tsconfig.json"
  - **pre_commit_hooks** (Level 2): Check for pre-commit hook configuration
    - Language-agnostic: `.pre-commit-config.yaml`, `.husky/` directory, `lefthook.yml`, `lint-staged` in package.json
    - Also check: `.git/hooks/pre-commit` exists (custom hook)
    - Evidence: "Found .pre-commit-config.yaml — pre-commit hooks configured"
    - Suggestion: "Add pre-commit hooks. Install: `pip install pre-commit && pre-commit install` or use Husky for JS projects"
  - **naming_consistency** (Level 2): LLM-assisted criterion
    - Rule-based fallback: Check if project follows language conventions:
      - Go: exported functions are PascalCase (check a sample of .go files)
      - TypeScript: files are camelCase or kebab-case (check file names)
      - Java: classes are PascalCase (check file names match class convention)
    - LLM mode: Send sample of 5–10 function/type names to LLM for consistency evaluation
    - Uses `FallbackEvaluator` from T5
    - Evidence: "Naming follows Go conventions (PascalCase exports)" / "Inconsistent naming detected: mix of camelCase and snake_case"
  - Write TDD tests per checker per language
  - For naming_consistency: test with MockProvider for LLM path, test rule-based fallback

  **Must NOT do**:
  - Do not execute pre-commit hooks — only check configuration exists
  - Do not deeply parse all source files for naming — sample-based check only

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Mix of file detection and light content parsing, includes first LLM-assisted criterion
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T7, T9-T19)
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3, T5 (needs FallbackEvaluator for naming_consistency)

  **References**:

  **Pattern References**:
  - Factory.ai CockroachDB: `strict_typing` 1/1 — "Go enforces strict typing by default, TypeScript UI has strict: true in tsconfig.json"
  - Factory.ai CockroachDB: `pre_commit_hooks` 0/1 — "No .pre-commit-config.yaml or Husky configuration found"
  - Factory.ai CockroachDB: `naming_consistency` 1/1 — "CLAUDE.md documents naming conventions, Go stdlib naming enforced"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests
  - [ ] naming_consistency works in both LLM and rule-based modes
  - [ ] strict_typing correctly parses tsconfig.json for strict flag
  - [ ] pre_commit_hooks detects multiple hook systems

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: strict_typing detects TypeScript strict mode
    Tool: Bash (go test)
    Preconditions: MapFS with tsconfig.json containing {"compilerOptions":{"strict":true}}
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestStrictTypingTS`
      2. Verify Result.Passed == true
    Expected Result: Pass — strict mode enabled
    Failure Indicators: Fails when strict is true
    Evidence: .sisyphus/evidence/task-8-strict-typing.txt

  Scenario: naming_consistency uses rule-based fallback when no LLM
    Tool: Bash (go test)
    Preconditions: MapFS with Go files, no LLM configured
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestNamingConsistencyFallback`
      2. Verify Result.Mode == "rule-based"
      3. Verify result includes evidence about naming patterns
    Expected Result: Rule-based evaluation completes, mode recorded
    Failure Indicators: Error when no LLM, mode not recorded
    Evidence: .sisyphus/evidence/task-8-naming-fallback.txt
  ```

  **Commit**: YES (groups with T7, T9, T10)
  - Message: `feat(style): add strict_typing, pre_commit_hooks, naming_consistency checkers`
  - Files: `internal/checker/style/strict_typing.go`, `internal/checker/style/pre_commit_hooks.go`, `internal/checker/style/naming_consistency.go`, `internal/checker/style/*_test.go`
  - Pre-commit: `go test ./internal/checker/style/...`

- [ ] 9. Style & Validation — cyclomatic_complexity + dead_code_detection + duplicate_code_detection

  **What to do**:
  - Implement 3 checkers in `internal/checker/style/`:
  - **cyclomatic_complexity** (Level 3): Check for complexity analysis tools in config
    - Go: `gocyclo` or `gocognit` in `.golangci.yml` linter list, or `cyclop` linter
    - TypeScript: `complexity` rule in ESLint config, or `sonarjs/cognitive-complexity`
    - Java: PMD complexity rules, SonarQube config
    - Evidence: "Found gocyclo in .golangci.yml linter configuration"
    - Suggestion: "Add complexity analysis. For Go: add `gocyclo` to golangci-lint. For TS: enable `complexity` ESLint rule"
  - **dead_code_detection** (Level 3): Check for dead code detection tools
    - Go: `deadcode` or `unused` in golangci-lint config
    - TypeScript: `knip` in package.json devDeps, or `ts-prune` config
    - Java: ProGuard, SpotBugs UnusedCode detector
    - Evidence: "Found knip in devDependencies — dead code detection configured"
  - **duplicate_code_detection** (Level 4): Check for duplicate code detection tools
    - Language-agnostic: `jscpd` config (`.jscpd.json`), PMD CPD in CI config
    - Go: `dupl` in golangci-lint
    - Evidence: "Found dupl linter in .golangci.yml"
  - All checkers: check for tool presence in config files, not execution
  - Write TDD tests with MapFS fixtures

  **Must NOT do**:
  - Do not run complexity analysis — only check tool configuration exists
  - Do not parse full source files for duplicates

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Config file parsing pattern, consistent with T7
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T7, T8, T10-T19)
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:

  **Pattern References**:
  - Factory.ai CockroachDB: `cyclomatic_complexity` 0/1 — "No evidence of complexity analysis tools like gocyclo"
  - Factory.ai CockroachDB: `dead_code_detection` 0/1 — "No dead code detection tools found in CI workflows"
  - Factory.ai CockroachDB: `duplicate_code_detection` 0/1 — "No duplicate code detection tools found in CI workflows"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests per language
  - [ ] Config parsing correctly identifies tools in golangci-lint/eslint configs
  - [ ] Evidence strings are specific about which tool was found

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: cyclomatic_complexity detects gocyclo in golangci-lint
    Tool: Bash (go test)
    Preconditions: MapFS with .golangci.yml containing gocyclo linter
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestCyclomaticComplexityGo`
      2. Verify Result.Passed == true, evidence mentions gocyclo
    Expected Result: Pass — gocyclo detected in config
    Failure Indicators: Fails to parse YAML, wrong evidence
    Evidence: .sisyphus/evidence/task-9-cyclomatic.txt

  Scenario: All 3 checkers fail when no tools configured
    Tool: Bash (go test)
    Preconditions: MapFS with go.mod only, no tool configs
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestAnalysisToolsMissing`
      2. Verify all 3 return Result.Passed == false
      3. Verify each has a meaningful suggestion
    Expected Result: All fail with suggestions
    Failure Indicators: Any pass without tool config
    Evidence: .sisyphus/evidence/task-9-missing-tools.txt
  ```

  **Commit**: YES (groups with T7, T8, T10)
  - Message: `feat(style): add cyclomatic_complexity, dead_code, duplicate_code checkers`
  - Files: `internal/checker/style/cyclomatic_complexity.go`, `internal/checker/style/dead_code_detection.go`, `internal/checker/style/duplicate_code_detection.go`, `internal/checker/style/*_test.go`
  - Pre-commit: `go test ./internal/checker/style/...`

- [ ] 10. Style & Validation — code_modularization + large_file_detection + tech_debt_tracking

  **What to do**:
  - Implement 3 checkers in `internal/checker/style/`:
  - **code_modularization** (Level 4): LLM-assisted criterion
    - Rule-based fallback: Check for module boundary enforcement tools
      - Go: Check for `depguard` in golangci-lint, or `internal/` directory usage pattern
      - TypeScript: Check for `@nx/enforce-module-boundaries` or import restrictions in ESLint
      - Java: Check for module-info.java (Java modules) or multi-module Maven/Gradle
    - LLM mode: Send directory structure to LLM for modularization assessment
    - Uses `FallbackEvaluator` from T5
  - **large_file_detection** (Level 5): Check for large file prevention
    - Language-agnostic: `.gitattributes` with LFS config, or `pre-commit` hook checking file sizes
    - Check for `.gitattributes` patterns like `*.bin filter=lfs`, `*.zip filter=lfs`
    - Evidence: "Found .gitattributes with Git LFS configuration"
    - Suggestion: "Add Git LFS for large files: `git lfs install && git lfs track '*.bin'`"
  - **tech_debt_tracking** (Level 5): Check for TODO/FIXME tracking
    - Check for: TODO scanner in CI (grep-based workflow step), or linter rules enforcing `TODO(TICKET-123)` format
    - Go: `godot` linter, `nolintlint` in golangci-lint
    - TypeScript: `no-warning-comments` ESLint rule
    - Evidence: "Found TODO tracking enforcement in CI workflow"
    - Suggestion: "Add tech debt tracking. Use linter rules to enforce `TODO(JIRA-123)` format"
  - Write TDD tests with MapFS fixtures

  **Must NOT do**:
  - Do not scan all source files for TODOs — check for tracking tool configuration only
  - Do not implement actual modularization analysis

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Mix of file detection and LLM-assisted evaluation
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T7-T9, T11-T19)
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3, T5

  **References**:

  **Pattern References**:
  - Factory.ai CockroachDB: `code_modularization` 0/1 — "No module boundary enforcement tools found"
  - Factory.ai CockroachDB: `large_file_detection` 1/1 — ".gitattributes exists for LFS"
  - Factory.ai CockroachDB: `tech_debt_tracking` 0/1 — "No TODO/FIXME scanner in CI workflows"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests
  - [ ] code_modularization works in both LLM and rule-based modes
  - [ ] large_file_detection correctly parses .gitattributes
  - [ ] tech_debt_tracking checks CI workflows for TODO scanning

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: large_file_detection finds Git LFS config
    Tool: Bash (go test)
    Preconditions: MapFS with .gitattributes containing "*.bin filter=lfs"
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestLargeFileDetection`
      2. Verify Result.Passed == true
    Expected Result: Pass — LFS configured
    Failure Indicators: Fails to parse .gitattributes
    Evidence: .sisyphus/evidence/task-10-large-file.txt

  Scenario: code_modularization uses fallback without LLM
    Tool: Bash (go test)
    Preconditions: MapFS with Go project using internal/ pattern
    Steps:
      1. Run `go test ./internal/checker/style/ -v -run TestCodeModularizationFallback`
      2. Verify Result.Mode == "rule-based"
    Expected Result: Rule-based evaluation completes
    Failure Indicators: Error, wrong mode
    Evidence: .sisyphus/evidence/task-10-modularization.txt
  ```

  **Commit**: YES (groups with T7, T8, T9)
  - Message: `feat(style): add code_modularization, large_file_detection, tech_debt_tracking checkers`
  - Files: `internal/checker/style/code_modularization.go`, `internal/checker/style/large_file_detection.go`, `internal/checker/style/tech_debt_tracking.go`, `internal/checker/style/*_test.go`
  - Pre-commit: `go test ./internal/checker/style/...`

- [ ] 11. Build System — build_cmd_doc + single_command_setup + deps_pinned

  **What to do**:
  - Implement 3 checkers in `internal/checker/build/`:
  - **build_cmd_doc** (Level 1): Check if build commands are documented
    - Search README.md, CONTRIBUTING.md, Makefile for build commands
    - Go: Look for `go build`, `make build`, `./dev build` in docs
    - TypeScript: Look for `npm run build`, `yarn build`, `pnpm build` in docs
    - Java: Look for `mvn`, `gradle`, `./gradlew` in docs
    - Evidence: "Build command documented in README.md: `make build`"
    - Suggestion: "Document build commands in README.md. Add a 'Getting Started' section with build instructions"
  - **single_command_setup** (Level 2): Check for one-command setup
    - Check for: `Makefile` with `setup`/`install` target, `script/setup`, `script/bootstrap`, `docker-compose.yml`, `devcontainer.json`
    - Also check README for "quick start" or "getting started" with single command
    - Evidence: "Found Makefile with `setup` target — single command setup available"
    - Suggestion: "Add a single setup command. Create a Makefile with `make setup` or a `script/setup` script"
  - **deps_pinned** (Level 1): Check for dependency lock files
    - Go: `go.sum` exists
    - TypeScript: `package-lock.json`, `yarn.lock`, `pnpm-lock.yaml`, `bun.lockb` exists
    - Java: `gradle.lockfile`, `pom.xml` with `<dependencyManagement>` version pinning
    - Evidence: "Found go.sum — dependencies pinned with exact versions"
    - Suggestion: "Pin dependencies. For Go: commit go.sum. For npm: `npm install` and commit package-lock.json"
  - Write TDD tests with MapFS fixtures

  **Must NOT do**:
  - Do not execute build commands — only check documentation and config exist
  - Do not verify build actually works

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Content search in docs + file existence checks
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T7-T10, T12-T19)
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:

  **Pattern References**:
  - Factory.ai CockroachDB: `build_cmd_doc` 1/1 — "CLAUDE.md documents ./dev build cockroach"
  - Factory.ai CockroachDB: `single_command_setup` 1/1 — "CLAUDE.md documents ./dev doctor for environment verification"
  - Factory.ai CockroachDB: `deps_pinned` 1/1 — "go.sum lockfile committed with exact dependency versions"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests per language
  - [ ] build_cmd_doc searches README content for build commands
  - [ ] deps_pinned detects lock files per language

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: deps_pinned detects go.sum
    Tool: Bash (go test)
    Preconditions: MapFS with go.mod + go.sum
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestDepsPinnedGo`
      2. Verify Result.Passed == true, evidence mentions go.sum
    Expected Result: Pass — go.sum detected
    Evidence: .sisyphus/evidence/task-11-deps-pinned.txt

  Scenario: build_cmd_doc fails when README has no build instructions
    Tool: Bash (go test)
    Preconditions: MapFS with README.md containing only project description, no build commands
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestBuildCmdDocMissing`
      2. Verify Result.Passed == false
      3. Verify suggestion mentions documenting build commands
    Expected Result: Fail — no build docs
    Evidence: .sisyphus/evidence/task-11-build-doc-missing.txt
  ```

  **Commit**: YES (groups with T12, T13, T14)
  - Message: `feat(build): add build_cmd_doc, single_command_setup, deps_pinned checkers`
  - Files: `internal/checker/build/build_cmd_doc.go`, `internal/checker/build/single_command_setup.go`, `internal/checker/build/deps_pinned.go`, `internal/checker/build/*_test.go`
  - Pre-commit: `go test ./internal/checker/build/...`

- [ ] 12. Build System — fast_ci_feedback + release_automation + deployment_frequency

  **What to do**:
  - Implement 3 checkers in `internal/checker/build/`:
  - **fast_ci_feedback** (Level 3): Check for CI configuration with fast feedback
    - Check for CI config files: `.github/workflows/*.yml`, `.gitlab-ci.yml`, `Jenkinsfile`, `.circleci/config.yml`
    - Basic validation: CI config exists and has test/lint steps
    - Evidence: "Found 3 GitHub Actions workflows — CI configured"
    - Suggestion: "Set up CI/CD. For GitHub: create `.github/workflows/ci.yml` with test and lint steps"
  - **release_automation** (Level 3): Check for automated release process
    - Check for: release workflow in CI, `goreleaser.yml`, `release-please` config, `semantic-release` config
    - Also check: `Makefile` with `release` target
    - Evidence: "Found .goreleaser.yml — release automation configured"
  - **deployment_frequency** (Level 4): Check for regular releases
    - Run `git tag --sort=-creatordate` to check recent tags
    - Pass if: at least 1 tag/release in last 90 days
    - Skip if: not a git repo
    - Evidence: "Last release: v2.1.0 (15 days ago) — regular deployment cadence"
    - Note: This criterion uses `os/exec` with `git` CLI, not `fs.FS`
  - Write TDD tests (deployment_frequency needs git mock via test helper)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:
  - Factory.ai CockroachDB: `fast_ci_feedback` 1/1 — "Essential CI workflow shows checks completing in 6-8 minutes"
  - Factory.ai CockroachDB: `release_automation` 1/1 — "update_releases.yaml workflow automates release file updates"
  - Factory.ai CockroachDB: `deployment_frequency` 1/1 — "Regular releases via gh release list"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests
  - [ ] CI config detection works across GitHub Actions, GitLab CI, etc.
  - [ ] deployment_frequency gracefully handles non-git repos (skip)

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: fast_ci_feedback detects GitHub Actions
    Tool: Bash (go test)
    Preconditions: MapFS with .github/workflows/ci.yml
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestFastCIFeedback`
      2. Verify Result.Passed == true
    Expected Result: Pass — CI config found
    Evidence: .sisyphus/evidence/task-12-ci-feedback.txt

  Scenario: deployment_frequency skips for non-git repo
    Tool: Bash (go test)
    Preconditions: RepoInfo with IsGitRepo == false
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestDeployFreqNoGit`
      2. Verify Result.Skipped == true
    Expected Result: Skipped — not a git repo
    Evidence: .sisyphus/evidence/task-12-deploy-skip.txt
  ```

  **Commit**: YES (groups with T11, T13, T14)
  - Message: `feat(build): add fast_ci_feedback, release_automation, deployment_frequency checkers`
  - Files: `internal/checker/build/fast_ci_feedback.go`, `internal/checker/build/release_automation.go`, `internal/checker/build/deployment_frequency.go`, `internal/checker/build/*_test.go`
  - Pre-commit: `go test ./internal/checker/build/...`

- [ ] 13. Build System — vcs_cli_tools + agentic_development + automated_pr_review

  **What to do**:
  - Implement 3 checkers in `internal/checker/build/`:
  - **vcs_cli_tools** (Level 5): Check for VCS CLI tools
    - Check for: `gh` CLI usage in docs/scripts, or `.github/` directory structure
    - Evidence: "GitHub CLI referenced in development workflow"
    - Suggestion: "Install GitHub CLI: `brew install gh` and document its usage"
  - **agentic_development** (Level 3): Check for AI agent support
    - Check for: `AGENTS.md`, `CLAUDE.md`, `.cursor/rules`, `.github/copilot-instructions.md`, `.claude/` directory
    - Evidence: "Found CLAUDE.md (886 lines) — comprehensive AI agent documentation"
    - Suggestion: "Create AGENTS.md documenting: build commands, test commands, architecture overview, coding conventions"
  - **automated_pr_review** (Level 4): Check for automated PR review tools
    - Check for: bot configs (`.github/bots/`, renovate config), CODEOWNERS, required reviewers in branch protection
    - Also check CI for: code review actions, linting on PRs
    - Evidence: "Found CODEOWNERS + CI lint-on-PR — automated review configured"
    - Suggestion: "Set up automated PR review. Add CODEOWNERS and CI checks that run on pull requests"
  - Write TDD tests with MapFS fixtures

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:
  - Factory.ai CockroachDB: `vcs_cli_tools` 1/1 — "gh CLI v2.52.0 installed and authenticated"
  - Factory.ai CockroachDB: `agentic_development` 1/1 — "Found CLAUDE.md with commit-helper skill"
  - Factory.ai CockroachDB: `automated_pr_review` 1/1 — "blathers-crl bot generates automated policy reviews"

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests
  - [ ] agentic_development detects AGENTS.md, CLAUDE.md, and .cursor/rules

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: agentic_development detects CLAUDE.md
    Tool: Bash (go test)
    Preconditions: MapFS with CLAUDE.md file
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestAgenticDevelopment`
      2. Verify Result.Passed == true
    Expected Result: Pass — agent docs found
    Evidence: .sisyphus/evidence/task-13-agentic.txt

  Scenario: agentic_development fails when no agent docs
    Tool: Bash (go test)
    Preconditions: MapFS with only go.mod
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestAgenticDevMissing`
      2. Verify Result.Passed == false, suggestion mentions AGENTS.md
    Expected Result: Fail — no agent documentation
    Evidence: .sisyphus/evidence/task-13-agentic-missing.txt
  ```

  **Commit**: YES (groups with T11, T12, T14)
  - Message: `feat(build): add vcs_cli_tools, agentic_development, automated_pr_review checkers`
  - Files: `internal/checker/build/vcs_cli_tools.go`, `internal/checker/build/agentic_development.go`, `internal/checker/build/automated_pr_review.go`, `internal/checker/build/*_test.go`
  - Pre-commit: `go test ./internal/checker/build/...`

- [ ] 14. Build System — build_performance_tracking + feature_flag_infra + release_notes_automation + unused_dependencies

  **What to do**:
  - Implement 4 checkers in `internal/checker/build/`:
  - **build_performance_tracking** (Level 4): Check for build timing/caching
    - Check for: build cache config (Bazel, Turborepo, Nx), CI timing artifacts
    - Evidence: "Found Turborepo config — build caching enabled"
  - **feature_flag_infrastructure** (Level 4): Check for feature flag system
    - Check for: feature flag library in deps (LaunchDarkly, Unleash, OpenFeature), custom flag config files
    - Go: Check go.mod for feature flag libraries
    - TypeScript: Check package.json for feature flag deps
    - Evidence: "Found @openfeature/js-sdk in dependencies"
  - **release_notes_automation** (Level 4): Check for automated release notes
    - Check for: `CHANGELOG.md`, release-please config, conventional-commits config, `cliff.toml` (git-cliff)
    - Evidence: "Found CHANGELOG.md + conventional-commits config"
  - **unused_dependencies_detection** (Level 5): Check for dependency audit tools
    - Go: `go mod tidy` in CI, or depcheck tools
    - TypeScript: `depcheck`, `knip` in package.json scripts
    - Java: `dependency:analyze` in Maven config
    - Evidence: "Found depcheck in package.json scripts"
  - Write TDD tests with MapFS fixtures

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **Acceptance Criteria**:
  - [ ] 4 checkers implemented with tests per language
  - [ ] Each checker correctly parses relevant config files

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: feature_flag_infra detects flag library in Go deps
    Tool: Bash (go test)
    Preconditions: MapFS with go.mod containing openfeature dependency
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestFeatureFlagGo`
      2. Verify Result.Passed == true
    Expected Result: Pass — feature flag library detected
    Evidence: .sisyphus/evidence/task-14-feature-flag.txt

  Scenario: All 4 checkers fail for bare-minimum repo
    Tool: Bash (go test)
    Preconditions: MapFS with only go.mod + main.go
    Steps:
      1. Run `go test ./internal/checker/build/ -v -run TestBuildAdvancedMissing`
      2. Verify all 4 return Passed == false
    Expected Result: All fail with suggestions
    Evidence: .sisyphus/evidence/task-14-all-missing.txt
  ```

  **Commit**: YES (groups with T11, T12, T13)
  - Message: `feat(build): add build_perf, feature_flags, release_notes, unused_deps checkers`
  - Files: `internal/checker/build/build_performance_tracking.go`, `internal/checker/build/feature_flag_infrastructure.go`, `internal/checker/build/release_notes_automation.go`, `internal/checker/build/unused_dependencies_detection.go`, `internal/checker/build/*_test.go`
  - Pre-commit: `go test ./internal/checker/build/...`

- [ ] 15. Testing — unit_tests_exist + unit_tests_runnable + test_naming_conventions + test_isolation

  **What to do**:
  - Implement 4 checkers in `internal/checker/testing/`:
  - **unit_tests_exist** (Level 1): Check for test file existence
    - Go: `*_test.go` files exist
    - TypeScript: `*.test.ts`, `*.spec.ts`, `__tests__/` directory
    - Java: `*Test.java`, `*Spec.java`, `src/test/` directory
    - Count test files, pass if > 0
    - Evidence: "Found 45 test files (*_test.go) across 12 packages"
  - **unit_tests_runnable** (Level 2): Check if test command is documented
    - Check README/CONTRIBUTING for test commands (`go test`, `npm test`, `mvn test`)
    - Also check: `Makefile` with `test` target, `package.json` with `test` script
    - Evidence: "Test command documented in README: `go test ./...`"
  - **test_naming_conventions** (Level 2): Check test file naming
    - Go: All test files follow `*_test.go` convention
    - TypeScript: Files in `__tests__/` or named `*.test.ts`/`*.spec.ts`
    - Java: Files follow `*Test.java` or `*Spec.java`
    - Evidence: "All 45 test files follow Go *_test.go convention"
  - **test_isolation** (Level 4): Check for parallel test support
    - Go: Check for `t.Parallel()` usage in test files, or `-race` flag in CI
    - TypeScript: Check Jest config for `maxWorkers`, or `--parallel` flag
    - Java: Check for `@Execution(CONCURRENT)` or parallel config in surefire
    - Evidence: "Found t.Parallel() in 30% of test functions — test isolation enabled"
  - Write TDD tests with MapFS fixtures

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:
  - Factory.ai CockroachDB: `unit_tests_exist` 1/1 — "2924 *_test.go files found"
  - Factory.ai CockroachDB: `unit_tests_runnable` 1/1 — "./dev test pkg/[package] command documented"
  - Factory.ai CockroachDB: `test_isolation` 1/1 — "Go tests support -race flag, parallel test execution"

  **Acceptance Criteria**:
  - [ ] 4 checkers implemented with tests per language
  - [ ] unit_tests_exist counts test files correctly
  - [ ] test_naming_conventions validates naming per language convention

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: unit_tests_exist counts Go test files
    Tool: Bash (go test)
    Preconditions: MapFS with 5 *_test.go files
    Steps:
      1. Run `go test ./internal/checker/testing/ -v -run TestUnitTestsExistGo`
      2. Verify Result.Passed == true, evidence mentions file count
    Expected Result: Pass — 5 test files found
    Evidence: .sisyphus/evidence/task-15-unit-tests.txt

  Scenario: unit_tests_exist fails for repo with no tests
    Tool: Bash (go test)
    Preconditions: MapFS with only main.go, no *_test.go files
    Steps:
      1. Run `go test ./internal/checker/testing/ -v -run TestUnitTestsExistNone`
      2. Verify Result.Passed == false
    Expected Result: Fail — no test files found
    Evidence: .sisyphus/evidence/task-15-no-tests.txt
  ```

  **Commit**: YES (groups with T16, T17)
  - Message: `feat(testing): add unit_tests, test_naming, test_isolation checkers`
  - Files: `internal/checker/testing/unit_tests_exist.go`, `internal/checker/testing/unit_tests_runnable.go`, `internal/checker/testing/test_naming_conventions.go`, `internal/checker/testing/test_isolation.go`, `internal/checker/testing/*_test.go`
  - Pre-commit: `go test ./internal/checker/testing/...`

- [ ] 16. Testing — integration_tests_exist + test_coverage_thresholds

  **What to do**:
  - Implement 2 checkers in `internal/checker/testing/`:
  - **integration_tests_exist** (Level 3): Check for integration/e2e tests
    - Check for: `integration/`, `e2e/`, `acceptance/` directories, test files with "integration" in name
    - Go: Check for `TestIntegration` or `TestAcceptance` function names, build tags `//go:build integration`
    - TypeScript: Check for Playwright/Cypress config, `e2e/` directory
    - Java: Check for `IT` suffix tests, `src/integrationTest/`
    - Evidence: "Found integration tests in pkg/acceptance/ directory"
  - **test_coverage_thresholds** (Level 3): Check for coverage configuration
    - Check for: coverage config in CI (`--coverprofile`, `--coverage`, `coverage` script in package.json)
    - Go: `go test -coverprofile` in CI config or Makefile
    - TypeScript: `jest --coverage` or `c8` config, `coverageThreshold` in jest config
    - Java: JaCoCo plugin in build config
    - Evidence: "Found coverage workflow in .github/workflows/ — coverage tracking enabled"
  - Write TDD tests

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **Acceptance Criteria**:
  - [ ] 2 checkers implemented with tests per language
  - [ ] Integration test detection works across directory patterns and naming conventions

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: integration_tests detects e2e directory
    Tool: Bash (go test)
    Preconditions: MapFS with e2e/ directory containing test files
    Steps:
      1. Run `go test ./internal/checker/testing/ -v -run TestIntegrationTestsExist`
      2. Verify Result.Passed == true
    Expected Result: Pass — integration tests found
    Evidence: .sisyphus/evidence/task-16-integration.txt

  Scenario: test_coverage detects Jest coverage config
    Tool: Bash (go test)
    Preconditions: MapFS with package.json containing jest.coverageThreshold
    Steps:
      1. Run `go test ./internal/checker/testing/ -v -run TestCoverageThresholdTS`
      2. Verify Result.Passed == true
    Expected Result: Pass — coverage threshold configured
    Evidence: .sisyphus/evidence/task-16-coverage.txt
  ```

  **Commit**: YES (groups with T15, T17)
  - Message: `feat(testing): add integration_tests_exist, test_coverage_thresholds checkers`
  - Files: `internal/checker/testing/integration_tests_exist.go`, `internal/checker/testing/test_coverage_thresholds.go`, `internal/checker/testing/*_test.go`
  - Pre-commit: `go test ./internal/checker/testing/...`

- [ ] 17. Testing — flaky_test_detection + test_performance_tracking

  **What to do**:
  - Implement 2 checkers in `internal/checker/testing/`:
  - **flaky_test_detection** (Level 4): Check for flaky test handling
    - Check for: retry/rerun mechanisms in CI, `--stress` flag usage, flaky test quarantine config
    - Go: `--count` flag for stress testing, `--flaky` labels in test framework
    - TypeScript: `jest --retry`, `--retries` in Playwright config
    - Java: `rerunFailingTestsCount` in surefire
    - Evidence: "Found test retry configuration in CI workflow"
  - **test_performance_tracking** (Level 4): Check for test timing tracking
    - Check for: `--durations` flag in CI, test timing artifacts, benchmark CI workflows
    - Go: `go test -bench` in CI, benchmark workflow files
    - TypeScript: `jest --verbose` with timing, performance budget config
    - Evidence: "Found benchmark CI workflow tracking test performance"
  - Write TDD tests

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **Acceptance Criteria**:
  - [ ] 2 checkers implemented with tests
  - [ ] Flaky test detection checks CI config for retry mechanisms

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: flaky_test_detection finds retry in CI
    Tool: Bash (go test)
    Preconditions: MapFS with CI workflow containing retry/rerun config
    Steps:
      1. Run `go test ./internal/checker/testing/ -v -run TestFlakyTestDetection`
      2. Verify Result.Passed == true
    Expected Result: Pass — retry mechanism found
    Evidence: .sisyphus/evidence/task-17-flaky.txt
  ```

  **Commit**: YES (groups with T15, T16)
  - Message: `feat(testing): add flaky_test_detection, test_performance_tracking checkers`
  - Files: `internal/checker/testing/flaky_test_detection.go`, `internal/checker/testing/test_performance_tracking.go`, `internal/checker/testing/*_test.go`
  - Pre-commit: `go test ./internal/checker/testing/...`

- [ ] 18. Documentation — readme + agents_md + documentation_freshness + skills

  **What to do**:
  - Implement 4 checkers in `internal/checker/docs/`:
  - **readme** (Level 1): Check README exists with meaningful content
    - Check for: `README.md`, `README`, `README.rst`, `README.txt`
    - Validate content: not empty, has at least 50 characters, has headings
    - Evidence: "Found README.md (2,450 bytes) with 5 sections"
    - Suggestion: "Create a README.md with: project description, installation, usage, contributing guide"
  - **agents_md** (Level 2): Check for AI agent documentation
    - Check for: `AGENTS.md`, `CLAUDE.md`, `.cursor/rules`, `.github/copilot-instructions.md`
    - If found, check minimum content (>100 chars, has code commands)
    - Evidence: "Found CLAUDE.md (886 lines) documenting build commands, test commands, and architecture"
    - Suggestion: "Create AGENTS.md documenting: build commands, test commands, architecture overview, coding conventions for AI agents"
  - **documentation_freshness** (Level 2): Check docs are up-to-date
    - Use `git log` to check last modification date of docs files (README.md, CONTRIBUTING.md, AGENTS.md)
    - Pass if: any doc file modified within last 180 days
    - Skip if: not a git repo
    - Evidence: "README.md last updated 23 days ago — documentation is fresh"
    - Suggestion: "Update documentation — last update was over 180 days ago"
  - **skills** (Level 4): Check for AI skill files
    - Check for: `.claude/skills/` directory, `.cursor/` directory with custom rules
    - Count skill files if directory exists
    - Evidence: "Found 2 skills in .claude/skills/"
    - Suggestion: "Create AI skills in .claude/skills/ to teach agents project-specific workflows"
  - Write TDD tests

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **References**:
  - Factory.ai CockroachDB: `readme` 1/1 — "README.md exists with comprehensive setup instructions"
  - Factory.ai CockroachDB: `agents_md` 1/1 — "CLAUDE.md exists at repo root (886 lines)"
  - Factory.ai CockroachDB: `documentation_freshness` 1/1 — "all modified within last 180 days"
  - Factory.ai CockroachDB: `skills` 1/1 — "2 skills configured in .claude/skills/"

  **Acceptance Criteria**:
  - [ ] 4 checkers implemented with tests
  - [ ] readme validates content quality (not just existence)
  - [ ] documentation_freshness uses git for date checking

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: readme passes for well-structured README
    Tool: Bash (go test)
    Preconditions: MapFS with README.md containing headings and 500+ chars
    Steps:
      1. Run `go test ./internal/checker/docs/ -v -run TestReadmeGood`
      2. Verify Result.Passed == true, evidence mentions sections
    Expected Result: Pass — README is substantive
    Evidence: .sisyphus/evidence/task-18-readme.txt

  Scenario: agents_md detects CLAUDE.md
    Tool: Bash (go test)
    Preconditions: MapFS with CLAUDE.md containing build commands
    Steps:
      1. Run `go test ./internal/checker/docs/ -v -run TestAgentsMd`
      2. Verify Result.Passed == true
    Expected Result: Pass — CLAUDE.md found
    Evidence: .sisyphus/evidence/task-18-agents-md.txt
  ```

  **Commit**: YES (groups with T19)
  - Message: `feat(docs): add readme, agents_md, documentation_freshness, skills checkers`
  - Files: `internal/checker/docs/readme.go`, `internal/checker/docs/agents_md.go`, `internal/checker/docs/documentation_freshness.go`, `internal/checker/docs/skills.go`, `internal/checker/docs/*_test.go`
  - Pre-commit: `go test ./internal/checker/docs/...`

- [ ] 19. Documentation — automated_doc_generation + service_flow_documented + api_schema_docs

  **What to do**:
  - Implement 3 checkers in `internal/checker/docs/`:
  - **automated_doc_generation** (Level 4): Check for doc generation tools
    - Go: `go generate` with doc tools, `godoc`, `pkgsite`
    - TypeScript: `typedoc`, `storybook`, `docusaurus` in devDeps
    - Java: `javadoc` plugin in build config, `dokka` for Kotlin
    - Also check: `docs/` directory with generation scripts
    - Evidence: "Found typedoc in devDependencies — automated doc generation configured"
  - **service_flow_documented** (Level 3): Check for architecture documentation
    - Check for: architecture diagrams (`.puml`, `.mmd`, `.drawio` files), `docs/architecture/`, `docs/design/`, `docs/rfcs/`
    - Also check: `ADR` (Architecture Decision Records) directory
    - Evidence: "Found 5 PlantUML diagrams in docs/tech-notes/"
    - Suggestion: "Document service flows. Create architecture diagrams using PlantUML or Mermaid"
  - **api_schema_docs** (Level 3): Check for API documentation
    - Check for: OpenAPI/Swagger files (`openapi.yaml`, `swagger.json`), GraphQL schema (`schema.graphql`)
    - Also check: `docs/api/` directory, Postman collections
    - Skip if: project doesn't appear to be an API (no HTTP handler/router imports)
    - Evidence: "Found openapi.yaml — API schema documented"
    - Suggestion: "Document your API with OpenAPI/Swagger specification"
  - Write TDD tests

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3
  - **Blocks**: T25, T26, T27
  - **Blocked By**: T2, T3

  **Acceptance Criteria**:
  - [ ] 3 checkers implemented with tests
  - [ ] api_schema_docs correctly skips for non-API projects

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: service_flow finds architecture diagrams
    Tool: Bash (go test)
    Preconditions: MapFS with docs/architecture/system.puml
    Steps:
      1. Run `go test ./internal/checker/docs/ -v -run TestServiceFlowDocumented`
      2. Verify Result.Passed == true
    Expected Result: Pass — architecture docs found
    Evidence: .sisyphus/evidence/task-19-service-flow.txt

  Scenario: api_schema_docs skips for CLI tool
    Tool: Bash (go test)
    Preconditions: MapFS with Go CLI project (no HTTP handlers)
    Steps:
      1. Run `go test ./internal/checker/docs/ -v -run TestApiSchemaDocsCLI`
      2. Verify Result.Skipped == true
    Expected Result: Skipped — not an API project
    Evidence: .sisyphus/evidence/task-19-api-skip.txt
  ```

  **Commit**: YES (groups with T18)
  - Message: `feat(docs): add automated_doc_gen, service_flow, api_schema_docs checkers`
  - Files: `internal/checker/docs/automated_doc_generation.go`, `internal/checker/docs/service_flow_documented.go`, `internal/checker/docs/api_schema_docs.go`, `internal/checker/docs/*_test.go`
  - Pre-commit: `go test ./internal/checker/docs/...`

- [ ] 20. HTML Reporter with Embedded Template

  **What to do**:
  - Implement `HTMLReporter` in `internal/reporter/html.go`:
    - Uses `html/template` for rendering
    - Template embedded via `//go:embed templates/report.html`
    - Outputs single self-contained HTML file (all CSS inline, no external deps)
    - Layout sections:
      1. Header: repo name, language, scan date, ari version
      2. Score card: Level badge (L1-L5), pass rate %, progress bar
      3. Level progression: L1-L5 bars showing pass rate per level
      4. Pillar breakdown: 4 cards showing per-pillar scores
      5. Criteria table: all 40 criteria with pass/fail/skip status, evidence, level
      6. Suggestions: prioritized list of failing criteria with fix instructions
      7. Footer: "Generated by ari" + timestamp
    - CSS: Clean, professional, light theme, responsive, uses CSS variables
    - No JavaScript whatsoever — static HTML only
  - Create template in `internal/reporter/templates/report.html`
  - Write TDD tests in `internal/reporter/html_test.go`:
    - Golden file comparison with `-update` flag for regeneration
    - Test: Output is valid HTML (parseable by `golang.org/x/net/html`)
    - Test: Contains all pillar names
    - Test: Contains score and level
    - Test: No external resource references (no `http://` or `https://` in href/src)
    - Test: Handles 0 criteria gracefully
    - Test: Handles all criteria passed (Level 5)

  **Must NOT do**:
  - No JavaScript in HTML output
  - No external CSS/font/image references
  - No dark mode (keep simple for MVP)
  - No interactive charts

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Template design + embedded assets + golden file testing
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with T21-T25)
  - **Blocks**: T25 (CLI needs reporters)
  - **Blocked By**: T4 (needs Score type), T6 (needs Reporter interface)

  **References**:
  - Go embed docs: https://pkg.go.dev/embed
  - Go html/template: https://pkg.go.dev/html/template
  - Factory.ai report visual: Level badges (L1-L5), pillar cards, criteria table with pass/fail

  **Acceptance Criteria**:
  - [ ] HTML report generates as single self-contained file
  - [ ] Valid HTML (parseable)
  - [ ] No external references
  - [ ] Golden file test passes
  - [ ] All sections present (header, score, pillars, criteria, suggestions)

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: HTML report is self-contained
    Tool: Bash (go test)
    Preconditions: Mock report data
    Steps:
      1. Run `go test ./internal/reporter/ -v -run TestHTMLReporterSelfContained`
      2. Generate HTML to buffer
      3. Parse with html.Parse — no errors
      4. Search for "http://" and "https://" in href/src attributes
    Expected Result: Valid HTML, zero external references
    Failure Indicators: Parse error, external URLs found
    Evidence: .sisyphus/evidence/task-20-html-self-contained.txt

  Scenario: HTML report shows correct level
    Tool: Bash (go test)
    Preconditions: Report with Level 3 score
    Steps:
      1. Run `go test ./internal/reporter/ -v -run TestHTMLReporterLevel`
      2. Generate HTML, check for "L3" or "Level 3" text
      3. Check for "Standardized" level name
    Expected Result: HTML contains correct level display
    Evidence: .sisyphus/evidence/task-20-html-level.txt
  ```

  **Commit**: YES (groups with T21)
  - Message: `feat(reporter): add HTML reporter with embedded template`
  - Files: `internal/reporter/html.go`, `internal/reporter/templates/report.html`, `internal/reporter/html_test.go`, `internal/reporter/testdata/golden.html`
  - Pre-commit: `go test ./internal/reporter/...`

- [ ] 21. Text Reporter (Non-TTY Fallback)

  **What to do**:
  - Implement `TextReporter` in `internal/reporter/text.go`:
    - Plain text output for non-TTY environments (piping, CI)
    - Format:
      ```
      ari — Agent Readiness Index
      ========================
      Repository: /path/to/repo
      Language: Go
      Level: 3 (Standardized)
      Pass Rate: 67% (27/40)
      
      Style & Validation: 9/12 (75%)
        ✓ lint_config — Found .golangci.yml
        ✓ formatter — gofmt built-in
        ✗ pre_commit_hooks — No pre-commit hooks configured
        ...
      
      Suggestions:
        1. [HIGH] pre_commit_hooks: Add pre-commit hooks...
        2. [MEDIUM] cyclomatic_complexity: Add complexity analysis...
      ```
    - Uses plain ASCII characters (✓, ✗) — no emoji, no ANSI colors
    - Suitable for `./ari --output text > report.txt`
  - Write TDD tests with golden file comparison

  **Must NOT do**:
  - No ANSI color codes — pure text
  - No box-drawing characters beyond basic ASCII

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Simple text formatting — straightforward implementation
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with T20, T22-T25)
  - **Blocks**: T25
  - **Blocked By**: T4, T6

  **Acceptance Criteria**:
  - [ ] Text output is clean, readable, pipe-friendly
  - [ ] No ANSI escape codes
  - [ ] Golden file test passes

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Text reporter produces pipe-friendly output
    Tool: Bash (go test)
    Preconditions: Mock report
    Steps:
      1. Run `go test ./internal/reporter/ -v -run TestTextReporterOutput`
      2. Verify output contains no ANSI codes (no \x1b[)
      3. Verify output is valid UTF-8 text
    Expected Result: Clean ASCII text output
    Evidence: .sisyphus/evidence/task-21-text-output.txt
  ```

  **Commit**: YES (groups with T20)
  - Message: `feat(reporter): add text reporter for non-TTY environments`
  - Files: `internal/reporter/text.go`, `internal/reporter/text_test.go`
  - Pre-commit: `go test ./internal/reporter/...`

- [ ] 22. TUI — Root Model + Progress View

  **What to do**:
  - Implement TUI root model in `internal/tui/model.go`:
    - Uses Bubbletea v2 (`charm.land/bubbletea/v2`)
    - `ViewID` enum: `ProgressView`, `ReportView`, `DetailView`
    - Root `Model` struct:
      - `currentView ViewID`
      - `progress *ProgressModel` — scanning progress
      - `report *ReportModel` — results display (nil until scan completes)
      - `detail *DetailModel` — criteria drill-down (nil until selected)
      - `err error` — error state
      - `quitting bool`
    - Root `View()` returns `tea.View` (v2 pattern):
      ```go
      func (m Model) View() tea.View {
          var content string
          switch m.currentView {
          case ProgressView: content = m.progress.View()
          case ReportView: content = m.report.View()
          case DetailView: content = m.detail.View()
          }
          return tea.NewView(content)
      }
      ```
    - Sub-models return `string` (not `tea.View`)
  - Implement `ProgressModel` in `internal/tui/views/progress.go`:
    - Shows: spinner + "Scanning repository..."
    - Shows: progress bar (N of M checkers completed)
    - Shows: current checker name being evaluated
    - Shows: scrolling log of completed checks with pass/fail
    - Uses `bubbles` components: `spinner`, `progress`
    - Receives `CheckerCompleteMsg` messages to update progress
  - Define message types in `internal/tui/messages.go`:
    - `ScanStartMsg`, `CheckerStartMsg{Name}`, `CheckerCompleteMsg{Result}`, `ScanCompleteMsg{Score, Report}`
    - `ErrorMsg{Error}`, `OpenBrowserMsg{Path}`, `QuitMsg`
  - Write TDD tests:
    - Test: Model starts in ProgressView
    - Test: View() returns tea.View (not string)
    - Test: ProgressModel updates on CheckerCompleteMsg
    - Test: Transitions to ReportView on ScanCompleteMsg

  **Must NOT do**:
  - Do not implement ReportView or DetailView here (T23)
  - Do not add fancy animations — spinner + progress bar only

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Bubbletea v2 architecture with proper Elm pattern, v2 is brand new
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with T20, T21, T23-T25)
  - **Blocks**: T23 (TUI views need root model), T25 (CLI needs TUI)
  - **Blocked By**: T4 (needs Score type for ScanCompleteMsg)

  **References**:
  - Bubbletea v2 API: `View()` returns `tea.View`, use `tea.NewView(content)`
  - Bubbles components: `github.com/charmbracelet/bubbles` (spinner, progress)
  - Lip Gloss: `github.com/charmbracelet/lipgloss` for styling

  **Acceptance Criteria**:
  - [ ] Root model compiles with Bubbletea v2
  - [ ] Progress view shows spinner + progress bar
  - [ ] State transitions work (Progress → Report)
  - [ ] View() returns tea.View

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: TUI starts in progress view
    Tool: Bash (go test)
    Preconditions: Model initialized
    Steps:
      1. Run `go test ./internal/tui/ -v -run TestModelInitialView`
      2. Create new Model, verify currentView == ProgressView
      3. Call View(), verify output contains spinner/progress text
    Expected Result: Model starts in ProgressView
    Evidence: .sisyphus/evidence/task-22-initial-view.txt

  Scenario: Progress view transitions on scan complete
    Tool: Bash (go test)
    Preconditions: Model in ProgressView
    Steps:
      1. Run `go test ./internal/tui/ -v -run TestProgressToReport`
      2. Send ScanCompleteMsg to Update()
      3. Verify currentView changed to ReportView
    Expected Result: Transition to ReportView
    Evidence: .sisyphus/evidence/task-22-transition.txt
  ```

  **Commit**: YES (groups with T23)
  - Message: `feat(tui): add Bubbletea v2 root model and progress view`
  - Files: `internal/tui/model.go`, `internal/tui/messages.go`, `internal/tui/views/progress.go`, `internal/tui/*_test.go`
  - Pre-commit: `go test ./internal/tui/...`

- [ ] 23. TUI — Report View + Detail View + Browser Opening

  **What to do**:
  - Implement `ReportModel` in `internal/tui/views/report.go`:
    - Shows overall level badge and pass rate
    - Lists 4 pillars with scores and pass/fail bars
    - Navigable list — arrow keys to select pillar
    - Press Enter to drill into pillar → switches to DetailView
    - Press `h` to generate and open HTML report
    - Press `j` to export JSON to stdout
    - Press `q` to quit
    - Uses `bubbles/list` or custom list component
    - Uses `lipgloss` for styling (colors, borders, padding)
  - Implement `DetailModel` in `internal/tui/views/detail.go`:
    - Shows all criteria for selected pillar
    - Each criterion: ✓/✗ status, name, evidence, level badge
    - Failed criteria show suggestion text
    - Press `Esc` or `Backspace` to return to ReportView
    - Scrollable with up/down keys
    - Uses `bubbles/viewport` for scrolling content
  - Implement browser opening in `internal/tui/browser.go`:
    - Uses `github.com/pkg/browser` for cross-platform browser opening
    - Generates HTML to temp file, opens in browser
    - Cleans up temp file on quit (or after 60s)
  - Write TDD tests:
    - Test: ReportModel displays all 4 pillar scores
    - Test: Arrow keys navigate pillar list
    - Test: Enter transitions to DetailView with correct pillar
    - Test: Esc returns from Detail to Report
    - Test: `q` triggers QuitMsg

  **Must NOT do**:
  - No mouse support in MVP
  - No window resizing handling beyond basic terminal width

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: UI design with lipgloss styling, list/viewport components, visual hierarchy
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with T20-T22, T24, T25)
  - **Blocks**: T25 (CLI needs all TUI views)
  - **Blocked By**: T4, T22 (needs root model + message types)

  **Acceptance Criteria**:
  - [ ] Report view shows all 4 pillar scores with visual bars
  - [ ] Navigation works (arrows, enter, esc, q)
  - [ ] Detail view shows per-criteria results with evidence
  - [ ] Browser opening works on macOS (at minimum)
  - [ ] Styling is clean and readable

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Report view shows pillar scores
    Tool: Bash (go test)
    Preconditions: ReportModel with mock score data
    Steps:
      1. Run `go test ./internal/tui/ -v -run TestReportViewPillars`
      2. Render View(), check output contains all 4 pillar names
      3. Check output contains pass rates
    Expected Result: All 4 pillars displayed with scores
    Evidence: .sisyphus/evidence/task-23-report-view.txt

  Scenario: Detail view shows criteria with evidence
    Tool: Bash (go test)
    Preconditions: DetailModel for Style & Validation pillar
    Steps:
      1. Run `go test ./internal/tui/ -v -run TestDetailViewCriteria`
      2. Render View(), check output contains criterion names
      3. Check pass/fail indicators present
    Expected Result: Criteria listed with status and evidence
    Evidence: .sisyphus/evidence/task-23-detail-view.txt
  ```

  **Commit**: YES (groups with T22)
  - Message: `feat(tui): add report view, detail view, and browser opening`
  - Files: `internal/tui/views/report.go`, `internal/tui/views/detail.go`, `internal/tui/browser.go`, `internal/tui/*_test.go`
  - Pre-commit: `go test ./internal/tui/...`

- [ ] 24. LLM Provider Implementations (OpenAI, Anthropic, Ollama)

  **What to do**:
  - Implement `OpenAIProvider` in `internal/llm/openai.go`:
    - Uses OpenAI Chat Completions API (`https://api.openai.com/v1/chat/completions`)
    - Default model: `gpt-4o-mini` (cost-effective for simple evaluations)
    - Sends structured JSON output request
    - Handles: rate limiting (429 → exponential backoff, 3 retries), auth errors (401)
    - Truncates input to 4000 tokens for large files
    - Uses `net/http` directly (no SDK dependency) to keep binary small
  - Implement `AnthropicProvider` in `internal/llm/anthropic.go`:
    - Uses Anthropic Messages API (`https://api.anthropic.com/v1/messages`)
    - Default model: `claude-sonnet-4-20250514`
    - Same error handling pattern as OpenAI
  - Implement `OllamaProvider` in `internal/llm/ollama.go`:
    - Uses Ollama API (`http://localhost:11434/api/chat`)
    - Default model: `llama3` (configurable via `ARI_LLM_MODEL`)
    - No auth needed (local)
    - Handles: connection refused (Ollama not running)
  - Factory function `NewProviderFromConfig` auto-selects provider based on `ARI_LLM_PROVIDER` env var
  - Write TDD tests:
    - All tests use `httptest.NewServer` to mock API responses
    - Test: OpenAI provider sends correct request format
    - Test: OpenAI provider handles 429 with retry
    - Test: Anthropic provider sends correct request format
    - Test: Ollama provider handles connection refused gracefully
    - Test: Factory function selects correct provider
    - NEVER call real APIs in tests

  **Must NOT do**:
  - Do not use SDK libraries (openai-go, anthropic-go) — use net/http directly
  - Do not implement streaming — simple request/response only
  - Do not call real APIs in tests

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: HTTP client implementation with error handling, retries, and mocking
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with T20-T23, T25)
  - **Blocks**: T26 (LLM wiring)
  - **Blocked By**: T5 (needs Provider interface)

  **References**:
  - OpenAI API: https://platform.openai.com/docs/api-reference/chat/create
  - Anthropic API: https://docs.anthropic.com/en/api/messages
  - Ollama API: https://github.com/ollama/ollama/blob/main/docs/api.md

  **Acceptance Criteria**:
  - [ ] 3 providers implemented, all satisfying Provider interface
  - [ ] All tests use httptest mock servers (no real API calls)
  - [ ] Rate limiting handled with exponential backoff
  - [ ] Factory function selects provider based on env var

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: OpenAI provider sends correct request
    Tool: Bash (go test)
    Preconditions: httptest server mocking OpenAI API
    Steps:
      1. Run `go test ./internal/llm/ -v -run TestOpenAIProvider`
      2. Verify request has correct headers (Authorization: Bearer)
      3. Verify request body has model, messages, temperature
      4. Verify response is parsed correctly
    Expected Result: Correct request/response flow
    Evidence: .sisyphus/evidence/task-24-openai.txt

  Scenario: Ollama provider handles connection refused
    Tool: Bash (go test)
    Preconditions: No Ollama server running (use unreachable host)
    Steps:
      1. Run `go test ./internal/llm/ -v -run TestOllamaConnectionRefused`
      2. Verify error is descriptive (not raw dial error)
      3. Verify FallbackEvaluator would catch this
    Expected Result: Clear error message about Ollama not running
    Evidence: .sisyphus/evidence/task-24-ollama-error.txt
  ```

  **Commit**: YES
  - Message: `feat(llm): add OpenAI, Anthropic, and Ollama provider implementations`
  - Files: `internal/llm/openai.go`, `internal/llm/anthropic.go`, `internal/llm/ollama.go`, `internal/llm/provider_test.go`
  - Pre-commit: `go test ./internal/llm/...`

- [ ] 25. CLI Entry Point (cmd/ari, Flag Parsing, Wiring)

  **What to do**:
  - Implement full CLI in `cmd/ari/main.go`:
    - Flag parsing:
      - `--path` (required): path to repository directory
      - `--output` (optional): `tui` (default), `json`, `html`, `text`
      - `--out` (optional): output file path (for html/json/text)
      - `--no-llm` (optional): skip LLM evaluation, use rule-based only
      - `--level-detail` (optional): show per-level breakdown in text output
      - `--version`: print ari version
      - `--help`: print usage
    - Wiring pipeline:
      1. Parse flags
      2. Validate `--path` exists and is a directory
      3. Create `os.DirFS(path)` for scanner
      4. Initialize LLM provider from env (or nil if `--no-llm`)
      5. Create checker registry with all 40 checkers
      6. Run scanner → detect language
      7. Run checker runner → collect results
      8. Run scorer → calculate level
      9. Build report
      10. Output based on `--output` flag:
          - `tui`: Launch Bubbletea program (default)
          - `json`: Write JSON to stdout or file
          - `html`: Generate HTML to file (default: `./ari-report.html`)
          - `text`: Write text to stdout or file
    - Exit codes: 0 = success, 1 = error, 2 = invalid args
    - Error handling: clear messages for common errors (path not found, permission denied)
  - Integrate ALL previously built modules:
    - `internal/scanner` → scan repo
    - `internal/checker` → run all 40 checkers
    - `internal/scorer` → calculate level
    - `internal/reporter` → output report
    - `internal/llm` → LLM evaluation (optional)
    - `internal/tui` → interactive display
  - Write integration test in `cmd/ari/main_test.go`:
    - Test: `--help` prints usage
    - Test: `--version` prints version
    - Test: `--path nonexistent` returns error
    - Test: `--output json --path testdata/sample` produces valid JSON

  **Must NOT do**:
  - Do not use cobra/viper — use `flag` stdlib for simplicity
  - Do not add subcommands — flat flag structure only

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Integration task wiring all modules together — needs comprehensive understanding of all components
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on nearly everything)
  - **Parallel Group**: Wave 4 (but last to complete — waits for all T7-T23)
  - **Blocks**: T27 (E2E tests)
  - **Blocked By**: T2-T4, T6, T7-T23 (needs all modules)

  **Acceptance Criteria**:
  - [ ] `go build ./cmd/ari` produces working binary
  - [ ] `./ari --help` prints usage with all flags
  - [ ] `./ari --path ./testdata/sample --output json` produces valid JSON
  - [ ] `./ari --path ./testdata/sample --output html --out /tmp/report.html` generates HTML
  - [ ] `./ari --path nonexistent` returns error with clear message
  - [ ] `./ari --no-llm --path ./testdata/sample --output json` works without LLM

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: CLI produces JSON output for sample repo
    Tool: Bash
    Preconditions: Binary built, testdata/sample-go-repo exists
    Steps:
      1. Run `go build ./cmd/ari`
      2. Run `./ari --path ./testdata/sample-go-repo --output json --no-llm`
      3. Parse output with `jq .level`
      4. Verify level is between 0 and 5
    Expected Result: Valid JSON with numeric level
    Failure Indicators: Parse error, level out of range
    Evidence: .sisyphus/evidence/task-25-cli-json.txt

  Scenario: CLI generates HTML report file
    Tool: Bash
    Preconditions: Binary built, testdata exists
    Steps:
      1. Run `./ari --path ./testdata/sample-go-repo --output html --out /tmp/ari-test.html --no-llm`
      2. Verify file exists: `test -f /tmp/ari-test.html`
      3. Verify file contains "<html"
    Expected Result: HTML file generated at specified path
    Failure Indicators: File not created, empty file, not HTML
    Evidence: .sisyphus/evidence/task-25-cli-html.txt

  Scenario: CLI shows error for invalid path
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./ari --path /nonexistent/path --output json 2>&1`
      2. Verify exit code is non-zero
      3. Verify stderr contains "not found" or "does not exist"
    Expected Result: Clear error message, non-zero exit
    Failure Indicators: Zero exit code, cryptic error
    Evidence: .sisyphus/evidence/task-25-cli-error.txt
  ```

  **Commit**: YES
  - Message: `feat(cmd): add CLI entry point with flag parsing and wiring`
  - Files: `cmd/ari/main.go`, `cmd/ari/main_test.go`
  - Pre-commit: `go build ./cmd/ari && go test ./cmd/ari/...`

- [ ] 26. Wire LLM into Applicable Criteria + FallbackEvaluator Integration

  **What to do**:
  - Identify all LLM-dependent criteria and wire in FallbackEvaluator:
    - `naming_consistency` (T8) — already has fallback, wire real LLM
    - `code_modularization` (T10) — already has fallback, wire real LLM
    - `documentation_freshness` (T18) — can be enhanced with LLM quality assessment
  - Create LLM prompt templates in `internal/llm/prompts.go`:
    - Standardized prompt format for each criterion
    - Require structured JSON response: `{"passed": bool, "evidence": string, "confidence": float}`
    - Include language-specific context in prompts
  - Wire `FallbackEvaluator` into checker runner:
    - If LLM provider configured → use LLM for applicable criteria
    - If no LLM → all criteria use rule-based (already implemented)
    - Runner tracks and reports mode per criterion
  - Add `--no-llm` flag support in runner (disable LLM even if configured)
  - Write integration tests:
    - Test: With LLM configured, applicable criteria use LLM mode
    - Test: With `--no-llm`, all criteria use rule-based mode
    - Test: LLM failure → fallback to rule-based gracefully
    - Test: LLM prompt templates produce valid prompts
    - All tests use MockProvider (no real API calls)

  **Must NOT do**:
  - Do not add new criteria — only wire LLM into existing ones
  - Do not implement prompt engineering beyond basic templates

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Integration across multiple modules with fallback logic
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with T27, T28)
  - **Blocks**: T27 (E2E tests)
  - **Blocked By**: T5, T7-T19, T24 (needs LLM interface, checkers, and providers)

  **Acceptance Criteria**:
  - [ ] LLM-dependent criteria use FallbackEvaluator
  - [ ] Mode is recorded per criterion ("llm" or "rule-based")
  - [ ] `--no-llm` flag forces all criteria to rule-based mode
  - [ ] LLM failure falls back gracefully

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: LLM mode recorded in criteria results
    Tool: Bash (go test)
    Preconditions: MockProvider configured as LLM
    Steps:
      1. Run `go test ./internal/checker/ -v -run TestLLMWiringMode`
      2. Run checkers with MockProvider
      3. Check naming_consistency result.Mode == "llm"
      4. Check lint_config result.Mode == "rule-based" (always rule-based)
    Expected Result: Correct mode per criterion
    Evidence: .sisyphus/evidence/task-26-llm-mode.txt

  Scenario: --no-llm forces all rule-based
    Tool: Bash (go test)
    Preconditions: MockProvider configured but --no-llm flag set
    Steps:
      1. Run `go test ./internal/checker/ -v -run TestNoLLMFlag`
      2. Run all checkers with no-llm option
      3. Verify ALL results have Mode == "rule-based"
    Expected Result: All criteria evaluated with rule-based mode
    Evidence: .sisyphus/evidence/task-26-no-llm.txt
  ```

  **Commit**: YES
  - Message: `feat(llm): wire LLM into applicable criteria with fallback evaluator`
  - Files: `internal/llm/prompts.go`, `internal/checker/runner.go` (update), `internal/llm/prompts_test.go`
  - Pre-commit: `go test ./...`

- [ ] 27. End-to-End Integration Tests + Testdata Fixtures

  **What to do**:
  - Create realistic test fixtures in `testdata/`:
    - `testdata/sample-go-repo/`: Minimal Go repo with go.mod, main.go, main_test.go, README.md, .golangci.yml
    - `testdata/sample-ts-repo/`: Minimal TS repo with package.json, tsconfig.json, src/index.ts, jest.config.js
    - `testdata/sample-java-repo/`: Minimal Java repo with pom.xml, src/main/java/App.java, README.md
    - `testdata/empty-repo/`: Empty directory
    - `testdata/no-git-repo/`: Repo without .git/ directory
    - `testdata/well-configured-repo/`: Repo with most criteria passing (for Level 3+ testing)
  - Write integration tests in `internal/integration_test.go`:
    - Test: Full pipeline on sample-go-repo → produces valid Score with level ≥ 1
    - Test: Full pipeline on empty-repo → Level 0, no crash
    - Test: Full pipeline on well-configured-repo → Level ≥ 3
    - Test: JSON output is valid and parseable
    - Test: HTML output is valid HTML
    - Test: All 40 criteria evaluated (or skipped with reason)
  - Write CLI integration test in `cmd/ari/integration_test.go`:
    - Test: Build binary, run against testdata, verify exit code + output
    - Test: `--output json` against sample-go-repo, parse JSON, verify structure
    - Test: `--output html` against sample-go-repo, verify HTML file created
    - Test: `--no-llm` works correctly
    - Use `go:build integration` tag for slow tests

  **Must NOT do**:
  - Do not include real large repos in testdata — keep fixtures minimal
  - Do not call real LLM APIs in integration tests

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Complex integration testing across all modules with realistic fixtures
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with T26, T28)
  - **Blocks**: F1-F4 (final verification)
  - **Blocked By**: T25, T26 (needs CLI and LLM wiring complete)

  **Acceptance Criteria**:
  - [ ] All testdata fixtures created with realistic file structures
  - [ ] Integration tests pass for all fixture types
  - [ ] Full pipeline works end-to-end (scan → check → score → report)
  - [ ] CLI integration test builds and runs binary successfully

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Full pipeline on Go repo produces valid report
    Tool: Bash
    Preconditions: All modules wired, testdata created
    Steps:
      1. Run `go build ./cmd/ari`
      2. Run `./ari --path ./testdata/sample-go-repo --output json --no-llm`
      3. Parse JSON output with jq
      4. Verify .level >= 1, .passRate > 0, .criteria | length == 40
    Expected Result: Valid JSON with 40 criteria, level ≥ 1
    Failure Indicators: Missing criteria, level 0 for valid repo, parse error
    Evidence: .sisyphus/evidence/task-27-e2e-go.txt

  Scenario: Empty repo returns Level 0 gracefully
    Tool: Bash
    Preconditions: testdata/empty-repo exists (empty directory)
    Steps:
      1. Run `./ari --path ./testdata/empty-repo --output json --no-llm`
      2. Parse JSON, verify .level == 0
      3. Verify no crash, exit code 0
    Expected Result: Level 0, all criteria failed or skipped
    Failure Indicators: Crash, non-zero exit, level > 0
    Evidence: .sisyphus/evidence/task-27-e2e-empty.txt
  ```

  **Commit**: YES
  - Message: `test(e2e): add end-to-end integration tests with testdata`
  - Files: `testdata/**/*`, `internal/integration_test.go`, `cmd/ari/integration_test.go`
  - Pre-commit: `go test ./...`

- [ ] 28. README + AGENTS.md + Goreleaser Setup

  **What to do**:
  - Create `README.md` for ari project:
    - Project description and motivation
    - Installation instructions (`go install`, binary download)
    - Quick start: `ari --path ./my-repo`
    - All CLI flags documented
    - Output format examples (JSON, HTML, text)
    - Maturity level descriptions (L1-L5)
    - Criteria reference table (40 criteria with pillar, level, description)
    - LLM configuration guide (env vars, provider setup)
    - Contributing guide section
  - Create `AGENTS.md` for ari itself (dogfooding):
    - Build command: `go build ./cmd/ari`
    - Test command: `go test ./...`
    - Architecture overview: cmd → internal/{scanner, checker, scorer, llm, reporter, tui}
    - Key patterns: fs.FS interface, Checker interface, FallbackEvaluator
    - Adding new checkers: step-by-step guide
    - Adding new pillars: step-by-step guide
  - Create `.goreleaser.yml` for binary distribution:
    - Build targets: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
    - Homebrew tap (optional, mark as TODO)
    - Checksum generation
    - Changelog from commits
  - Create `.golangci.yml` linter config for ari itself

  **Must NOT do**:
  - Do not create a website or landing page
  - Do not set up CI/CD (future task)

  **Recommended Agent Profile**:
  - **Category**: `writing`
    - Reason: Documentation-heavy task requiring clear technical writing
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with T26, T27)
  - **Blocks**: F1-F4 (final verification)
  - **Blocked By**: T25 (needs CLI flags documented)

  **Acceptance Criteria**:
  - [ ] README has all sections: install, quick start, flags, criteria table
  - [ ] AGENTS.md has build/test commands and architecture overview
  - [ ] .goreleaser.yml produces multi-platform builds
  - [ ] .golangci.yml configured for ari codebase

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: README covers all CLI flags
    Tool: Bash
    Preconditions: README.md written
    Steps:
      1. Run `./ari --help 2>&1` and capture all flag names
      2. Search README.md for each flag name
      3. Verify every flag is documented in README
    Expected Result: All flags documented
    Failure Indicators: Undocumented flags
    Evidence: .sisyphus/evidence/task-28-readme-flags.txt

  Scenario: Goreleaser builds for all targets
    Tool: Bash
    Preconditions: .goreleaser.yml configured
    Steps:
      1. Run `goreleaser check` (or `goreleaser build --snapshot --clean`)
      2. Verify no config errors
      3. Verify builds for darwin/amd64 and linux/amd64
    Expected Result: Valid goreleaser config, builds succeed
    Failure Indicators: Config errors, build failures
    Evidence: .sisyphus/evidence/task-28-goreleaser.txt
  ```

  **Commit**: YES
  - Message: `docs: add README, AGENTS.md, and goreleaser config`
  - Files: `README.md`, `AGENTS.md`, `.goreleaser.yml`, `.golangci.yml`
  - Pre-commit: `go build ./cmd/ari`

---

## Final Verification Wave (MANDATORY — after ALL implementation tasks)

> 4 review agents run in PARALLEL. ALL must APPROVE. Rejection → fix → re-run.

- [ ] F1. **Plan Compliance Audit** — `oracle`

  **What to do**:
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, run command). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/. Compare deliverables against plan.

  **Recommended Agent Profile**:
  - **Category**: `oracle` (via `subagent_type`)

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave FINAL (with F2, F3, F4)
  - **Blocks**: None
  - **Blocked By**: T27, T28

  **Acceptance Criteria**:
  - [ ] All "Must Have" items verified as present
  - [ ] All "Must NOT Have" items verified as absent
  - [ ] All 28 tasks have corresponding implementation

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: All Must Have items are implemented
    Tool: Bash
    Preconditions: All implementation tasks completed
    Steps:
      1. Run `go build ./cmd/ari` — verify binary builds
      2. Run `./ari --path ./testdata/sample-go-repo --output json --no-llm | jq '.criteria | length'` — verify 40 criteria
      3. Run `./ari --path ./testdata/sample-go-repo --output html --out /tmp/f1-test.html` — verify HTML output
      4. Run `ls internal/checker/style/ internal/checker/build/ internal/checker/testing/ internal/checker/docs/` — verify all pillar directories exist with checker files
      5. Run `go test ./... -count=1` — verify all tests pass
    Expected Result: Binary builds, 40 criteria in JSON, HTML generated, all pillar directories exist, tests pass
    Failure Indicators: Build fails, criteria count != 40, missing pillar directories, test failures
    Evidence: .sisyphus/evidence/f1-must-have-audit.txt

  Scenario: All Must NOT Have items are absent
    Tool: Bash
    Preconditions: Codebase complete
    Steps:
      1. Run `grep -r "github.com/google/go-github" . --include="*.go"` — verify no GitHub API client
      2. Run `grep -r "monorepo\|sub-application\|SubApp" . --include="*.go" -l` — verify no monorepo logic
      3. Run `grep -r "os\.Open\|os\.ReadFile\|os\.ReadDir" ./internal/checker/ --include="*.go"` — verify no direct os calls in checkers
      4. Run `grep -r "\.ari\.yaml\|\.ari\.yml\|\.ari\.json" . --include="*.go"` — verify no config file parsing
    Expected Result: All grep commands return empty (no matches)
    Failure Indicators: Any grep returns matches — forbidden pattern present
    Evidence: .sisyphus/evidence/f1-must-not-have-audit.txt
  ```

  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`

  **What to do**:
  Run `go vet ./...` + `golangci-lint run` + `go test ./...`. Review all files for: `//nolint` without reason, empty error handling, `fmt.Print` in production code, unused imports. Check AI slop: excessive comments, over-abstraction, generic names (data/result/item/temp). Verify `fs.FS` usage — no direct `os` calls in checker code.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave FINAL (with F1, F3, F4)
  - **Blocks**: None
  - **Blocked By**: T27, T28

  **Acceptance Criteria**:
  - [ ] `go vet ./...` passes with zero warnings
  - [ ] `go test ./...` passes with zero failures
  - [ ] No `//nolint` without justification comment
  - [ ] No `fmt.Print` in non-test, non-main production code

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Code passes static analysis
    Tool: Bash
    Preconditions: All code written
    Steps:
      1. Run `go vet ./... 2>&1` — capture output
      2. Run `go test ./... -count=1 2>&1` — capture output
      3. Run `grep -rn "//nolint" . --include="*.go" | grep -v "//nolint:" | grep -v "_test.go"` — find unjustified nolint
      4. Run `grep -rn 'fmt\.Print' ./internal/ --include="*.go" | grep -v "_test.go"` — find fmt.Print in prod code
    Expected Result: vet clean, all tests pass, no unjustified nolint, no fmt.Print in prod
    Failure Indicators: Vet warnings, test failures, unjustified nolint found, fmt.Print in prod
    Evidence: .sisyphus/evidence/f2-code-quality.txt

  Scenario: All checkers use fs.FS exclusively
    Tool: Bash
    Preconditions: All checker code written
    Steps:
      1. Run `grep -rn "os\.Open\|os\.Stat\|os\.ReadFile\|os\.ReadDir\|ioutil\." ./internal/checker/ --include="*.go" | grep -v "_test.go"`
      2. Verify zero matches
    Expected Result: No direct os filesystem calls in checker code
    Failure Indicators: Any match — checker bypasses fs.FS
    Evidence: .sisyphus/evidence/f2-fs-compliance.txt
  ```

  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Tests [N pass/N fail] | Files [N clean/N issues] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high`

  **What to do**:
  Start from clean state. Build `go build ./cmd/ari`. Execute against testdata repos. Test all output modes and edge cases.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave FINAL (with F1, F2, F4)
  - **Blocks**: None
  - **Blocked By**: T27, T28

  **Acceptance Criteria**:
  - [ ] TUI launches and displays correctly
  - [ ] All 3 output formats work (JSON, HTML, text)
  - [ ] Edge cases handled gracefully (empty repo, no git, invalid path)
  - [ ] `--no-llm` works correctly

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: JSON output is valid and complete
    Tool: Bash
    Preconditions: Binary built, testdata repos exist
    Steps:
      1. Run `go build ./cmd/ari`
      2. Run `./ari --path ./testdata/sample-go-repo --output json --no-llm > /tmp/f3-json.json`
      3. Run `jq '.level' /tmp/f3-json.json` — verify numeric level
      4. Run `jq '.criteria | length' /tmp/f3-json.json` — verify 40 criteria
      5. Run `jq '.criteria[] | select(.passed == false) | .suggestion' /tmp/f3-json.json | head -3` — verify suggestions exist for failures
    Expected Result: Valid JSON, level 1-5, 40 criteria, suggestions for failures
    Failure Indicators: jq parse error, wrong criteria count, missing suggestions
    Evidence: .sisyphus/evidence/f3-json-output.txt

  Scenario: HTML report generates and is self-contained
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./ari --path ./testdata/sample-go-repo --output html --out /tmp/f3-report.html --no-llm`
      2. Run `test -f /tmp/f3-report.html && echo "EXISTS" || echo "MISSING"`
      3. Run `grep -c "https\?://" /tmp/f3-report.html` — count external refs (should be 0 or very low)
      4. Run `wc -c /tmp/f3-report.html` — verify non-trivial size (>1000 bytes)
    Expected Result: HTML file exists, no/minimal external refs, substantial content
    Failure Indicators: File missing, external URLs found, tiny file
    Evidence: .sisyphus/evidence/f3-html-output.txt

  Scenario: Empty repo handled gracefully
    Tool: Bash
    Preconditions: testdata/empty-repo exists (empty directory)
    Steps:
      1. Run `./ari --path ./testdata/empty-repo --output json --no-llm 2>&1`
      2. Verify exit code 0 (no crash)
      3. Parse JSON, verify level == 0
    Expected Result: Level 0, graceful handling, no crash
    Failure Indicators: Non-zero exit, panic, crash
    Evidence: .sisyphus/evidence/f3-empty-repo.txt

  Scenario: Invalid path returns clear error
    Tool: Bash
    Preconditions: Binary built
    Steps:
      1. Run `./ari --path /nonexistent/totally/fake --output json 2>&1`
      2. Capture exit code with `echo $?`
      3. Verify exit code != 0
      4. Verify stderr contains descriptive error message
    Expected Result: Non-zero exit, clear error about path not found
    Failure Indicators: Zero exit code, cryptic error, panic
    Evidence: .sisyphus/evidence/f3-invalid-path.txt
  ```

  Output: `Scenarios [N/N pass] | Integration [N/N] | Edge Cases [N tested] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`

  **What to do**:
  For each task: read "What to do", read actual implementation. Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT do" compliance. Verify exactly 40 criteria, no more. Flag unaccounted changes.

  **Recommended Agent Profile**:
  - **Category**: `deep`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave FINAL (with F1, F2, F3)
  - **Blocks**: None
  - **Blocked By**: T27, T28

  **Acceptance Criteria**:
  - [ ] Exactly 40 criteria registered in default registry
  - [ ] No files modified outside of plan scope
  - [ ] All "Must NOT do" items verified as absent

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Exactly 40 criteria are registered
    Tool: Bash
    Preconditions: All checkers implemented and registered
    Steps:
      1. Run `grep -c "Register\|register" ./internal/checker/registry.go` or count NewDefaultRegistry entries
      2. Run `./ari --path ./testdata/sample-go-repo --output json --no-llm | jq '.criteria | length'`
      3. Verify count == 40
    Expected Result: Exactly 40 criteria — no more, no less
    Failure Indicators: Count != 40 — scope creep or missing criteria
    Evidence: .sisyphus/evidence/f4-criteria-count.txt

  Scenario: No forbidden patterns in codebase
    Tool: Bash
    Preconditions: All code written
    Steps:
      1. Run `grep -rn "monorepo\|subApplication\|SubApp" ./internal/ --include="*.go" | grep -v "_test.go" | grep -v "// not implemented"` — no monorepo logic
      2. Run `grep -rn "remediat\|autofix\|auto-fix\|AutoFix" ./internal/ --include="*.go" | grep -v "_test.go"` — no auto-remediation
      3. Run `grep -rn "github\.com/google/go-github\|octokit" . --include="*.go"` — no GitHub API client
      4. Count total .go files: `find . -name "*.go" -not -path "./.sisyphus/*" | wc -l` — verify reasonable count (not hundreds of unplanned files)
    Expected Result: No forbidden patterns, reasonable file count
    Failure Indicators: Forbidden patterns found, excessive unplanned files
    Evidence: .sisyphus/evidence/f4-scope-fidelity.txt
  ```

  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

| Wave | Commit | Message | Pre-commit |
|------|--------|---------|-----------|
| 1 | T1 | `feat(core): scaffold ari project with core types and interfaces` | `go build ./...` |
| 2 | T2 | `feat(scanner): add repo scanner with fs.FS and language detection` | `go test ./internal/scanner/...` |
| 2 | T3 | `feat(checker): add checker registry and runner engine` | `go test ./internal/checker/...` |
| 2 | T4 | `feat(scorer): add 5-level maturity scoring with gated progression` | `go test ./internal/scorer/...` |
| 2 | T5 | `feat(llm): add multi-provider LLM interface with mock and fallback` | `go test ./internal/llm/...` |
| 2 | T6 | `feat(reporter): add reporter interface and JSON reporter` | `go test ./internal/reporter/...` |
| 3 | T7-T10 | `feat(style): add style & validation checkers` | `go test ./internal/checker/style/...` |
| 3 | T11-T14 | `feat(build): add build system checkers` | `go test ./internal/checker/build/...` |
| 3 | T15-T17 | `feat(testing): add testing checkers` | `go test ./internal/checker/testing/...` |
| 3 | T18-T19 | `feat(docs): add documentation checkers` | `go test ./internal/checker/docs/...` |
| 4 | T20-T21 | `feat(reporter): add HTML and text reporters` | `go test ./internal/reporter/...` |
| 4 | T22-T23 | `feat(tui): add Bubbletea TUI with progress, report, and detail views` | `go test ./internal/tui/...` |
| 4 | T24 | `feat(llm): add OpenAI, Anthropic, and Ollama provider implementations` | `go test ./internal/llm/...` |
| 4 | T25 | `feat(cmd): add CLI entry point with flag parsing and wiring` | `go build ./cmd/ari && go test ./...` |
| 5 | T26 | `feat(llm): wire LLM into applicable criteria with fallback evaluator` | `go test ./...` |
| 5 | T27 | `test(e2e): add end-to-end integration tests with testdata` | `go test ./...` |
| 5 | T28 | `docs: add README, AGENTS.md, and goreleaser config` | `go build ./cmd/ari` |

---

## Success Criteria

### Verification Commands
```bash
go build ./cmd/ari                                    # Binary builds successfully
go test ./... -count=1                                # All tests pass
go vet ./...                                          # No vet warnings
./ari --help                                          # Shows usage
./ari --path ./testdata/sample-go-repo                # TUI launches
./ari --path ./testdata/sample-go-repo --output json  # Valid JSON output
./ari --path ./testdata/sample-go-repo --output html --out /tmp/report.html  # HTML generated
test -f /tmp/report.html                              # HTML file exists
```

### Final Checklist
- [ ] All 40 criteria implemented across 4 pillars
- [ ] All criteria have TDD tests (pass/fail/skip per language)
- [ ] 5-level scoring with 80% gated progression works
- [ ] TUI shows progress → report → detail views
- [ ] HTML report is self-contained (no external deps)
- [ ] JSON output includes all criteria results + level + evidence
- [ ] LLM integration works with OpenAI/Anthropic/Ollama
- [ ] Rule-based fallback works when no LLM configured
- [ ] Go, TypeScript, Java repos are all detected and evaluated
- [ ] No GitHub API calls anywhere in codebase
- [ ] No monorepo logic anywhere in codebase
- [ ] No auto-remediation logic anywhere in codebase
- [ ] All "Must NOT Have" guardrails verified
