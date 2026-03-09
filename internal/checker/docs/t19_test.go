package docs

import (
	"testing"
	"testing/fstest"

	"github.com/bbik/ari/internal/checker"
)

// ── SERVICE FLOW DOCUMENTED ───────────────────────────────────────────────────

func TestServiceFlowDocumented(t *testing.T) {
	repo := fstest.MapFS{
		"docs/architecture/system.puml": &fstest.MapFile{Data: []byte("@startuml\nactor User\n@enduml")},
	}
	c := &ServiceFlowDocumentedChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestServiceFlowMermaid(t *testing.T) {
	repo := fstest.MapFS{
		"docs/design/flow.mmd": &fstest.MapFile{Data: []byte("graph TD\n  A --> B")},
	}
	c := &ServiceFlowDocumentedChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true (mermaid), got false; evidence: %s", r.Evidence)
	}
}

func TestServiceFlowMissing(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo")},
	}
	c := &ServiceFlowDocumentedChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed {
		t.Errorf("expected Passed=false, got true; evidence: %s", r.Evidence)
	}
}

// ── API SCHEMA DOCS ───────────────────────────────────────────────────────────

func TestApiSchemaDocsCLI(t *testing.T) {
	// Go repo with no net/http import → should be skipped.
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main

import "fmt"

func main() { fmt.Println("hello") }
`)},
	}
	c := &ApiSchemaDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Skipped {
		t.Errorf("expected Skipped=true (no net/http), got Skipped=false; evidence: %s", r.Evidence)
	}
}

func TestApiSchemaDocsFound(t *testing.T) {
	// Go repo with net/http AND openapi.yaml present.
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main

import "net/http"

func main() { http.ListenAndServe(":8080", nil) }
`)},
		"openapi.yaml": &fstest.MapFile{Data: []byte("openapi: 3.0.0\ninfo:\n  title: My API")},
	}
	c := &ApiSchemaDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Passed {
		t.Errorf("expected Passed=true, got false; evidence: %s", r.Evidence)
	}
}

func TestApiSchemaDocsAPINoSchema(t *testing.T) {
	// Go repo with net/http but no schema file → should fail.
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main

import "net/http"

func main() { http.ListenAndServe(":8080", nil) }
`)},
	}
	c := &ApiSchemaDocsChecker{}
	r, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if r.Passed || r.Skipped {
		t.Errorf("expected Passed=false, Skipped=false; got Passed=%v Skipped=%v; evidence: %s", r.Passed, r.Skipped, r.Evidence)
	}
}
