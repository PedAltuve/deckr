package tools

import (
	"context"
	"fmt"
)

type PathResolver interface {
	Abs(path string) (string, error)
	IsDir(path string) (bool, error)
}

type Registry interface {
	Exists(ctx context.Context, name string) (bool, error)
	Save(ctx context.Context, tool ManagedTool) error
}

type Backend interface {
	InitTool(ctx context.Context, name string) (string, error)
	ImportDeck(ctx context.Context, toolName, repoPath, deckName, sourcePath string) error
	ActivateDeck(ctx context.Context, targetPath, repoPath, deckName string) error
}

type Service struct {
	Paths    PathResolver
	Registry Registry
	Backend  Backend
}

type InitInput struct {
	Name       string
	TargetPath string
}

type InitResult struct {
	Name       string
	TargetPath string
	ActiveDeck string
	RepoPath   string
}

type ManagedTool struct {
	Name       string
	TargetPath string
	ActiveDeck string
	RepoPath   string
}

func (s *Service) Init(ctx context.Context, in InitInput) (InitResult, error) {
	if in.Name == "" {
		return InitResult{}, fmt.Errorf("tool name is required")
	}

	targetPath, err := s.Paths.Abs(in.TargetPath)
	if err != nil {
		return InitResult{}, fmt.Errorf("resolve target path: %w", err)
	}

	isDir, err := s.Paths.IsDir(targetPath)
	if err != nil {
		return InitResult{}, fmt.Errorf("check target path: %w", err)
	}
	if !isDir {
		return InitResult{}, fmt.Errorf("target path must be an existing directory")
	}

	exists, err := s.Registry.Exists(ctx, in.Name)
	if err != nil {
		return InitResult{}, fmt.Errorf("check managed tool: %w", err)
	}
	if exists {
		return InitResult{}, fmt.Errorf("tool %q is already managed", in.Name)
	}

	repoPath, err := s.Backend.InitTool(ctx, in.Name)
	if err != nil {
		return InitResult{}, fmt.Errorf("init backend: %w", err)
	}

	if err := s.Backend.ImportDeck(ctx, in.Name, repoPath, "default", targetPath); err != nil {
		return InitResult{}, fmt.Errorf("import default deck: %w", err)
	}

	if err := s.Backend.ActivateDeck(ctx, targetPath, repoPath, "default"); err != nil {
		return InitResult{}, fmt.Errorf("activate default deck: %w", err)
	}

	tool := ManagedTool{
		Name:       in.Name,
		TargetPath: targetPath,
		ActiveDeck: "default",
		RepoPath:   repoPath,
	}

	if err := s.Registry.Save(ctx, tool); err != nil {
		return InitResult{}, fmt.Errorf("save tool metadata: %w", err)
	}

	return InitResult{
		Name:       tool.Name,
		TargetPath: tool.TargetPath,
		ActiveDeck: tool.ActiveDeck,
		RepoPath:   tool.RepoPath,
	}, nil
}
