package scanner

import (
	"context"
	"io/fs"
	"os/exec"
	"strings"
)

func DetectGitRepo(ctx context.Context, rootPath string) (isGitRepo bool, commitHash, branch string) {
	out, err := runGit(ctx, rootPath, "rev-parse", "--is-inside-work-tree")
	if err != nil || strings.TrimSpace(out) != "true" {
		return false, "", ""
	}

	commitHash, branch = DetectGitInfo(ctx, rootPath)
	return true, commitHash, branch
}

func DetectGitInfo(ctx context.Context, rootPath string) (commitHash, branch string) {
	if hash, err := runGit(ctx, rootPath, "rev-parse", "--short", "HEAD"); err == nil {
		commitHash = strings.TrimSpace(hash)
	}

	if br, err := runGit(ctx, rootPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		branch = strings.TrimSpace(br)
	}

	return commitHash, branch
}

func runGit(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}

func HasGitDirFS(repo fs.FS) bool {
	_, err := fs.Stat(repo, ".git")
	return err == nil
}
