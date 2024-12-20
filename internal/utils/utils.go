package utils

import (
	"fmt"
	"os"
	"strings"
)

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
