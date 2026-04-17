package tools

import (
	"context"
	"fmt"
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
	getWasCalled  bool
	getResult     ManagedTool
	getErr        error
	getName       string
}

func (f *fakeRegistry) Exists(ctx context.Context, name string) (bool, error) {
	return f.exists, f.existsErr
}
func (f *fakeRegistry) Save(ctx context.Context, tool ManagedTool) error {
	f.saveWasCalled = true
	f.saved = tool
	return f.saveErr
}
func (f *fakeRegistry) Get(ctx context.Context, name string) (ManagedTool, error) {
	f.getWasCalled = true
	f.getName = name
	return f.getResult, f.getErr
}

type fakeBackend struct {
	repoPath        string
	initErr         error
	importErr       error
	activateErr     error
	importCalled    bool
	activeCalled    bool
	createErr       error
	createWasCalled bool
	createRepoPath  string
	createDeckName  string
	createFromDeck  string
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
func (f *fakeBackend) CreateDeck(ctx context.Context, repoPath, deckName, fromDeck string) error {
	f.createWasCalled = true
	f.createRepoPath = repoPath
	f.createDeckName = deckName
	f.createFromDeck = fromDeck
	return f.createErr
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

func TestServiceCurrentSuccess(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
		},
	}

	svc := &Service{
		Registry: reg,
	}

	got, err := svc.Current(context.Background(), "nvim")
	if err != nil {
		t.Fatalf("Current() error = %v", err)
	}

	if got != "default" {
		t.Fatalf("Current() = %q, want %q", got, "default")
	}

	if !reg.getWasCalled {
		t.Fatal("expected Get to be called")
	}

	if reg.getName != "nvim" {
		t.Fatalf("Get() called with %q, want %q", reg.getName, "nvim")
	}
}

func TestServiceCurrentReturnsRegistryError(t *testing.T) {
	reg := &fakeRegistry{
		getErr: fmt.Errorf("boom"),
	}

	svc := &Service{
		Registry: reg,
	}

	_, err := svc.Current(context.Background(), "nvim")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !reg.getWasCalled {
		t.Fatal("expected Get to be called")
	}
}

func TestServiceCurrentRequiresToolName(t *testing.T) {
	reg := &fakeRegistry{}

	svc := &Service{
		Registry: reg,
	}

	_, err := svc.Current(context.Background(), "")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if reg.getWasCalled {
		t.Fatal("did not expected Get to be called")
	}
}

func TestServiceCreateUsesActiveDeckByDefault(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
			RepoPath:   "/repo/nvim",
		},
	}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	got, err := svc.Create(context.Background(), CreateInput{
		Tool:    "nvim",
		NewDeck: "work",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !reg.getWasCalled {
		t.Fatal("expected Get to be called")
	}
	if !backend.createWasCalled {
		t.Fatal("expected CreateDeck to be called")
	}
	if backend.createRepoPath != "/repo/nvim" {
		t.Fatalf("CreateDeck repoPath = %q, want %q", backend.createRepoPath, "/repo/nvim")
	}
	if backend.createDeckName != "work" {
		t.Fatalf("CreateDeck deckName = %q, want %q", backend.createDeckName, "work")
	}
	if backend.createFromDeck != "default" {
		t.Fatalf("CreateDeck fromDeck = %q, want %q", backend.createFromDeck, "default")
	}
	if got.Tool != "nvim" {
		t.Fatalf("Tool = %q, want %q", got.Tool, "nvim")
	}
	if got.Deck != "work" {
		t.Fatalf("Deck = %q, want %q", got.Deck, "work")
	}
	if got.SourceDeck != "default" {
		t.Fatalf("SourceDeck = %q, want %q", got.SourceDeck, "default")
	}
	if got.ActiveDeck != "default" {
		t.Fatalf("ActiveDeck = %q, want %q", got.ActiveDeck, "default")
	}
	if reg.saveWasCalled {
		t.Fatal("did not expect Save to be called")
	}
}

func TestServiceCreateUsesExplicitFromDeck(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
			RepoPath:   "/repo/nvim",
		},
	}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	got, err := svc.Create(context.Background(), CreateInput{
		Tool:     "nvim",
		NewDeck:  "work",
		FromDeck: "minimal",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if backend.createFromDeck != "minimal" {
		t.Fatalf("CreateDeck fromDeck = %q, want %q", backend.createFromDeck, "minimal")
	}
	if got.SourceDeck != "minimal" {
		t.Fatalf("SourceDeck = %q, want %q", got.SourceDeck, "minimal")
	}
}

func TestServiceCreateEmptyDeck(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
			RepoPath:   "/repo/nvim",
		},
	}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	got, err := svc.Create(context.Background(), CreateInput{
		Tool:    "nvim",
		NewDeck: "scratch",
		Empty:   true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if backend.createFromDeck != "" {
		t.Fatalf("CreateDeck fromDeck = %q, want empty string", backend.createFromDeck)
	}
	if got.SourceDeck != "" {
		t.Fatalf("SourceDeck = %q, want empty string", got.SourceDeck)
	}
}

func TestServiceCreateRejectsFromDeckAndEmptyTogether(t *testing.T) {
	reg := &fakeRegistry{}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		Tool:     "nvim",
		NewDeck:  "work",
		FromDeck: "default",
		Empty:    true,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if reg.getWasCalled {
		t.Fatal("did not expect Get to be called")
	}
	if backend.createWasCalled {
		t.Fatal("did not expect CreateDeck to be called")
	}
}

func TestServiceCreateRejectsSameSourceAndDestinationDeck(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
			RepoPath:   "/repo/nvim",
		},
	}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		Tool:     "nvim",
		NewDeck:  "default",
		FromDeck: "default",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if backend.createWasCalled {
		t.Fatal("did not expect CreateDeck to be called")
	}
}

func TestServiceCreateReturnsRegistryError(t *testing.T) {
	reg := &fakeRegistry{
		getErr: fmt.Errorf("boom"),
	}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		Tool:    "nvim",
		NewDeck: "work",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !reg.getWasCalled {
		t.Fatal("expected Get to be called")
	}
	if backend.createWasCalled {
		t.Fatal("did not expect CreateDeck to be called")
	}
}

func TestServiceCreateReturnsBackendError(t *testing.T) {
	reg := &fakeRegistry{
		getResult: ManagedTool{
			Name:       "nvim",
			ActiveDeck: "default",
			RepoPath:   "/repo/nvim",
		},
	}
	backend := &fakeBackend{
		createErr: fmt.Errorf("boom"),
	}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		Tool:    "nvim",
		NewDeck: "work",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !backend.createWasCalled {
		t.Fatal("expected CreateDeck to be called")
	}
}

func TestServiceCreateRequiresToolName(t *testing.T) {
	reg := &fakeRegistry{}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		NewDeck: "work",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if reg.getWasCalled {
		t.Fatal("did not expect Get to be called")
	}
	if backend.createWasCalled {
		t.Fatal("did not expect CreateDeck to be called")
	}
}
func TestServiceCreateRequiresNewDeckName(t *testing.T) {
	reg := &fakeRegistry{}
	backend := &fakeBackend{}
	svc := &Service{
		Registry: reg,
		Backend:  backend,
	}
	_, err := svc.Create(context.Background(), CreateInput{
		Tool: "nvim",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if reg.getWasCalled {
		t.Fatal("did not expect Get to be called")
	}
	if backend.createWasCalled {
		t.Fatal("did not expect CreateDeck to be called")
	}
}
