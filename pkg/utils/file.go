package utils

import (
	"os"
	"path/filepath"
)

func ExpandPath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

func FileExists(filePath string) bool {
	filePath, err := ExpandPath(filePath)
	if err != nil {
		return false
	}

	if _, err := os.Stat(filePath); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
