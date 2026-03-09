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
