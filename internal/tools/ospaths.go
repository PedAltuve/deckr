package tools

import (
	"errors"
	"os"
	"path/filepath"
)

type OSPaths struct{}

func (OSPaths) IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

func (OSPaths) Abs(path string) (string, error) {
	return filepath.Abs(path)
}
