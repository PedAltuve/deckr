package tools

import (
	"context"
	"testing"
)

type fakePaths struct {
	abs    string
	absErr error
	isDir  bool
	dirErr error
}

func (f fakePaths) Abs(path string) (string, error) {
	return f.abs, f.absErr
}
func (f fakePaths) IsDir(path string) (bool, error) {
	return f.isDir, f.dirErr
}

type fakeRegistry struct {
	exists        bool
	existsErr     error
	saveErr       error
	saved         ManagedTool
	saveWasCalled bool
}

func (f *fakeRegistry) Exists(ctx context.Context, name string) (bool, error) {
	return f.exists, f.existsErr
}
func (f *fakeRegistry) Save(ctx context.Context, tool ManagedTool) error {
	f.saveWasCalled = true
	f.saved = tool
	return f.saveErr
}

type fakeBackend struct {
	repoPath     string
	initErr      error
	importErr    error
	activateErr  error
	importCalled bool
	activeCalled bool
}

func (f *fakeBackend) InitTool(ctx context.Context, name string) (string, error) {
	return f.repoPath, f.initErr
}
func (f *fakeBackend) ImportDeck(ctx context.Context, toolName, repoPath, deckName, sourcePath string) error {
	f.importCalled = true
	return f.importErr
}
func (f *fakeBackend) ActivateDeck(ctx context.Context, targetPath, repoPath, deckName string) error {
	f.activeCalled = true
	return f.activateErr
}

func TestServiceInitSuccess(t *testing.T) {
	reg := &fakeRegistry{}
	backend := &fakeBackend{repoPath: "/repo/nvim"}
	svc := &Service{
		Paths: fakePaths{
			abs:   "/configs/nvim",
			isDir: true,
		},
		Registry: reg,
		Backend:  backend,
	}
	got, err := svc.Init(context.Background(), InitInput{
		Name:       "nvim",
		TargetPath: "./nvim",
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if got.Name != "nvim" {
		t.Fatalf("Name = %q, want %q", got.Name, "nvim")
	}
	if got.TargetPath != "/configs/nvim" {
		t.Fatalf("TargetPath = %q, want %q", got.TargetPath, "/configs/nvim")
	}
	if got.ActiveDeck != "default" {
		t.Fatalf("ActiveDeck = %q, want %q", got.ActiveDeck, "default")
	}
	if got.RepoPath != "/repo/nvim" {
		t.Fatalf("RepoPath = %q, want %q", got.RepoPath, "/repo/nvim")
	}
	if !backend.importCalled {
		t.Fatal("expected ImportDeck to be called")
	}
	if !backend.activeCalled {
		t.Fatal("expected ActivateDeck to be called")
	}
	if !reg.saveWasCalled {
		t.Fatal("expected Save to be called")
	}
}

func TestServiceInitRejectsAlreadyManagedTool(t *testing.T) {
	svc := &Service{
		Paths: fakePaths{
			abs:   "/configs/nvim",
			isDir: true,
		},
		Registry: &fakeRegistry{
			exists: true,
		},
		Backend: &fakeBackend{},
	}
	_, err := svc.Init(context.Background(), InitInput{
		Name:       "nvim",
		TargetPath: "./nvim",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
