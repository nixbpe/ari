# ari — Agent Readiness Index

> Evaluate how ready your codebase is for AI coding agents.

## What is ari?

ari scans a local repository and evaluates it across 40 criteria in 4 pillars, 
assigning a maturity level from 1 (Functional) to 5 (Autonomous).

## Installation

### From source
```bash
go install github.com/nixbpe/ari/cmd/ari@latest
```

### Build locally
```bash
git clone https://github.com/nixbpe/ari
cd ari
go build ./cmd/ari
```

## Quick Start

```bash
# Evaluate current directory (interactive TUI)
ari --path .

# JSON output
ari --path . --output json

# HTML report
ari --path . --output html --out report.html

# Skip LLM evaluation
ari --path . --no-llm
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | (required) | Path to repository to evaluate |
| `--output` | `tui` | Output format: `tui`, `json`, `html`, `text` |
| `--out` | | Output file path (for html/json/text) |
| `--no-llm` | false | Skip LLM evaluation, use rule-based only |
| `--level-detail` | false | Show per-level breakdown in text output |
| `--version` | | Print ari version |
| `--help` | | Show usage |

## Maturity Levels

| Level | Name | Description |
|-------|------|-------------|
| 1 | Functional | Basic tooling in place |
| 2 | Documented | Code and processes documented |
| 3 | Standardized | Consistent standards enforced |
| 4 | Optimized | Performance and quality optimized |
| 5 | Autonomous | Fully ready for AI agents |

Progression is gated: you must achieve ≥80% at each level before advancing.

## Criteria Reference

### Style & Validation (12 criteria)
| ID | Level | Description |
|----|-------|-------------|
| lint_config | L1 | Linter configuration file present |
| formatter | L1 | Code formatter configured |
| type_check | L1 | Type checking enabled |
| strict_typing | L2 | Strict type checking enforced |
| pre_commit_hooks | L2 | Pre-commit hooks configured |
| naming_consistency | L2 | Naming conventions consistent |
| cyclomatic_complexity | L3 | Complexity limits configured |
| dead_code_detection | L3 | Dead code detection enabled |
| duplicate_code_detection | L4 | Duplicate code detection configured |
| code_modularization | L4 | Module boundaries enforced |
| large_file_detection | L3 | Large file detection configured |
| tech_debt_tracking | L3 | Tech debt tracking in place |

### Build System (13 criteria)
| ID | Level | Description |
|----|-------|-------------|
| build_cmd_doc | L1 | Build command documented |
| single_command_setup | L1 | Single command setup available |
| deps_pinned | L1 | Dependencies pinned/locked |
| fast_ci_feedback | L2 | CI feedback under 10 minutes |
| release_automation | L2 | Release process automated |
| deployment_frequency | L3 | Deployment frequency tracked |
| vcs_cli_tools | L2 | VCS CLI tools configured |
| agentic_development | L3 | AI agent development docs present |
| automated_pr_review | L3 | Automated PR review configured |
| build_performance_tracking | L4 | Build performance tracked |
| feature_flag_infra | L4 | Feature flag infrastructure present |
| release_notes_automation | L4 | Release notes automated |
| unused_dependencies | L4 | Unused dependency detection |

### Testing (8 criteria)
| ID | Level | Description |
|----|-------|-------------|
| unit_tests_exist | L1 | Unit tests present |
| unit_tests_runnable | L1 | Unit tests can be run |
| test_naming_conventions | L2 | Test naming conventions followed |
| test_isolation | L2 | Tests are isolated |
| integration_tests_exist | L3 | Integration tests present |
| test_coverage_thresholds | L3 | Coverage thresholds configured |
| flaky_test_detection | L4 | Flaky test detection in place |
| test_performance_tracking | L4 | Test performance tracked |

### Documentation (7 criteria)
| ID | Level | Description |
|----|-------|-------------|
| readme | L1 | README.md present |
| agents_md | L2 | AGENTS.md present |
| documentation_freshness | L2 | Documentation updated recently |
| skills | L3 | AI agent skills documented |
| automated_doc_generation | L3 | Documentation auto-generated |
| service_flow_documented | L4 | Service flow diagrams present |
| api_schema_docs | L4 | API schema documented |

## LLM Configuration

ari supports optional LLM evaluation for criteria that benefit from semantic analysis:

```bash
# OpenAI
export ARI_LLM_PROVIDER=openai
export ARI_API_KEY=sk-...
ari --path .

# Anthropic
export ARI_LLM_PROVIDER=anthropic
export ARI_API_KEY=sk-ant-...
ari --path .

# Ollama (local)
export ARI_LLM_PROVIDER=ollama
ari --path .

# Custom model
export ARI_LLM_MODEL=gpt-4o
ari --path .
```

Without LLM configuration, ari uses rule-based evaluation for all criteria.

## Supported Languages

- Go
- TypeScript / JavaScript
- Java / Kotlin

## Output Formats

### TUI (default)
Interactive terminal UI with progress view, report view, and drill-down detail view.

### JSON
```bash
ari --path . --output json | jq '.level'
```

### HTML
Self-contained HTML report with inline CSS. No external dependencies.

### Text
Plain text output suitable for CI logs and non-TTY environments.

## Contributing

See [AGENTS.md](AGENTS.md) for development setup and architecture overview.
