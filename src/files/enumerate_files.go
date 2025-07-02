package files

import (
	"log"
	"os"
	"path/filepath"
)

// EnumerateFiles enumerates all files in a directory and its subdirectories, and returns a slice of file pointers.
func EnumerateFiles(dir string) ([]*File, error) {
	var files []*File

	// Traverse the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			// Skip this file or directory and continue
			return nil
		}

		// Check if the path is a file (not a directory)
		if !info.IsDir() {
			file := NewFile(path, info.Name(), info)
			files = append(files, file)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
