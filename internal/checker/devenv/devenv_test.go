package devenv_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/devenv"
)

func TestEnvTemplateFound(t *testing.T) {
	repo := fstest.MapFS{
		".env.example": &fstest.MapFile{Data: []byte("DATABASE_URL=\nAPI_KEY=")},
	}
	c := &devenv.EnvTemplateChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Evidence)
	}
}

func TestEnvTemplateNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.EnvTemplateChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestEnvTemplateSampleVariant(t *testing.T) {
	repo := fstest.MapFS{
		".env.sample": &fstest.MapFile{Data: []byte("PORT=3000")},
	}
	c := &devenv.EnvTemplateChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for .env.sample, got fail: %s", result.Evidence)
	}
}

func TestDevcontainerFound(t *testing.T) {
	repo := fstest.MapFS{
		".devcontainer/devcontainer.json": &fstest.MapFile{Data: []byte(`{"name":"dev"}`)},
	}
	c := &devenv.DevcontainerChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Evidence)
	}
}

func TestDevcontainerNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.DevcontainerChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestVersionPinningNvmrc(t *testing.T) {
	repo := fstest.MapFS{
		".nvmrc": &fstest.MapFile{Data: []byte("18.17.0")},
	}
	c := &devenv.VersionPinningChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for .nvmrc, got fail: %s", result.Evidence)
	}
}

func TestVersionPinningToolVersions(t *testing.T) {
	repo := fstest.MapFS{
		".tool-versions": &fstest.MapFile{Data: []byte("nodejs 18.17.0\npython 3.11.0")},
	}
	c := &devenv.VersionPinningChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for .tool-versions, got fail: %s", result.Evidence)
	}
}

func TestVersionPinningNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.VersionPinningChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestLocalServicesSetupFound(t *testing.T) {
	repo := fstest.MapFS{
		"docker-compose.yml": &fstest.MapFile{Data: []byte("services:\n  db:\n    image: postgres")},
	}
	c := &devenv.LocalServicesSetupChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass, got fail: %s", result.Evidence)
	}
}

func TestLocalServicesSetupNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.LocalServicesSetupChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestIDEConfigVSCode(t *testing.T) {
	repo := fstest.MapFS{
		".vscode/settings.json": &fstest.MapFile{Data: []byte(`{"editor.formatOnSave": true}`)},
	}
	c := &devenv.IDEConfigChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for .vscode/settings.json, got fail: %s", result.Evidence)
	}
}

func TestIDEConfigEditorConfig(t *testing.T) {
	repo := fstest.MapFS{
		".editorconfig": &fstest.MapFile{Data: []byte("[*]\nindent_size = 2")},
	}
	c := &devenv.IDEConfigChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for .editorconfig, got fail: %s", result.Evidence)
	}
}

func TestIDEConfigNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.IDEConfigChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestDevcontainerQualityWithPostCreateCommand(t *testing.T) {
	repo := fstest.MapFS{
		".devcontainer/devcontainer.json": &fstest.MapFile{Data: []byte(`{
			"name": "Go Dev",
			"postCreateCommand": "go mod download",
			"extensions": ["golang.go"]
		}`)},
	}
	c := &devenv.DevcontainerQualityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass with postCreateCommand, got fail: %s", result.Evidence)
	}
}

func TestDevcontainerQualityMissingPostCreateCommand(t *testing.T) {
	repo := fstest.MapFS{
		".devcontainer/devcontainer.json": &fstest.MapFile{Data: []byte(`{"name": "Go Dev"}`)},
	}
	c := &devenv.DevcontainerQualityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail without postCreateCommand, got pass")
	}
}

func TestDevcontainerQualityNoDevcontainer(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.DevcontainerQualityChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail with no devcontainer, got pass")
	}
}

func TestDatabaseSchemaMigrationsDir(t *testing.T) {
	repo := fstest.MapFS{
		"migrations/001_init.sql": &fstest.MapFile{Data: []byte("CREATE TABLE users (id INT);")},
	}
	c := &devenv.DatabaseSchemaChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for migrations/ dir, got fail: %s", result.Evidence)
	}
}

func TestDatabaseSchemaPrisma(t *testing.T) {
	repo := fstest.MapFS{
		"prisma/schema.prisma": &fstest.MapFile{Data: []byte("datasource db { provider = \"postgresql\" }")},
	}
	c := &devenv.DatabaseSchemaChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageTypeScript)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed {
		t.Errorf("expected pass for prisma/schema.prisma, got fail: %s", result.Evidence)
	}
}

func TestDatabaseSchemaNotFound(t *testing.T) {
	repo := fstest.MapFS{}
	c := &devenv.DatabaseSchemaChecker{}
	result, err := c.Check(context.Background(), repo, checker.LanguageGo)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Errorf("expected fail, got pass")
	}
}

func TestAllCheckersReturnDevEnvPillar(t *testing.T) {
	repo := fstest.MapFS{}
	checkers := []checker.Checker{
		&devenv.EnvTemplateChecker{},
		&devenv.DevcontainerChecker{},
		&devenv.VersionPinningChecker{},
		&devenv.LocalServicesSetupChecker{},
		&devenv.IDEConfigChecker{},
		&devenv.DevcontainerQualityChecker{},
		&devenv.DatabaseSchemaChecker{},
	}
	for _, c := range checkers {
		if c.Pillar() != checker.PillarDevEnvironment {
			t.Errorf("checker %s returned pillar %v, want PillarDevEnvironment", c.ID(), c.Pillar())
		}
		result, err := c.Check(context.Background(), repo, checker.LanguageGo)
		if err != nil {
			t.Errorf("checker %s returned error: %v", c.ID(), err)
		}
		if result == nil {
			t.Errorf("checker %s returned nil result", c.ID())
		}
	}
}
