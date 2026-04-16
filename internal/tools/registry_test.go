package tools

import (
	"context"
	"os"
	"testing"
)

func TestFileRegistryExists(t *testing.T) {
	t.Run("missing file returns false", func(t *testing.T) {
		reg := NewFileRegistry(t.TempDir())

		exists, err := reg.Exists(context.Background(), "nvim")
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}

		if exists {
			t.Fatalf("Exists() = true, want false")
		}
	})

	t.Run("existing file returns true", func(t *testing.T) {
		dir := t.TempDir()
		reg := NewFileRegistry(dir)

		err := os.WriteFile(reg.toolFile("nvim"), []byte(`{}`), 0o644)
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}

		exists, err := reg.Exists(context.Background(), "nvim")
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if !exists {
			t.Fatalf("Exists() = false, want true")
		}
	})
}
