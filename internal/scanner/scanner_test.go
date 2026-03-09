package scanner_test

import (
	"context"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/scanner"
)

func TestScanGoRepo(t *testing.T) {
	fsys := fstest.MapFS{
		"go.mod":        &fstest.MapFile{Data: []byte("module example.com/foo\ngo 1.26\n")},
		"main.go":       &fstest.MapFile{Data: []byte("package main\n")},
		"main_test.go":  &fstest.MapFile{Data: []byte("package main_test\n")},
		".git/HEAD":     &fstest.MapFile{Data: []byte("ref: refs/heads/main\n")},
		"pkg/helper.go": &fstest.MapFile{Data: []byte("package pkg\n")},
	}

	s := scanner.NewScanner()
	info, err := s.Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if info.Language != checker.LanguageGo {
		t.Fatalf("Language = %v, want %v", info.Language, checker.LanguageGo)
	}

	if !info.IsGitRepo {
		t.Fatal("IsGitRepo = false, want true")
	}

	if containsPath(info.Files, ".git/HEAD") {
		t.Fatal("expected .git/HEAD to be excluded")
	}

	if !containsPath(info.Files, "main.go") {
		t.Fatal("expected main.go to be included")
	}
}

func TestScanTSRepo(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json":   &fstest.MapFile{Data: []byte("{}")},
		"src/index.ts":   &fstest.MapFile{Data: []byte("export const x = 1\n")},
		"src/view.tsx":   &fstest.MapFile{Data: []byte("export const View = () => null\n")},
		"README.md":      &fstest.MapFile{Data: []byte("# app\n")},
		"node_modules/a": &fstest.MapFile{Data: []byte("ignore\n")},
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if info.Language != checker.LanguageTypeScript {
		t.Fatalf("Language = %v, want %v", info.Language, checker.LanguageTypeScript)
	}
}

func TestScanJavaRepo(t *testing.T) {
	fsys := fstest.MapFS{
		"pom.xml":                                 &fstest.MapFile{Data: []byte("<project/>\n")},
		"src/main/java/com/example/Main.java":     &fstest.MapFile{Data: []byte("class Main {}\n")},
		"src/test/java/com/example/MainTest.java": &fstest.MapFile{Data: []byte("class MainTest {}\n")},
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if info.Language != checker.LanguageJava {
		t.Fatalf("Language = %v, want %v", info.Language, checker.LanguageJava)
	}
}

func TestScanUnknownLanguage(t *testing.T) {
	fsys := fstest.MapFS{
		"README.md":      &fstest.MapFile{Data: []byte("# repo\n")},
		"docs/spec.adoc": &fstest.MapFile{Data: []byte("= Spec\n")},
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if info.Language != checker.LanguageUnknown {
		t.Fatalf("Language = %v, want %v", info.Language, checker.LanguageUnknown)
	}
}

func TestScanEmptyRepo(t *testing.T) {
	fsys := fstest.MapFS{}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(info.Files) != 0 {
		t.Fatalf("len(Files) = %d, want 0", len(info.Files))
	}

	if info.Language != checker.LanguageUnknown {
		t.Fatalf("Language = %v, want %v", info.Language, checker.LanguageUnknown)
	}

	if info.IsGitRepo {
		t.Fatal("IsGitRepo = true, want false")
	}
}

func TestScanIgnoresGitDir(t *testing.T) {
	fsys := fstest.MapFS{
		"go.mod":      &fstest.MapFile{Data: []byte("module example.com/foo\ngo 1.26\n")},
		".git/HEAD":   &fstest.MapFile{Data: []byte("ref: refs/heads/main\n")},
		".git/config": &fstest.MapFile{Data: []byte("[core]\n")},
		"cmd/main.go": &fstest.MapFile{Data: []byte("package main\n")},
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if containsPath(info.Files, ".git/HEAD") || containsPath(info.Files, ".git/config") {
		t.Fatal("expected .git entries to be excluded")
	}
}

func TestScanIgnoresNodeModules(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json":                &fstest.MapFile{Data: []byte("{}")},
		"src/index.ts":                &fstest.MapFile{Data: []byte("export {}\n")},
		"node_modules/lib/index.js":   &fstest.MapFile{Data: []byte("module.exports = {}\n")},
		"node_modules/.bin/some-tool": &fstest.MapFile{Data: []byte("#!/bin/sh\n")},
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if containsPath(info.Files, "node_modules/lib/index.js") {
		t.Fatal("expected node_modules files to be excluded")
	}

	if !containsPath(info.Files, "src/index.ts") {
		t.Fatal("expected src/index.ts to be included")
	}
}

func TestScanFileLimitRespected(t *testing.T) {
	fsys := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module example.com/foo\ngo 1.26\n")},
	}

	for i := 0; i < 6000; i++ {
		fsys[fmt.Sprintf("src/file%d.go", i)] = &fstest.MapFile{Data: []byte("package main\n")}
	}

	info, err := scanner.NewScanner().Scan(context.Background(), fsys)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(info.Files) > scanner.DefaultFileLimit {
		t.Fatalf("len(Files) = %d, want <= %d", len(info.Files), scanner.DefaultFileLimit)
	}
}

func TestDetectLanguagePriority(t *testing.T) {
	files := []scanner.FileInfo{
		{Path: "package.json", Extension: ".json"},
		{Path: "go.mod", Extension: ".mod"},
		{Path: "src/app.ts", Extension: ".ts"},
	}

	got := scanner.DetectLanguage(files)
	if got != checker.LanguageGo {
		t.Fatalf("DetectLanguage() = %v, want %v", got, checker.LanguageGo)
	}
}

func TestDetectLanguageFallback(t *testing.T) {
	files := []scanner.FileInfo{
		{Path: "pkg/a.go", Extension: ".go"},
		{Path: "pkg/b.go", Extension: ".go"},
		{Path: "web/view.ts", Extension: ".ts"},
	}

	got := scanner.DetectLanguage(files)
	if got != checker.LanguageGo {
		t.Fatalf("DetectLanguage() = %v, want %v", got, checker.LanguageGo)
	}
}

func containsPath(files []scanner.FileInfo, path string) bool {
	for _, f := range files {
		if f.Path == path {
			return true
		}
	}
	return false
}
