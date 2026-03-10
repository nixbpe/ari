# AGENTS.md — ari Development Guide

## Build & Test

```bash
# Build
go build ./cmd/ari

# Test all packages
go test ./...

# Run against current directory
./ari --path . --output json --no-llm
```

## Architecture

```
cmd/ari/main.go              — CLI entry point, flag parsing, pipeline wiring
internal/
  scanner/                   — Repository scanning, language detection
  checker/                   — Checker interface, registry, runner
    all/                     — RegisterAll() — registers all 72 checkers
    helpers.go               — Shared helpers: CIWorkflowContains, DepFileContains, FileExistsAny, FileContentContains
    style/                   — 12 Style & Validation checkers
    build/                   — 13 Build System checkers
    testing/                 — 8 Testing checkers
    docs/                    — 7 Documentation checkers
    devenv/                  — 7 Dev Environment checkers
    observability/           — 8 Debugging & Observability checkers
    security/                — 7 Security checkers
    taskdiscovery/           — 5 Task Discovery checkers
    analytics/               — 5 Product & Analytics checkers
  scorer/                    — 5-level maturity scoring with gated progression
  llm/                       — Multi-provider LLM interface with fallback
  reporter/                  — HTML, JSON, Text reporters
  tui/                       — Bubbletea v2 interactive TUI
```

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

### FallbackEvaluator
LLM evaluation uses a fallback chain: LLM → rule-based → never fails.

### Adding a New Checker
1. Create `internal/checker/{pillar}/{id}.go`
2. Implement the `Checker` interface
3. Register in `internal/checker/all/all.go`
4. Add tests in `{pillar}/{id}_test.go`

### Adding a New Pillar
1. Add constant to `internal/checker/checker.go`
2. Create `internal/checker/{pillar}/` directory
3. Implement checkers
4. Register in `internal/checker/all/all.go`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ARI_LLM_PROVIDER` | LLM provider: `openai`, `anthropic`, `ollama` |
| `ARI_API_KEY` | API key for OpenAI or Anthropic |
| `ARI_LLM_MODEL` | Override default model |
| `ARI_LLM_BASE_URL` | Override API base URL |
