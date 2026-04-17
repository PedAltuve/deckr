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
	Get(ctx context.Context, name string) (ManagedTool, error)
}

type Backend interface {
	InitTool(ctx context.Context, name string) (string, error)
	ImportDeck(ctx context.Context, toolName, repoPath, deckName, sourcePath string) error
	ActivateDeck(ctx context.Context, targetPath, repoPath, deckName string) error
	CreateDeck(ctx context.Context, repoPath, deckName, fromDeck string) error
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

type CreateInput struct {
	Tool     string
	NewDeck  string
	FromDeck string
	Empty    bool
}

type CreateResult struct {
	Tool       string
	Deck       string
	SourceDeck string
	ActiveDeck string
}

type SwitchInput struct {
	Tool   string
	ToDeck string
}

type SwitchResult struct {
	Tool       string
	ActiveDeck string
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

func (s *Service) Current(ctx context.Context, tool string) (string, error) {
	if tool == "" {
		return "", fmt.Errorf("tool name is required")
	}

	managedTool, err := s.Registry.Get(ctx, tool)
	if err != nil {
		return "", err
	}

	return managedTool.ActiveDeck, nil
}

func (s *Service) Create(ctx context.Context, in CreateInput) (CreateResult, error) {
	if in.Tool == "" {
		return CreateResult{}, fmt.Errorf("tool name is required")
	}
	if in.NewDeck == "" {
		return CreateResult{}, fmt.Errorf("new deck name is required")
	}
	if in.Empty && in.FromDeck != "" {
		return CreateResult{}, fmt.Errorf("fromDeck and empty are mutually exclusive")
	}

	tool, err := s.Registry.Get(ctx, in.Tool)
	if err != nil {
		return CreateResult{}, fmt.Errorf("get managed tool: %w", err)
	}

	sourceDeck := in.FromDeck

	switch {
	case in.Empty:
		sourceDeck = ""
	case sourceDeck == "":
		sourceDeck = tool.ActiveDeck
	}

	if sourceDeck != "" && sourceDeck == in.NewDeck {
		return CreateResult{}, fmt.Errorf("new deck and source deck cannot be the name")
	}

	if err := s.Backend.CreateDeck(ctx, tool.RepoPath, in.NewDeck, sourceDeck); err != nil {
		return CreateResult{}, fmt.Errorf("create deck: %w", err)
	}

	return CreateResult{
		Tool:       tool.Name,
		Deck:       in.NewDeck,
		SourceDeck: sourceDeck,
		ActiveDeck: tool.ActiveDeck,
	}, nil
}

func (s *Service) Switch(ctx context.Context, in SwitchInput) (SwitchResult, error) {
	if in.Tool == "" {
		return SwitchResult{}, fmt.Errorf("tool name is required")
	}
	if in.ToDeck == "" {
		return SwitchResult{}, fmt.Errorf("new deck name is required")
	}

	tool, err := s.Registry.Get(ctx, in.Tool)
	if err != nil {
		return SwitchResult{}, fmt.Errorf("get managed tool: %w", err)
	}

	if err := s.Backend.ActivateDeck(ctx, tool.TargetPath, tool.RepoPath, in.ToDeck); err != nil {
		return SwitchResult{}, fmt.Errorf("activate %s deck: %w", in.ToDeck, err)
	}

	previousActiveDeck := tool.ActiveDeck
	tool.ActiveDeck = in.ToDeck

	if err := s.Registry.Save(ctx, tool); err != nil {
		if restoreErr := s.Backend.ActivateDeck(ctx, tool.TargetPath, tool.RepoPath, previousActiveDeck); restoreErr != nil {
			return SwitchResult{}, fmt.Errorf("save failed, rollback failed")
		}
		return SwitchResult{}, fmt.Errorf("save %s deck: %w", in.ToDeck, err)
	}

	return SwitchResult{
		Tool:       in.Tool,
		ActiveDeck: in.ToDeck,
	}, nil
}
