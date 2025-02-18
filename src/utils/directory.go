package utils

import (
	"fmt"
	"os"
)

// EnsureDirectoriesExist ensures that all given directory paths exist.
// If a directory doesn't exist, it will be created with permissions 0755.
func EnsureDirectoriesExist(dirs ...string) error {
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
