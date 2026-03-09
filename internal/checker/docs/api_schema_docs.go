package docs

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/nixbpe/ari/internal/checker"
)

// ApiSchemaDocsChecker checks whether the repository documents its API with an
// OpenAPI/Swagger specification or a GraphQL schema file.
type ApiSchemaDocsChecker struct{}

func (c *ApiSchemaDocsChecker) ID() checker.CheckerID  { return "api_schema_docs" }
func (c *ApiSchemaDocsChecker) Pillar() checker.Pillar { return checker.PillarDocumentation }
func (c *ApiSchemaDocsChecker) Level() checker.Level   { return checker.LevelStandardized }
func (c *ApiSchemaDocsChecker) Name() string           { return "API Schema Documentation" }
func (c *ApiSchemaDocsChecker) Description() string {
	return "Checks for OpenAPI/Swagger spec or GraphQL schema files documenting the repository API"
}
func (c *ApiSchemaDocsChecker) Suggestion() string {
	return "Document your API with OpenAPI/Swagger specification"
}

// schemaRoots are well-known top-level schema files.
var schemaRoots = []string{
	"openapi.yaml",
	"openapi.json",
	"swagger.yaml",
	"swagger.json",
	"schema.graphql",
}

// schemaDirs are directories that may contain GraphQL or API schema files.
var schemaDirs = []string{"api", "docs/api"}

func (c *ApiSchemaDocsChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	// Skip logic — if the repo has no HTTP-related indicator, skip.
	if shouldSkipAPICheck(repo, lang) {
		result.Skipped = true
		result.SkipReason = "No HTTP server indicators found — not an API service"
		return result, nil
	}

	// Check well-known root-level schema files.
	for _, name := range schemaRoots {
		if _, err := fs.Stat(repo, name); err == nil {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — API schema documented", name)
			return result, nil
		}
	}

	// Check for *.graphql files under api/ or docs/api/.
	for _, dir := range schemaDirs {
		if _, err := fs.Stat(repo, dir); err != nil {
			continue
		}
		found, name := findGraphQLFile(repo, dir)
		if found {
			result.Passed = true
			result.Evidence = fmt.Sprintf("Found %s — API schema documented", name)
			return result, nil
		}
	}

	result.Passed = false
	result.Evidence = "No API schema documentation found"
	return result, nil
}

// findGraphQLFile walks a directory tree and returns the first *.graphql file found.
func findGraphQLFile(repo fs.FS, dir string) (bool, string) {
	var found bool
	var foundPath string
	_ = fs.WalkDir(repo, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || found {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".graphql") {
			found = true
			foundPath = path
		}
		return nil
	})
	return found, foundPath
}

// shouldSkipAPICheck returns true when the language has no HTTP-server indicators.
//
// Rules:
//   - Java: never skip (always assume it exposes an API).
//   - Go: skip when no .go file imports "net/http".
//   - TypeScript: skip when package.json lacks express/fastify/hapi in dependencies.
//   - Unknown: don't skip (conservative).
func shouldSkipAPICheck(repo fs.FS, lang checker.Language) bool {
	switch lang {
	case checker.LanguageJava:
		return false
	case checker.LanguageGo:
		return !goHasNetHTTP(repo)
	case checker.LanguageTypeScript:
		return !tsHasHTTPFramework(repo)
	default:
		return false
	}
}

// goHasNetHTTP reports whether any .go file (non-test) imports "net/http".
func goHasNetHTTP(repo fs.FS) bool {
	var found bool
	_ = fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || found {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, err := fs.ReadFile(repo, path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), `"net/http"`) {
			found = true
		}
		return nil
	})
	return found
}

// tsHasHTTPFramework reports whether package.json lists express, fastify, or hapi
// in its dependencies or devDependencies.
func tsHasHTTPFramework(repo fs.FS) bool {
	data, err := fs.ReadFile(repo, "package.json")
	if err != nil {
		return false
	}
	content := string(data)
	for _, fw := range []string{"express", "fastify", "hapi"} {
		if strings.Contains(content, fw) {
			return true
		}
	}
	return false
}
