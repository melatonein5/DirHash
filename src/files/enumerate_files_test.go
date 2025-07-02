package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnumerateFiles_SingleFile(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test EnumerateFiles
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	file := files[0]
	if file.FileName != "test.txt" {
		t.Errorf("Expected filename test.txt, got %s", file.FileName)
	}
	if file.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, file.Path)
	}
	if file.Size != 12 { // "test content" is 12 bytes
		t.Errorf("Expected size 12, got %d", file.Size)
	}
	if file.Hashes == nil {
		t.Error("Expected Hashes map to be initialized")
	}
	if len(file.Hashes) != 0 {
		t.Errorf("Expected empty Hashes map, got %d entries", len(file.Hashes))
	}
}

func TestEnumerateFiles_MultipleFiles(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.go"}
	for i, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		content := []byte("content " + string(rune('1'+i)))
		err = os.WriteFile(filePath, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test EnumerateFiles
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Fatalf("Expected %d files, got %d", len(testFiles), len(files))
	}

	// Check that all files are present (order may vary)
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[file.FileName] = true
		if file.Size <= 0 {
			t.Errorf("File %s should have positive size, got %d", file.FileName, file.Size)
		}
	}

	for _, expectedFile := range testFiles {
		if !foundFiles[expectedFile] {
			t.Errorf("Expected file %s not found", expectedFile)
		}
	}
}

func TestEnumerateFiles_WithSubdirectories(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files in root and subdirectory
	files := map[string]string{
		"root.txt":           "root content",
		"subdir/nested.txt":  "nested content",
		"subdir/another.go":  "go content",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, relPath)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", relPath, err)
		}
	}

	// Test EnumerateFiles
	result, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	if len(result) != len(files) {
		t.Fatalf("Expected %d files, got %d", len(files), len(result))
	}

	// Check that all files are found with correct paths
	foundPaths := make(map[string]bool)
	for _, file := range result {
		// Convert absolute path back to relative for comparison
		relPath, err := filepath.Rel(tmpDir, file.Path)
		if err != nil {
			t.Fatalf("Failed to get relative path: %v", err)
		}
		foundPaths[relPath] = true

		// Check content size matches
		expectedContent, exists := files[relPath]
		if !exists {
			t.Errorf("Unexpected file found: %s", relPath)
			continue
		}
		if file.Size != int64(len(expectedContent)) {
			t.Errorf("File %s: expected size %d, got %d", relPath, len(expectedContent), file.Size)
		}
	}

	for expectedPath := range files {
		if !foundPaths[expectedPath] {
			t.Errorf("Expected file %s not found", expectedPath)
		}
	}
}

func TestEnumerateFiles_EmptyDirectory(t *testing.T) {
	// Create empty temporary directory
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test EnumerateFiles on empty directory
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("Expected 0 files in empty directory, got %d", len(files))
	}
}

func TestEnumerateFiles_NonexistentDirectory(t *testing.T) {
	// Test EnumerateFiles on nonexistent directory
	// The function should return an error when the initial directory doesn't exist
	files, err := EnumerateFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
	if files != nil {
		t.Error("Expected nil files for nonexistent directory")
	}
}

func TestEnumerateFiles_DirectoryOnly(t *testing.T) {
	// Create temporary directory with only subdirectories (no files)
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories
	subDirs := []string{"dir1", "dir2", "dir1/subdir"}
	for _, dir := range subDirs {
		err = os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Test EnumerateFiles
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("Expected 0 files (directories should be ignored), got %d", len(files))
	}
}