package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashFiles_SingleAlgorithm(t *testing.T) {
	// Create temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file with known content
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "hello world"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Enumerate files
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test MD5 hashing
	hashedFiles, err := HashFiles(files, []int{0}) // MD5 = 0
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if len(hashedFiles) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(hashedFiles))
	}

	file := hashedFiles[0]
	if len(file.Hashes) != 1 {
		t.Fatalf("Expected 1 hash, got %d", len(file.Hashes))
	}

	md5Hash, exists := file.Hashes["md5"]
	if !exists {
		t.Fatal("MD5 hash not found")
	}

	// Known MD5 of "hello world" is 5eb63bbbe01eeed093cb22bb8f5acdc3
	expectedMD5 := "5eb63bbbe01eeed093cb22bb8f5acdc3"
	if md5Hash != expectedMD5 {
		t.Errorf("Expected MD5 %s, got %s", expectedMD5, md5Hash)
	}
}

func TestHashFiles_MultipleAlgorithms(t *testing.T) {
	// Create temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file with known content
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "hello world"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Enumerate files
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test multiple algorithms
	algorithms := []int{0, 1, 2, 3} // MD5, SHA1, SHA256, SHA512
	hashedFiles, err := HashFiles(files, algorithms)
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if len(hashedFiles) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(hashedFiles))
	}

	file := hashedFiles[0]
	if len(file.Hashes) != 4 {
		t.Fatalf("Expected 4 hashes, got %d", len(file.Hashes))
	}

	// Check that all expected algorithms are present
	expectedAlgorithms := []string{"md5", "sha1", "sha256", "sha512"}
	for _, algo := range expectedAlgorithms {
		if _, exists := file.Hashes[algo]; !exists {
			t.Errorf("Hash for algorithm %s not found", algo)
		}
	}

	// Verify known hash values for "hello world"
	expectedHashes := map[string]string{
		"md5":    "5eb63bbbe01eeed093cb22bb8f5acdc3",
		"sha1":   "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
		"sha256": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		"sha512": "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
	}

	for algo, expectedHash := range expectedHashes {
		if actualHash := file.Hashes[algo]; actualHash != expectedHash {
			t.Errorf("Algorithm %s: expected %s, got %s", algo, expectedHash, actualHash)
		}
	}
}

func TestHashFiles_MultipleFiles(t *testing.T) {
	// Create temporary directory with multiple test files
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with different content
	testFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Enumerate files
	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test hashing
	hashedFiles, err := HashFiles(files, []int{0, 2}) // MD5 and SHA256
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if len(hashedFiles) != len(testFiles) {
		t.Fatalf("Expected %d files, got %d", len(testFiles), len(hashedFiles))
	}

	// Check that each file has the expected hashes
	for _, file := range hashedFiles {
		if len(file.Hashes) != 2 {
			t.Errorf("File %s: expected 2 hashes, got %d", file.FileName, len(file.Hashes))
		}

		if _, exists := file.Hashes["md5"]; !exists {
			t.Errorf("File %s: MD5 hash not found", file.FileName)
		}

		if _, exists := file.Hashes["sha256"]; !exists {
			t.Errorf("File %s: SHA256 hash not found", file.FileName)
		}

		// Verify hashes are not empty
		if file.Hashes["md5"] == "" {
			t.Errorf("File %s: MD5 hash is empty", file.FileName)
		}
		if file.Hashes["sha256"] == "" {
			t.Errorf("File %s: SHA256 hash is empty", file.FileName)
		}
	}
}

func TestHashFiles_EmptyFileList(t *testing.T) {
	// Test with empty file list
	emptyFiles := []*File{}
	hashedFiles, err := HashFiles(emptyFiles, []int{0})
	if err != nil {
		t.Fatalf("HashFiles failed with empty list: %v", err)
	}

	if len(hashedFiles) != 0 {
		t.Fatalf("Expected 0 files, got %d", len(hashedFiles))
	}
}

func TestHashFiles_InvalidAlgorithm(t *testing.T) {
	// Create temporary directory with test file
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test with invalid algorithm ID
	_, err = HashFiles(files, []int{999})
	if err == nil {
		t.Error("Expected error for invalid algorithm, got nil")
	}
}

func TestHashFiles_NoValidAlgorithms(t *testing.T) {
	// Create temporary directory with test file
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test with empty algorithm list
	_, err = HashFiles(files, []int{})
	if err == nil {
		t.Error("Expected error for empty algorithm list, got nil")
	}
}

func TestHashFiles_NonexistentFile(t *testing.T) {
	// Create a File struct pointing to a nonexistent file
	nonexistentFile := &File{
		FileName: "nonexistent.txt",
		Path:     "/nonexistent/path/file.txt",
		Size:     0,
		Hashes:   make(map[string]string),
	}

	files := []*File{nonexistentFile}
	_, err := HashFiles(files, []int{0})
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestCalculateAllHashes_EmptyFile(t *testing.T) {
	// Create temporary empty file
	tmpDir, err := os.MkdirTemp("", "dirhash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(testFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	files, err := EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnumerateFiles failed: %v", err)
	}

	// Test hashing empty file
	hashedFiles, err := HashFiles(files, []int{0}) // MD5
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if len(hashedFiles) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(hashedFiles))
	}

	file := hashedFiles[0]
	md5Hash, exists := file.Hashes["md5"]
	if !exists {
		t.Fatal("MD5 hash not found")
	}

	// MD5 of empty string is d41d8cd98f00b204e9800998ecf8427e
	expectedMD5 := "d41d8cd98f00b204e9800998ecf8427e"
	if md5Hash != expectedMD5 {
		t.Errorf("Expected MD5 %s, got %s", expectedMD5, md5Hash)
	}
}
