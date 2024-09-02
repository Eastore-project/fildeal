package dealutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// getReaders returns a slice of io.ReadSeeker for all files in the given folder
func GetReaders(folder string) ([]io.ReadSeeker, error) {
    readers := make([]io.ReadSeeker, 0)
    files, err := os.ReadDir(folder)
    if err != nil {
        return nil, fmt.Errorf("failed to read folder %s: %w", folder, err)
    }
    for _, file := range files {
        if file.Type().IsRegular() {
            filePath := filepath.Join(folder, file.Name())
            r, err := os.Open(filePath)
            if err != nil {
                return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
            }
            readers = append(readers, r)
        }
    }
    return readers, nil
}