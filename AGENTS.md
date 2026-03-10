# AGENTS.md — ari Development Guide

## Build & Test

```bash
# Build
go build ./cmd/ari

# Test all packages
go test ./...

# Lint
golangci-lint run ./...

# Run against current directory
./ari --path . --output json --no-llm
```

## Makefile Targets

| Target | Command |
|--------|---------|
| `make build` | `go build ./cmd/ari` |
| `make test` | `go test -race --count=1 -coverprofile=coverage.out -covermode=atomic ./...` |
| `make lint` | `golangci-lint run ./...` |
| `make coverage` | Test + `go tool cover -func=coverage.out` |
| `make bench` | `go test -bench=./... -benchmem -run=^$ ./...` |
| `make setup` | `go mod download` + install go-test-coverage |

## Architecture

```
cmd/ari/main.go              — CLI entry point, flag parsing, pipeline wiring
internal/
  scanner/                   — Repository scanning, language detection
  checker/                   — Checker interface, registry, runner
    all/                     — RegisterAll() — registers all 72 checkers
    helpers.go               — Shared helpers: CIWorkflowContains, DepFileContains, FileExistsAny, FileContentContains
    style/                   — 12 checkers → Constraints & Governance pillar
    build/                   — 13 checkers → split across all 4 pillars
    testing/                 — 8 checkers → Verification & Feedback pillar
    docs/                    — 7 checkers → Context & Intent pillar
    devenv/                  — 7 checkers → Environment & Infra pillar
    observability/           — 8 checkers → Verification & Feedback pillar
    security/                — 7 checkers → Constraints & Governance pillar
    taskdiscovery/           — 5 checkers → Context & Intent pillar
    analytics/               — 5 checkers → Context & Intent pillar
  scorer/                    — 5-level maturity scoring with gated progression (>=80% per level)
  llm/                       — Multi-provider LLM interface with fallback chain
  reporter/                  — HTML, JSON, Text reporters
  tui/                       — Bubbletea v2 interactive TUI (progress, report, detail views)
```

## Pillars (4 MECE categories)

Checkers are organized by source directory (style/, build/, etc.) but mapped to 4 evaluation pillars:

| Pillar | Description | Source packages |
|--------|-------------|-----------------|
| Context & Intent | Documentation, discovery, analytics | docs, taskdiscovery, analytics, build (partial) |
| Environment & Infra | Dev setup, build tooling, CI config | devenv, build (partial) |
| Constraints & Governance | Linting, security, code quality | style, security, build (partial) |
| Verification & Feedback | Testing, observability, release | testing, observability, build (partial) |

## Key Patterns

### fs.FS Interface
All checker code uses `fs.FS` for file access — never `os.*` directly.
This enables testing with `testing/fstest.MapFS`.

### Checker Interface
```go
type Checker interface {
    ID() CheckerID
    Pillar() Pillar
    Level() Level
    Name() string
    Description() string
    Check(ctx context.Context, repo fs.FS, lang Language) (*Result, error)
}
```

### SuggestionProvider
Checkers may optionally implement `SuggestionProvider` to return fix suggestions for failing criteria:
```go
type SuggestionProvider interface {
    Suggestion() string
}
```

### FallbackEvaluator
LLM evaluation uses a fallback chain: LLM provider -> rule-based evaluator -> never fails.
Configured via environment variables; nil evaluator means rule-based only.

### Runner
Sequential checker execution with `OnStart`/`OnDone` callbacks for progress reporting.
Handles panics gracefully. Skips checkers that don't support the detected language via `SupportsLanguage()`.

### Shared Helpers (`internal/checker/helpers.go`)
| Helper | Purpose |
|--------|---------|
| `CIWorkflowContains(repo, keywords)` | Scan `.github/workflows/*.yml` for keywords |
| `DepFileContains(repo, lang, packages)` | Check dependency files for package names |
| `FileExistsAny(repo, paths)` | Check if any candidate path exists |
| `FileContentContains(repo, path, keywords)` | Check file content for keywords |

## Adding a New Checker

1. Create `internal/checker/{package}/{id}.go`
2. Implement the `Checker` interface (set correct `Pillar()` return value)
3. Optionally implement `SuggestionProvider`
4. Register in `internal/checker/all/all.go`
5. Add tests in `{package}/{id}_test.go` using `testing/fstest.MapFS`

## Adding a New Pillar

1. Add constant to `internal/checker/checker.go` (in the `Pillar` iota block)
2. Add `String()` case
3. Update `runner.go` pillar range validation
4. Create checker package directory if needed
5. Register checkers in `internal/checker/all/all.go`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ARI_LLM_PROVIDER` | LLM provider: `anthropic`, `ollama` |
| `ARI_API_KEY` | API key for the provider |
| `ARI_LLM_MODEL` | Override default model |
| `ARI_LLM_BASE_URL` | Override API base URL |

## CI

GitHub Actions (`.github/workflows/ci.yml`): lint (golangci-lint v2.1.6) + test (race, coverage, benchmarks).
Release via GoReleaser (`.goreleaser.yml`): cross-compile for linux/darwin/windows (amd64/arm64).
Pre-commit hooks: `gofmt`, `go vet`, `golangci-lint`.
