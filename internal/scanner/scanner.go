package scanner

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/bbik/ari/internal/checker"
)

const DefaultFileLimit = 5000

var ignoreDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".idea":        true,
	"build":        true,
	"dist":         true,
	".next":        true,
	"__pycache__":  true,
	".cache":       true,
	"target":       true,
}

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

type DefaultScanner struct {
	FileLimit int
}

func NewScanner() *DefaultScanner {
	return &DefaultScanner{FileLimit: DefaultFileLimit}
}

func (s *DefaultScanner) Scan(ctx context.Context, repo fs.FS) (*RepoInfo, error) {
	info := &RepoInfo{}
	count := 0

	limit := s.FileLimit
	if limit <= 0 {
		limit = DefaultFileLimit
	}

	err := fs.WalkDir(repo, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			log.Printf("warning: skipping %s: %v", path, walkErr)
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if d.IsDir() {
			if ignoreDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}

		if count >= limit {
			return nil
		}

		fi, statErr := d.Info()
		if statErr != nil {
			log.Printf("warning: could not stat %s: %v", path, statErr)
			return nil
		}

		info.Files = append(info.Files, FileInfo{
			Path:      path,
			Size:      fi.Size(),
			IsDir:     false,
			Extension: strings.ToLower(filepath.Ext(path)),
		})
		count++
		return nil
	})

	if err != nil {
		if ctx.Err() != nil {
			return nil, err
		}
		return nil, err
	}

	info.Language = DetectLanguage(info.Files)
	info.IsGitRepo = HasGitDirFS(repo)

	return info, nil
}
