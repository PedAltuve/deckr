package tools

import (
	"context"
	"fmt"
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
	_ = sourcePath

	if deckName == "" {
		return fmt.Errorf("deck name is required")
	}

	deckPath := filepath.Join(repoPath, "decks", deckName)

	if err := os.MkdirAll(deckPath, 0o755); err != nil {
		return fmt.Errorf("create deck directory: %w", err)
	}

	return fmt.Errorf("ImportDeck not implemented yet")
}

func (b *LocalBackend) ActivateDeck(ctx context.Context, targetPath, repoPath, deckName string) error {
	_ = ctx
	_ = targetPath
	_ = repoPath
	_ = deckName

	return fmt.Errorf("ActivateDeck not implemented yet")
}
