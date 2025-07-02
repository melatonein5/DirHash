package files

import (
	"os"
	"testing"
	"time"
)

func TestNewFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some test data
	testData := "test file content"
	_, err = tmpFile.WriteString(testData)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Get file info
	fileInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Test NewFile
	file := NewFile(tmpFile.Name(), fileInfo.Name(), fileInfo)

	if file.Path != tmpFile.Name() {
		t.Errorf("Expected path %s, got %s", tmpFile.Name(), file.Path)
	}
	if file.FileName != fileInfo.Name() {
		t.Errorf("Expected filename %s, got %s", fileInfo.Name(), file.FileName)
	}
	if file.Size != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), file.Size)
	}
	if file.Hashes == nil {
		t.Error("Expected Hashes map to be initialized")
	}
	if len(file.Hashes) != 0 {
		t.Errorf("Expected empty Hashes map, got %d entries", len(file.Hashes))
	}
	if file.ModTime.IsZero() {
		t.Error("Expected ModTime to be set")
	}
}

func TestGetSupportedAlgorithms(t *testing.T) {
	algorithms := GetSupportedAlgorithms()

	expectedAlgorithms := []HashAlgorithm{
		{ID: 0, Name: "md5"},
		{ID: 1, Name: "sha1"},
		{ID: 2, Name: "sha256"},
		{ID: 3, Name: "sha512"},
	}

	if len(algorithms) != len(expectedAlgorithms) {
		t.Fatalf("Expected %d algorithms, got %d", len(expectedAlgorithms), len(algorithms))
	}

	for i, expected := range expectedAlgorithms {
		if algorithms[i].ID != expected.ID {
			t.Errorf("Algorithm %d: expected ID %d, got %d", i, expected.ID, algorithms[i].ID)
		}
		if algorithms[i].Name != expected.Name {
			t.Errorf("Algorithm %d: expected name %s, got %s", i, expected.Name, algorithms[i].Name)
		}
	}
}

func TestFileStructFields(t *testing.T) {
	// Test that File struct has all expected fields
	file := &File{
		FileName: "test.txt",
		Path:     "/path/to/test.txt",
		Size:     1024,
		ModTime:  time.Now(),
		Hashes:   make(map[string]string),
	}

	if file.FileName != "test.txt" {
		t.Errorf("Expected FileName test.txt, got %s", file.FileName)
	}
	if file.Path != "/path/to/test.txt" {
		t.Errorf("Expected Path /path/to/test.txt, got %s", file.Path)
	}
	if file.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", file.Size)
	}
	if file.Hashes == nil {
		t.Error("Expected Hashes to be initialized")
	}
}

func TestFileHashesMap(t *testing.T) {
	file := &File{
		FileName: "test.txt",
		Path:     "/path/to/test.txt",
		Size:     1024,
		ModTime:  time.Now(),
		Hashes:   make(map[string]string),
	}

	// Test adding hashes
	file.Hashes["md5"] = "d41d8cd98f00b204e9800998ecf8427e"
	file.Hashes["sha256"] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	if len(file.Hashes) != 2 {
		t.Errorf("Expected 2 hashes, got %d", len(file.Hashes))
	}

	expectedMd5 := "d41d8cd98f00b204e9800998ecf8427e"
	if file.Hashes["md5"] != expectedMd5 {
		t.Errorf("Expected MD5 %s, got %s", expectedMd5, file.Hashes["md5"])
	}

	expectedSha256 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if file.Hashes["sha256"] != expectedSha256 {
		t.Errorf("Expected SHA256 %s, got %s", expectedSha256, file.Hashes["sha256"])
	}
}

func TestHashAlgorithmStruct(t *testing.T) {
	algo := HashAlgorithm{
		ID:   2,
		Name: "sha256",
	}

	if algo.ID != 2 {
		t.Errorf("Expected ID 2, got %d", algo.ID)
	}
	if algo.Name != "sha256" {
		t.Errorf("Expected Name sha256, got %s", algo.Name)
	}
}
