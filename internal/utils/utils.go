package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath replaces a leading '~' in the path with the user's home directory and returns the expanded path.
// Returns an error if the home directory cannot be determined.
func ExpandPath(path string) (string, error) {
	// Check if the path starts with ~
	if strings.HasPrefix(path, "~") {
		// Get the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to retrieve home directory: %w", err)
		}
		// Replace "~" with the home directory
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path, nil
}

// WriteToFile writes a byte slice to a file
func WriteToFile(filePath string, data []byte) error {
	var err error
	base := filepath.Dir(filePath)
	err = MkdirP(base)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	// Create or open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}

// MkdirP creates a directory with all necessary parents (like `mkdir -p`).
func MkdirP(path string) error {
	// os.MkdirAll creates the directory along with any necessary parents if they don't exist.
	// If the directory already exists, it does nothing and does not return an error.
	fmt.Printf("Creating directory path: %s\n", path)
	err := os.MkdirAll(path, os.ModePerm) // os.ModePerm sets permissions to 0777
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// FileExists checks if a file exists at the given path
func FileExists(path string) bool {
	// Use os.Stat to check if the file or directory exists
	info, err := os.Stat(path)

	// File does not exist if an error occurs and it is not a permission error
	if os.IsNotExist(err) {
		return false
	}

	// Check if the found path is actually a file and not a directory
	return !info.IsDir()
}

// CopySlice Helper function to copy a slice of any type
func CopySlice[T any](original []T) []T {
	// Create a new slice with the same length as the original
	copied := make([]T, len(original))
	copy(copied, original) // Copies the data from the original slice to the new slice
	return copied
}
