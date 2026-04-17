package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileRegistry struct {
	baseDir string
}

func NewFileRegistry(baseDir string) *FileRegistry {
	return &FileRegistry{baseDir: baseDir}
}

func (r *FileRegistry) Exists(ctx context.Context, name string) (bool, error) {
	_ = ctx

	toolFile := r.toolFile(name)

	_, err := os.Stat(toolFile)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("stat tool file: %w", err)
}

func (r *FileRegistry) Save(ctx context.Context, tool ManagedTool) error {
	_ = ctx

	if tool.Name == "" {
		return fmt.Errorf("tool name is required")
	}

	if err := os.MkdirAll(r.baseDir, 0o755); err != nil {
		return fmt.Errorf("create registry directory: %w", err)
	}

	data, err := json.MarshalIndent(tool, "", "	")
	if err != nil {
		return fmt.Errorf("marshal tool metadata: %w", err)
	}

	if err := os.WriteFile(r.toolFile(tool.Name), data, 0o644); err != nil {
		return fmt.Errorf("write tool metadata: %w", err)
	}

	return nil
}

func (r *FileRegistry) Get(ctx context.Context, name string) (ManagedTool, error) {
	path := r.toolFile(name)

	data, err := os.ReadFile(path)
	if err != nil {
		return ManagedTool{}, fmt.Errorf("error reading file: %w", err)
	}

	var m ManagedTool

	if err := json.Unmarshal(data, &m); err != nil {
		return ManagedTool{}, fmt.Errorf("unmarshal tool metadata: %w", err)
	}

	return m, nil
}

func (r *FileRegistry) toolFile(name string) string {
	return filepath.Join(r.baseDir, name+".json")
}
