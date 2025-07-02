package files

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func createTestFiles() []*File {
	return []*File{
		{
			FileName: "file1.txt",
			Path:     "/test/path/file1.txt",
			Size:     1024,
			ModTime:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Hashes: map[string]string{
				"md5":    "d41d8cd98f00b204e9800998ecf8427e",
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
		{
			FileName: "file2.go",
			Path:     "/test/path/file2.go",
			Size:     2048,
			ModTime:  time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			Hashes: map[string]string{
				"md5":    "5d41402abc4b2a76b9719d911017c592",
				"sha1":   "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
				"sha256": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			},
		},
	}
}

func TestWriteOutput_StandardFormat(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test files
	testFiles := createTestFiles()

	// Write output
	err = WriteOutput(testFiles, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteOutput failed: %v", err)
	}

	// Read and verify output
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Check header
	expectedHeader := []string{"Path", "FileName", "Size", "Hash", "HashType"}
	if len(records) == 0 {
		t.Fatal("No records found in output file")
	}
	header := records[0]
	if len(header) != len(expectedHeader) {
		t.Fatalf("Expected %d header columns, got %d", len(expectedHeader), len(header))
	}
	for i, expected := range expectedHeader {
		if header[i] != expected {
			t.Errorf("Header column %d: expected %s, got %s", i, expected, header[i])
		}
	}

	// Check that we have the right number of data rows
	// Each file should have one row per hash type
	expectedRows := 2 + 3 // file1 has 2 hashes, file2 has 3 hashes
	actualRows := len(records) - 1 // subtract header
	if actualRows != expectedRows {
		t.Errorf("Expected %d data rows, got %d", expectedRows, actualRows)
	}
}

func TestWriteOutputCondensed(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test files
	testFiles := createTestFiles()

	// Write condensed output
	err = WriteOutputCondensed(testFiles, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteOutputCondensed failed: %v", err)
	}

	// Read and verify output
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Check we have header + 2 data rows (one per file)
	if len(records) != 3 {
		t.Fatalf("Expected 3 records (header + 2 files), got %d", len(records))
	}

	// Check header structure
	header := records[0]
	if header[0] != "Path" || header[1] != "FileName" || header[2] != "Size" {
		t.Error("Unexpected header format in condensed output")
	}

	// Check that header contains hash algorithm columns
	headerStr := strings.Join(header, ",")
	if !strings.Contains(headerStr, "md5") || !strings.Contains(headerStr, "sha256") {
		t.Error("Header should contain hash algorithm columns")
	}
}

func TestWriteOutputForIOC(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test files
	testFiles := createTestFiles()

	// Write IOC output
	err = WriteOutputForIOC(testFiles, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteOutputForIOC failed: %v", err)
	}

	// Read and verify output
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Check we have header + 2 data rows
	if len(records) != 3 {
		t.Fatalf("Expected 3 records (header + 2 files), got %d", len(records))
	}

	// Check IOC header format
	expectedHeader := []string{"file_path", "file_name", "file_size", "md5", "sha1", "sha256", "sha512"}
	header := records[0]
	if len(header) != len(expectedHeader) {
		t.Fatalf("Expected %d header columns, got %d", len(expectedHeader), len(header))
	}
	for i, expected := range expectedHeader {
		if header[i] != expected {
			t.Errorf("IOC header column %d: expected %s, got %s", i, expected, header[i])
		}
	}

	// Check data rows
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) != len(expectedHeader) {
			t.Errorf("Row %d: expected %d columns, got %d", i, len(expectedHeader), len(row))
		}

		// Check that file_path and file_name are not empty
		if row[0] == "" || row[1] == "" {
			t.Errorf("Row %d: file_path or file_name is empty", i)
		}

		// Check that file_size is numeric
		if row[2] == "" {
			t.Errorf("Row %d: file_size is empty", i)
		}
	}
}

func TestWriteOutput_EmptyFileList(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test with empty file list
	emptyFiles := []*File{}

	// Write output
	err = WriteOutput(emptyFiles, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteOutput failed with empty list: %v", err)
	}

	// Read and verify output
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Should have only header row
	if len(records) != 1 {
		t.Fatalf("Expected 1 record (header only), got %d", len(records))
	}
}

func TestWriteOutput_FilesWithoutHashes(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create files without hashes
	testFiles := []*File{
		{
			FileName: "nohash.txt",
			Path:     "/test/nohash.txt",
			Size:     512,
			ModTime:  time.Now(),
			Hashes:   make(map[string]string), // Empty hashes
		},
	}

	// Write output
	err = WriteOutput(testFiles, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteOutput failed: %v", err)
	}

	// Read and verify output
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// Should have header + 1 row with N/A values
	if len(records) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(records))
	}

	dataRow := records[1]
	if dataRow[3] != "N/A" || dataRow[4] != "N/A" {
		t.Error("Expected N/A values for missing hashes")
	}
}

func TestGetHashOrEmpty(t *testing.T) {
	hashes := map[string]string{
		"md5":    "test_md5_hash",
		"sha256": "test_sha256_hash",
	}

	// Test existing hash
	result := getHashOrEmpty(hashes, "md5")
	if result != "test_md5_hash" {
		t.Errorf("Expected 'test_md5_hash', got '%s'", result)
	}

	// Test missing hash
	result = getHashOrEmpty(hashes, "sha512")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}

	// Test with nil map
	result = getHashOrEmpty(nil, "md5")
	if result != "" {
		t.Errorf("Expected empty string with nil map, got '%s'", result)
	}
}

func TestWriteOutput_InvalidPath(t *testing.T) {
	// Test with invalid output path
	testFiles := createTestFiles()
	
	// Use invalid path (directory that doesn't exist)
	invalidPath := "/nonexistent/directory/output.csv"
	
	err := WriteOutput(testFiles, invalidPath)
	if err == nil {
		t.Error("Expected error for invalid output path, got nil")
	}
}

func TestWriteOutput_LargeFile(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dirhash_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create file with many hashes
	largeFile := &File{
		FileName: "large.txt",
		Path:     "/test/large.txt",
		Size:     1048576, // 1MB
		ModTime:  time.Now(),
		Hashes: map[string]string{
			"md5":    "large_md5_hash_value_here",
			"sha1":   "large_sha1_hash_value_here",
			"sha256": "large_sha256_hash_value_here",
			"sha512": "large_sha512_hash_value_here",
		},
	}

	testFiles := []*File{largeFile}

	// Test all output formats
	formats := []func([]*File, string) error{
		WriteOutput,
		WriteOutputCondensed,
		WriteOutputForIOC,
	}

	for i, writeFunc := range formats {
		testFile := filepath.Join(filepath.Dir(tmpFile.Name()), 
			"test_large_"+string(rune('0'+i))+".csv")
		defer os.Remove(testFile)

		err = writeFunc(testFiles, testFile)
		if err != nil {
			t.Errorf("Write function %d failed: %v", i, err)
		}

		// Verify file was created and has content
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("Output file %s was not created", testFile)
		}
	}
}