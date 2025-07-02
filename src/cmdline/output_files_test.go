package cmdline

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/melatonein5/DirHash/src/files"
)

func createTestFileList() []*files.File {
	return []*files.File{
		{
			FileName: "test1.txt",
			Path:     "/path/to/test1.txt",
			Size:     1024,
			ModTime:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Hashes: map[string]string{
				"md5":    "d41d8cd98f00b204e9800998ecf8427e",
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
		{
			FileName: "test2.go",
			Path:     "/path/to/test2.go",
			Size:     2048,
			ModTime:  time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			Hashes: map[string]string{
				"md5":    "5d41402abc4b2a76b9719d911017c592",
				"sha1":   "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed",
				"sha256": "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				"sha512": "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
			},
		},
	}
}

func captureOutput(fn func()) string {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run function
	fn()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestOutputFiles_StandardFormat(t *testing.T) {
	testFiles := createTestFileList()

	output := captureOutput(func() {
		OutputFiles(testFiles)
	})

	// Check that output contains expected elements
	if !strings.Contains(output, "File Name") {
		t.Error("Output should contain 'File Name' header")
	}
	if !strings.Contains(output, "Hash Type") {
		t.Error("Output should contain 'Hash Type' header")
	}
	if !strings.Contains(output, "test1.txt") {
		t.Error("Output should contain test1.txt")
	}
	if !strings.Contains(output, "test2.go") {
		t.Error("Output should contain test2.go")
	}

	// Should have separate rows for each hash type
	lines := strings.Split(output, "\n")
	contentLines := 0
	for _, line := range lines {
		if strings.Contains(line, "test1.txt") || strings.Contains(line, "test2.go") {
			contentLines++
		}
	}
	
	// test1.txt has 2 hashes, test2.go has 4 hashes = 6 total rows
	if contentLines != 6 {
		t.Errorf("Expected 6 content lines (2+4 hashes), got %d", contentLines)
	}
}

func TestOutputFilesCondensed(t *testing.T) {
	testFiles := createTestFileList()

	output := captureOutput(func() {
		OutputFilesCondensed(testFiles)
	})

	// Check headers
	if !strings.Contains(output, "File Name") {
		t.Error("Output should contain 'File Name' header")
	}
	if !strings.Contains(output, "Hashes") {
		t.Error("Output should contain 'Hashes' header")
	}

	// Check file names appear
	if !strings.Contains(output, "test1.txt") {
		t.Error("Output should contain test1.txt")
	}
	if !strings.Contains(output, "test2.go") {
		t.Error("Output should contain test2.go")
	}

	// Should contain hash format "algo:hash"
	if !strings.Contains(output, "md5:") {
		t.Error("Output should contain 'md5:' format")
	}
	if !strings.Contains(output, "sha256:") {
		t.Error("Output should contain 'sha256:' format")
	}

	// Should have pipe separator
	if !strings.Contains(output, " | ") {
		t.Error("Output should contain ' | ' separator")
	}

	// Should have exactly 2 file rows (one per file)
	lines := strings.Split(output, "\n")
	fileLines := 0
	for _, line := range lines {
		if strings.Contains(line, "test1.txt") || strings.Contains(line, "test2.go") {
			fileLines++
		}
	}
	
	if fileLines != 2 {
		t.Errorf("Expected 2 file lines in condensed format, got %d", fileLines)
	}
}

func TestOutputFilesIOC(t *testing.T) {
	testFiles := createTestFileList()

	output := captureOutput(func() {
		OutputFilesIOC(testFiles)
	})

	// Check IOC headers
	expectedHeaders := []string{"File Path", "File Name", "Size", "MD5", "SHA1", "SHA256", "SHA512"}
	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("Output should contain '%s' header", header)
		}
	}

	// Check file data appears
	if !strings.Contains(output, "test1.txt") {
		t.Error("Output should contain test1.txt")
	}
	if !strings.Contains(output, "test2.go") {
		t.Error("Output should contain test2.go")
	}

	// Check that hash values appear in output
	if !strings.Contains(output, "d41d8cd98f00b204e9800998ecf8427e") {
		t.Error("Output should contain MD5 hash from test file")
	}

	// Should have exactly 2 file rows
	lines := strings.Split(output, "\n")
	fileLines := 0
	for _, line := range lines {
		if strings.Contains(line, "test1.txt") || strings.Contains(line, "test2.go") {
			fileLines++
		}
	}
	
	if fileLines != 2 {
		t.Errorf("Expected 2 file lines in IOC format, got %d", fileLines)
	}
}

func TestOutputFiles_EmptyList(t *testing.T) {
	emptyFiles := []*files.File{}

	output := captureOutput(func() {
		OutputFiles(emptyFiles)
	})

	if !strings.Contains(output, "No files to display") {
		t.Error("Output should contain 'No files to display' message for empty list")
	}
}

func TestOutputFilesCondensed_EmptyList(t *testing.T) {
	emptyFiles := []*files.File{}

	output := captureOutput(func() {
		OutputFilesCondensed(emptyFiles)
	})

	if !strings.Contains(output, "No files to display") {
		t.Error("Condensed output should contain 'No files to display' message for empty list")
	}
}

func TestOutputFilesIOC_EmptyList(t *testing.T) {
	emptyFiles := []*files.File{}

	output := captureOutput(func() {
		OutputFilesIOC(emptyFiles)
	})

	if !strings.Contains(output, "No files to display") {
		t.Error("IOC output should contain 'No files to display' message for empty list")
	}
}

func TestGetHashOrNA(t *testing.T) {
	hashes := map[string]string{
		"md5":    "test_md5_hash",
		"sha256": "test_sha256_hash",
	}

	// Test existing hash
	result := getHashOrNA(hashes, "md5")
	if result != "test_md5_hash" {
		t.Errorf("Expected 'test_md5_hash', got '%s'", result)
	}

	// Test missing hash
	result = getHashOrNA(hashes, "sha512")
	if result != "N/A" {
		t.Errorf("Expected 'N/A', got '%s'", result)
	}

	// Test with nil map
	result = getHashOrNA(nil, "md5")
	if result != "N/A" {
		t.Errorf("Expected 'N/A' with nil map, got '%s'", result)
	}
}

func TestOutputFilesIOC_MissingHashes(t *testing.T) {
	// Create file with only some hashes
	testFiles := []*files.File{
		{
			FileName: "partial.txt",
			Path:     "/path/to/partial.txt",
			Size:     512,
			ModTime:  time.Now(),
			Hashes: map[string]string{
				"md5": "d41d8cd98f00b204e9800998ecf8427e",
				// Missing sha1, sha256, sha512
			},
		},
	}

	output := captureOutput(func() {
		OutputFilesIOC(testFiles)
	})

	// Should contain the available hash
	if !strings.Contains(output, "d41d8cd98f00b204e9800998ecf8427e") {
		t.Error("Output should contain available MD5 hash")
	}

	// Should contain N/A for missing hashes
	naCount := strings.Count(output, "N/A")
	if naCount < 3 { // Should have N/A for sha1, sha256, sha512
		t.Errorf("Expected at least 3 'N/A' entries for missing hashes, got %d", naCount)
	}
}

func TestOutputFiles_NoHashes(t *testing.T) {
	// Create file with no hashes
	testFiles := []*files.File{
		{
			FileName: "nohash.txt",
			Path:     "/path/to/nohash.txt",
			Size:     256,
			ModTime:  time.Now(),
			Hashes:   make(map[string]string), // Empty map
		},
	}

	output := captureOutput(func() {
		OutputFiles(testFiles)
	})

	// Should contain filename and path
	if !strings.Contains(output, "nohash.txt") {
		t.Error("Output should contain filename even without hashes")
	}

	// Should contain N/A for hash and hash type
	if !strings.Contains(output, "N/A") {
		t.Error("Output should contain 'N/A' for missing hash information")
	}
}

func TestOutputFilesCondensed_SingleHash(t *testing.T) {
	// Create file with single hash
	testFiles := []*files.File{
		{
			FileName: "single.txt",
			Path:     "/path/to/single.txt",
			Size:     128,
			ModTime:  time.Now(),
			Hashes: map[string]string{
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
	}

	output := captureOutput(func() {
		OutputFilesCondensed(testFiles)
	})

	// Should contain the hash in correct format
	if !strings.Contains(output, "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855") {
		t.Error("Output should contain hash in 'algo:hash' format")
	}

	// Should not contain pipe separator for single hash
	lineWithFile := ""
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "single.txt") {
			lineWithFile = line
			break
		}
	}
	
	if strings.Contains(lineWithFile, " | ") {
		t.Error("Single hash output should not contain pipe separator")
	}
}