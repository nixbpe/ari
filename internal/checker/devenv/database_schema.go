package devenv

import (
	"context"
	"io/fs"

	"github.com/nixbpe/ari/internal/checker"
)

type DatabaseSchemaChecker struct{}

func (c *DatabaseSchemaChecker) ID() checker.CheckerID  { return "database_schema" }
func (c *DatabaseSchemaChecker) Pillar() checker.Pillar { return checker.PillarEnvInfra }
func (c *DatabaseSchemaChecker) Level() checker.Level   { return checker.LevelOptimized }
func (c *DatabaseSchemaChecker) Name() string           { return "Database Schema" }
func (c *DatabaseSchemaChecker) Description() string {
	return "Checks that database schema or migrations are defined and tracked in version control"
}
func (c *DatabaseSchemaChecker) Suggestion() string {
	return "Add database migrations (e.g., migrations/ directory) or schema files (schema.sql, prisma/schema.prisma) to version control"
}

func (c *DatabaseSchemaChecker) Check(ctx context.Context, repo fs.FS, lang checker.Language) (*checker.Result, error) {
	result := &checker.Result{
		ID:         c.ID(),
		Name:       c.Name(),
		Level:      c.Level(),
		Pillar:     c.Pillar(),
		Mode:       "rule-based",
		Suggestion: c.Suggestion(),
	}

	migrationDirs := []string{"migrations", "db/migrate"}
	for _, dir := range migrationDirs {
		if _, err := fs.ReadDir(repo, dir); err == nil {
			result.Passed = true
			result.Evidence = "Found " + dir + "/ directory"
			return result, nil
		}
	}

	schemaFiles := []string{
		"db/schema.sql",
		"schema.sql",
		"prisma/schema.prisma",
		"sqlc.yaml",
		"flyway.conf",
	}
	found, path := checker.FileExistsAny(repo, schemaFiles)
	if found {
		result.Passed = true
		result.Evidence = "Found " + path
	} else {
		result.Passed = false
		result.Evidence = "No database schema or migration files found"
	}
	return result, nil
}
