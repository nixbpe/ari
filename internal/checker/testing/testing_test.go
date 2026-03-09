package checktesting_test

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
	checktesting "github.com/bbik/ari/internal/checker/testing"
)

func TestUnitTestsExistGo(t *testing.T) {
	repo := fstest.MapFS{
		"main.go":                              &fstest.MapFile{Data: []byte("package main")},
		"main_test.go":                         &fstest.MapFile{Data: []byte("package main")},
		"internal/foo/foo_test.go":             &fstest.MapFile{Data: []byte("package foo")},
		"internal/bar/bar_test.go":             &fstest.MapFile{Data: []byte("package bar")},
		"internal/baz/baz_test.go":             &fstest.MapFile{Data: []byte("package baz")},
		"internal/qux/qux_integration_test.go": &fstest.MapFile{Data: []byte("package qux")},
	}
	c := &checktesting.UnitTestsExistChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "5") {
		t.Errorf("evidence should mention count 5, got: %s", res.Evidence)
	}
}

func TestUnitTestsExistNone(t *testing.T) {
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte("package main")},
	}
	c := &checktesting.UnitTestsExistChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestUnitTestsExistTS(t *testing.T) {
	repo := fstest.MapFS{
		"src/index.ts":    &fstest.MapFile{Data: []byte("export {}")},
		"src/foo.test.ts": &fstest.MapFile{Data: []byte("test('foo', () => {})")},
	}
	c := &checktesting.UnitTestsExistChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestUnitTestsRunnableGo(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project\n\nRun tests:\n\n```\ngo test ./...\n```\n")},
	}
	c := &checktesting.UnitTestsRunnableChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "README.md") {
		t.Errorf("evidence should mention README.md, got: %s", res.Evidence)
	}
}

func TestUnitTestsRunnableMissing(t *testing.T) {
	repo := fstest.MapFS{
		"README.md": &fstest.MapFile{Data: []byte("# My Project\n\nSome description.\n")},
	}
	c := &checktesting.UnitTestsRunnableChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}

func TestTestNamingGo(t *testing.T) {
	repo := fstest.MapFS{
		"main.go":                  &fstest.MapFile{Data: []byte("package main")},
		"main_test.go":             &fstest.MapFile{Data: []byte("package main")},
		"internal/foo/foo_test.go": &fstest.MapFile{Data: []byte("package foo")},
	}
	c := &checktesting.TestNamingConventionsChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
}

func TestTestNamingNoTests(t *testing.T) {
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte("package main")},
	}
	c := &checktesting.TestNamingConventionsChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Errorf("expected Skipped=true, got false; evidence: %s", res.Evidence)
	}
}

func TestTestIsolationGo(t *testing.T) {
	repo := fstest.MapFS{
		"main_test.go": &fstest.MapFile{Data: []byte(`package main

import "testing"

func TestFoo(t *testing.T) {
	t.Parallel()
}
`)},
	}
	c := &checktesting.TestIsolationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", res.Evidence)
	}
	if !strings.Contains(res.Evidence, "t.Parallel()") {
		t.Errorf("evidence should mention t.Parallel(), got: %s", res.Evidence)
	}
}

func TestTestIsolationMissing(t *testing.T) {
	repo := fstest.MapFS{
		"main_test.go": &fstest.MapFile{Data: []byte(`package main

import "testing"

func TestFoo(t *testing.T) {
	_ = t
}
`)},
	}
	c := &checktesting.TestIsolationChecker{}
	res, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", res.Evidence)
	}
}
