package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	showHelp := flag.Bool("help", false, "Show usage")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showHelp || flag.NFlag() == 0 && len(flag.Args()) == 0 {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("ari version %s\n", version)
		os.Exit(0)
	}

	printHelp()
}

func printHelp() {
	fmt.Println("ari - Agent Readiness Index")
	fmt.Println()
	fmt.Println("Usage: ari [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help        Show this help message")
	fmt.Println("  --version     Show version")
	fmt.Println("  --path        Path to repository to evaluate (required)")
	fmt.Println("  --output      Output format: tui (default), json, html, text")
	fmt.Println("  --out         Output file path (for html/json/text)")
	fmt.Println("  --no-llm      Skip LLM evaluation, use rule-based only")
}
