package utils

import (
	"bytes"
	"fmt"
	"os"
)


func FindOffset(fileAPath, fileBPath string) (int, error) {
    // Read File A
    fileA, err := os.ReadFile(fileAPath)
    if err != nil {
        return -1, fmt.Errorf("failed to read file A: %w", err)
    }

    // Read File B
    fileB, err := os.ReadFile(fileBPath)
    if err != nil {
        return -1, fmt.Errorf("failed to read file B: %w", err)
    }

    // Get the length of File B
    lenB := len(fileB)

    // Iterate through File A with a window size equal to the length of File B
    for i := 0; i <= len(fileA)-lenB; i++ {
        if bytes.Equal(fileA[i:i+lenB], fileB) {
            return i, nil
        }
    }

    // If no match is found
    return -1, fmt.Errorf("child file not found in parent")
}