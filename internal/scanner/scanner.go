package scanner

import (
	"context"
	"io/fs"

	"github.com/bbik/ari/internal/checker"
)

// FileInfo holds metadata about a scanned file
type FileInfo struct {
	Path      string
	Size      int64
	IsDir     bool
	Extension string
}

// RepoInfo holds the results of scanning a repository
type RepoInfo struct {
	Files      []FileInfo
	Language   checker.Language
	IsGitRepo  bool
	RootPath   string
	CommitHash string
	Branch     string
}

// Scanner scans a repository and returns metadata
type Scanner interface {
	Scan(ctx context.Context, repo fs.FS) (*RepoInfo, error)
}
