package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/all"
	"github.com/nixbpe/ari/internal/llm"
	"github.com/nixbpe/ari/internal/reporter"
	"github.com/nixbpe/ari/internal/scanner"
	"github.com/nixbpe/ari/internal/scorer"
	"github.com/nixbpe/ari/internal/tui"
)

const ariVersion = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ari", flag.ContinueOnError)
	fs.SetOutput(stderr)

	pathFlag := fs.String("path", "", "Path to repository to evaluate (required)")
	outputFlag := fs.String("output", "tui", "Output format: tui, json, html, text")
	outFlag := fs.String("out", "", "Output file path (for html/json/text)")
	noLLMFlag := fs.Bool("no-llm", false, "Skip LLM evaluation, use rule-based only")
	levelDetailFlag := fs.Bool("level-detail", false, "Show per-level breakdown in text output")
	showVersionFlag := fs.Bool("version", false, "Print ari version")
	showHelpFlag := fs.Bool("help", false, "Show this help message")

	// Suppress unused warning until level-detail is fully implemented
	_ = levelDetailFlag

	printUsage := func() {
		fmt.Fprintln(stdout, "ari - Agent Readiness Index")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Usage: ari [flags]")
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "Flags:")
		// Capture PrintDefaults output and redirect to stdout.
		var buf bytes.Buffer
		prev := fs.Output()
		fs.SetOutput(&buf)
		fs.PrintDefaults()
		fs.SetOutput(prev)
		fmt.Fprint(stdout, buf.String())
	}

	fs.Usage = printUsage

	if err := fs.Parse(args); err != nil {
		// flag.ErrHelp is returned when -h is passed (unregistered short flag).
		// fs.Usage was already called by the flag package.
		if err == flag.ErrHelp {
			return 0
		}
		// Other parse errors are already printed to stderr by the flag package.
		return 2
	}

	if *showHelpFlag {
		printUsage()
		return 0
	}

	if *showVersionFlag {
		fmt.Fprintf(stdout, "ari version %s\n", ariVersion)
		return 0
	}

	// No flags at all → show help
	if fs.NFlag() == 0 && len(fs.Args()) == 0 {
		printUsage()
		return 0
	}

	// Validate --path
	if *pathFlag == "" {
		fmt.Fprintf(stderr, "error: --path is required\n")
		return 2
	}

	fi, err := os.Stat(*pathFlag)
	if err != nil {
		fmt.Fprintf(stderr, "error: path %q does not exist\n", *pathFlag)
		return 2
	}

	if !fi.IsDir() {
		fmt.Fprintf(stderr, "error: path %q is not a directory\n", *pathFlag)
		return 2
	}

	ctx := context.Background()
	repoFS := os.DirFS(*pathFlag)

	// LLM evaluator (optional — nil if --no-llm or no env config)
	var eval llm.Evaluator
	if !*noLLMFlag {
		cfg := llm.ConfigFromEnv()
		if cfg.Provider != "" {
			p, provErr := llm.NewProviderFromConfig(cfg)
			if provErr != nil {
				fmt.Fprintf(stderr, "warning: LLM setup failed: %v (continuing without LLM)\n", provErr)
			} else {
				eval = &llm.FallbackEvaluator{
					Primary:  &llm.ProviderEvaluator{Provider: p},
					Fallback: &llm.RuleBasedEvaluator{},
				}
			}
		}
	}

	// Registry — all 40 checkers registered with optional LLM evaluator
	registry := checker.NewDefaultRegistry()
	all.RegisterAll(registry, eval)

	// Runner
	rnr := &checker.Runner{Registry: registry}

	// Scan repository
	sc := scanner.NewScanner()
	repoInfo, scanErr := sc.Scan(ctx, repoFS)
	if scanErr != nil {
		fmt.Fprintf(stderr, "error: scan failed: %v\n", scanErr)
		return 1
	}
	repoInfo.RootPath = *pathFlag

	// Run checkers
	results, runErr := rnr.Run(ctx, repoFS, repoInfo)
	if runErr != nil {
		fmt.Fprintf(stderr, "error: checker run failed: %v\n", runErr)
		return 1
	}

	// Score
	score := scorer.New().Calculate(results)

	// Build report
	report := reporter.BuildReport(repoInfo, score, results)

	return outputReport(ctx, *outputFlag, *outFlag, score, report, results, stdout, stderr)
}

// outputReport renders the report in the requested format.
// Exit codes: 0 = success, 1 = runtime error, 2 = bad format argument.
func outputReport(
	ctx context.Context,
	format, outPath string,
	score *scorer.Score,
	report *reporter.Report,
	results []*checker.Result,
	stdout, stderr io.Writer,
) int {
	// If --out was provided, open that file for writing.
	var w io.Writer
	var closeFile func() error

	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			fmt.Fprintf(stderr, "error: create output file %q: %v\n", outPath, err)
			return 1
		}
		w = f
		closeFile = f.Close
	}

	// Ensure file is closed after writing regardless of return path.
	defer func() {
		if closeFile != nil {
			if err := closeFile(); err != nil {
				// Nothing useful to do from a deferred close; log to stderr.
				fmt.Fprintf(stderr, "warning: close output file: %v\n", err)
			}
		}
	}()

	switch format {
	case "json":
		if w == nil {
			w = stdout
		}
		rep := &reporter.JSONReporter{}
		if err := rep.Report(ctx, report, w); err != nil {
			fmt.Fprintf(stderr, "error: JSON report: %v\n", err)
			return 1
		}

	case "html":
		if w == nil {
			// Default output file when --out is not specified.
			const defaultHTMLPath = "ari-report.html"
			f, err := os.Create(defaultHTMLPath)
			if err != nil {
				fmt.Fprintf(stderr, "error: create HTML output file %q: %v\n", defaultHTMLPath, err)
				return 1
			}
			w = f
			closeFile = f.Close
		}
		rep := &reporter.HTMLReporter{}
		if err := rep.Report(ctx, report, w); err != nil {
			fmt.Fprintf(stderr, "error: HTML report: %v\n", err)
			return 1
		}

	case "text":
		if w == nil {
			w = stdout
		}
		rep := &reporter.TextReporter{}
		if err := rep.Report(ctx, report, w); err != nil {
			fmt.Fprintf(stderr, "error: text report: %v\n", err)
			return 1
		}

	case "tui":
		// TUI mode: close any pre-opened file (not expected, but safe).
		if closeFile != nil {
			_ = closeFile()
			closeFile = nil
		}
		model := tui.NewModel()
		p := tea.NewProgram(model)
		// Scan is already complete; send the result immediately.
		go func() {
			p.Send(tui.ScanCompleteMsg{Score: score, Report: report, Results: results})
		}()
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(stderr, "error: TUI: %v\n", err)
			return 1
		}

	default:
		fmt.Fprintf(stderr, "error: unknown output format %q\n", format)
		return 2
	}

	return 0
}
