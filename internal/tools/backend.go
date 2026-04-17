package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalBackend struct {
	baseDir string
}

func NewLocalBackend(baseDir string) *LocalBackend {
	return &LocalBackend{baseDir: baseDir}
}
func (b *LocalBackend) InitTool(ctx context.Context, name string) (string, error) {
	_ = ctx
	if name == "" {
		return "", fmt.Errorf("tool name is required")
	}
	repoPath := filepath.Join(b.baseDir, name)
	decksPath := filepath.Join(repoPath, "decks")
	if err := os.MkdirAll(decksPath, 0o755); err != nil {
		return "", fmt.Errorf("create backend directories: %w", err)
	}
	return repoPath, nil
}
func (b *LocalBackend) ImportDeck(ctx context.Context, toolName, repoPath, deckName, sourcePath string) error {
	_ = ctx
	_ = toolName
	if deckName == "" {
		return fmt.Errorf("deck name is required")
	}
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("stat source path: %w", err)
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("source path must be a directory")
	}
	deckPath := filepath.Join(repoPath, "decks", deckName)
	if _, err := os.Stat(deckPath); err == nil {
		return fmt.Errorf("deck path already exists: %s", deckPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat deck path: %w", err)
	}
	if err := os.MkdirAll(deckPath, 0o755); err != nil {
		return fmt.Errorf("create deck directory: %w", err)
	}
	if err := copyDir(sourcePath, deckPath); err != nil {
		return fmt.Errorf("copy source config into deck: %w", err)
	}
	return nil
}
func (b *LocalBackend) ActivateDeck(ctx context.Context, targetPath, repoPath, deckName string) error {
	_ = ctx
	if deckName == "" {
		return fmt.Errorf("deck name is required")
	}
	if targetPath == "" {
		return fmt.Errorf("target path is required")
	}
	deckPath := filepath.Join(repoPath, "decks", deckName)
	deckPath, err := filepath.Abs(deckPath)
	if err != nil {
		return fmt.Errorf("resolve deck path: %w", err)
	}
	deckInfo, err := os.Stat(deckPath)
	if err != nil {
		return fmt.Errorf("stat deck path: %w", err)
	}
	if !deckInfo.IsDir() {
		return fmt.Errorf("deck path must be a directory")
	}
	targetInfo, err := os.Lstat(targetPath)
	switch {
	case err == nil:
		if targetInfo.Mode()&os.ModeSymlink != 0 {
			currentTarget, err := os.Readlink(targetPath)
			if err != nil {
				return fmt.Errorf("read existing symlink: %w", err)
			}
			if !filepath.IsAbs(currentTarget) {
				currentTarget = filepath.Join(filepath.Dir(targetPath), currentTarget)
			}
			currentTarget, err = filepath.Abs(currentTarget)
			if err != nil {
				return fmt.Errorf("resolve existing symlink target: %w", err)
			}
			if currentTarget == deckPath {
				return nil
			}
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("remove existing symlink: %w", err)
			}
		} else {
			backupPath := targetPath + ".deckr.bak"
			if _, err := os.Lstat(backupPath); err == nil {
				return fmt.Errorf("backup path already exists: %s", backupPath)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("check backup path: %w", err)
			}
			if err := os.Rename(targetPath, backupPath); err != nil {
				return fmt.Errorf("backup target path: %w", err)
			}
			if err := os.Symlink(deckPath, targetPath); err != nil {
				rollbackErr := os.Rename(backupPath, targetPath)
				if rollbackErr != nil {
					return fmt.Errorf("create symlink: %w (rollback failed: %v)", err, rollbackErr)
				}
				return fmt.Errorf("create symlink: %w", err)
			}
			return nil
		}
	case errors.Is(err, os.ErrNotExist):
		// target path does not exist, continue
	default:
		return fmt.Errorf("inspect target path: %w", err)
	}
	if err := os.Symlink(deckPath, targetPath); err != nil {
		return fmt.Errorf("create symlink: %w", err)
	}
	return nil
}

func (b *LocalBackend) CreateDeck(ctx context.Context, repoPath, deckName, fromDeck string) error {
	_ = ctx

	if repoPath == "" {
		return fmt.Errorf("repo path is required")
	}
	if deckName == "" {
		return fmt.Errorf("deck name is required")
	}

	newDeckPath := filepath.Join(repoPath, "decks", deckName)

	if _, err := os.Stat(newDeckPath); err == nil {
		return fmt.Errorf("deck already exists: %s", deckName)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat new deck path: %w", err)
	}

	if err := os.MkdirAll(newDeckPath, 0o755); err != nil {
		return fmt.Errorf("create new deck directory: %w", err)
	}

	if fromDeck == "" {
		return nil
	}

	sourceDeckPath := filepath.Join(repoPath, "decks", fromDeck)

	sourceInfo, err := os.Stat(sourceDeckPath)
	if err != nil {
		removeErr := os.RemoveAll(newDeckPath)
		if removeErr != nil {
			return fmt.Errorf("stat source deck: %w (cleanup failed: %v)", err, removeErr)
		}
		return fmt.Errorf("stat source deck: %w", err)
	}

	if !sourceInfo.IsDir() {
		removeErr := os.RemoveAll(newDeckPath)
		if removeErr != nil {
			return fmt.Errorf("source deck path must be a directory (cleanup failed: %v", removeErr)
		}
		return fmt.Errorf("source and destination deck cannot be the same")
	}

	if err := copyDir(sourceDeckPath, newDeckPath); err != nil {
		removeErr := os.RemoveAll(newDeckPath)
		if removeErr != nil {
			return fmt.Errorf("copy source deck into new deck: %w (cleanup failed: %v)", err, removeErr)
		}
		return fmt.Errorf("copy source deck into new deck: %w", err)
	}
	return nil
}

func copyDir(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", srcDir, err)
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks are not supported during import: %s", srcPath)
		}
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("stat directory %s: %w", srcPath, err)
			}
			if err := os.MkdirAll(dstPath, info.Mode().Perm()); err != nil {
				return fmt.Errorf("create directory %s: %w", dstPath, err)
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat file %s: %w", srcPath, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported file type: %s", srcPath)
		}
		if err := copyFile(srcPath, dstPath, info.Mode().Perm()); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(srcPath, dstPath string, perm os.FileMode) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source file %s: %w", srcPath, err)
	}
	defer func() {
		closeErr := srcFile.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open destination file %s: %w", dstPath, err)
	}
	defer func() {
		closeErr := dstFile.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy file %s -> %s: %w", srcPath, dstPath, err)
	}
	return nil
}
