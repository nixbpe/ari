package observability

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
)

var ctx = context.Background()

func TestStructuredLoggingGoFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire go.uber.org/zap v1.0.0\n")},
	}
	c := &StructuredLoggingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for zap in go.mod, got evidence: %s", result.Evidence)
	}
}

func TestStructuredLoggingNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/pkg/errors v1.0.0\n")},
	}
	c := &StructuredLoggingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no logging library, got evidence: %s", result.Evidence)
	}
}

func TestStructuredLoggingTypeScriptFound(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"winston":"^3.0.0"}}`)},
	}
	c := &StructuredLoggingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for winston in package.json, got evidence: %s", result.Evidence)
	}
}

func TestStructuredLoggingPinoFound(t *testing.T) {
	repo := fstest.MapFS{
		"package.json": &fstest.MapFile{Data: []byte(`{"dependencies":{"pino":"^8.0.0"}}`)},
	}
	c := &StructuredLoggingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for pino in package.json, got evidence: %s", result.Evidence)
	}
}

func TestHealthChecksGoFound(t *testing.T) {
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main
func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}`)},
	}
	c := &HealthChecksChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for /healthz route, got evidence: %s", result.Evidence)
	}
}

func TestHealthChecksNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main
func main() {
	http.HandleFunc("/api/users", handler)
}`)},
	}
	c := &HealthChecksChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no health route, got evidence: %s", result.Evidence)
	}
}

func TestHealthChecksTypeScriptFound(t *testing.T) {
	repo := fstest.MapFS{
		"server.ts": &fstest.MapFile{Data: []byte(`app.get('/health', (req, res) => { res.json({ok: true}) })`)},
	}
	c := &HealthChecksChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for /health in ts file, got evidence: %s", result.Evidence)
	}
}

func TestErrorTrackingFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/getsentry/sentry-go v0.18.0\n")},
	}
	c := &ErrorTrackingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for sentry in go.mod, got evidence: %s", result.Evidence)
	}
}

func TestErrorTrackingNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/pkg/errors v1.0.0\n")},
	}
	c := &ErrorTrackingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no error tracking, got evidence: %s", result.Evidence)
	}
}

func TestDistributedTracingFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire go.opentelemetry.io/otel v1.0.0\n")},
	}
	c := &DistributedTracingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for opentelemetry in go.mod, got evidence: %s", result.Evidence)
	}
}

func TestDistributedTracingNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/gorilla/mux v1.8.0\n")},
	}
	c := &DistributedTracingChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no tracing library, got evidence: %s", result.Evidence)
	}
}

func TestMetricsCollectionFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/prometheus/client_golang v1.15.0\n")},
	}
	c := &MetricsCollectionChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for prometheus in go.mod, got evidence: %s", result.Evidence)
	}
}

func TestMetricsCollectionNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/gorilla/mux v1.8.0\n")},
	}
	c := &MetricsCollectionChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no metrics library, got evidence: %s", result.Evidence)
	}
}

func TestAlertingConfiguredFileFound(t *testing.T) {
	repo := fstest.MapFS{
		"alertmanager.yml": &fstest.MapFile{Data: []byte("route:\n  receiver: default\n")},
	}
	c := &AlertingConfiguredChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for alertmanager.yml, got evidence: %s", result.Evidence)
	}
}

func TestAlertingConfiguredDirFound(t *testing.T) {
	repo := fstest.MapFS{
		"alerts/cpu_high.yml": &fstest.MapFile{Data: []byte("alert: CPUHigh\n")},
	}
	c := &AlertingConfiguredChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for alerts/ directory, got evidence: %s", result.Evidence)
	}
}

func TestAlertingConfiguredNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\n")},
	}
	c := &AlertingConfiguredChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no alerting config, got evidence: %s", result.Evidence)
	}
}

func TestProfilingInstrumentationPprofFound(t *testing.T) {
	repo := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main
import _ "net/http/pprof"
func main() {}`)},
		"go.mod": &fstest.MapFile{Data: []byte("module foo\n")},
	}
	c := &ProfilingInstrumentationChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for net/http/pprof import, got evidence: %s", result.Evidence)
	}
}

func TestProfilingInstrumentationFgprofFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/felixge/fgprof v0.9.3\n")},
	}
	c := &ProfilingInstrumentationChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for fgprof in go.mod, got evidence: %s", result.Evidence)
	}
}

func TestProfilingInstrumentationNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\nrequire github.com/gorilla/mux v1.8.0\n")},
		"main.go": &fstest.MapFile{Data: []byte(`package main
import "net/http"
func main() { http.ListenAndServe(":8080", nil) }`)},
	}
	c := &ProfilingInstrumentationChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no profiling, got evidence: %s", result.Evidence)
	}
}

func TestRunbooksDocumentedFileFound(t *testing.T) {
	repo := fstest.MapFS{
		"RUNBOOK.md": &fstest.MapFile{Data: []byte("# Runbook\n\n## Incident Response\n\nThis runbook describes how to handle incidents in production.\n")},
	}
	c := &RunbooksDocumentedChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for RUNBOOK.md, got evidence: %s", result.Evidence)
	}
}

func TestRunbooksDocumentedDirFound(t *testing.T) {
	repo := fstest.MapFS{
		"runbooks/on-call.md":   &fstest.MapFile{Data: []byte("# On-Call Guide\n")},
		"runbooks/incidents.md": &fstest.MapFile{Data: []byte("# Incidents\n")},
	}
	c := &RunbooksDocumentedChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for runbooks/ directory, got evidence: %s", result.Evidence)
	}
}

func TestRunbooksDocumentedNotFound(t *testing.T) {
	repo := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module foo\n")},
	}
	c := &RunbooksDocumentedChecker{}
	result, err := c.Check(ctx, repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail for no runbooks, got evidence: %s", result.Evidence)
	}
}

func TestAllCheckersReturnObservabilityPillar(t *testing.T) {
	checkers := []checker.Checker{
		&StructuredLoggingChecker{},
		&HealthChecksChecker{},
		&ErrorTrackingChecker{},
		&DistributedTracingChecker{},
		&MetricsCollectionChecker{},
		&AlertingConfiguredChecker{},
		&ProfilingInstrumentationChecker{},
		&RunbooksDocumentedChecker{},
	}
	for _, c := range checkers {
		if c.Pillar() != checker.PillarObservability {
			t.Errorf("checker %s returned pillar %v, want PillarObservability", c.ID(), c.Pillar())
		}
	}
}
