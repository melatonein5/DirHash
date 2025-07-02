package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/melatonein5/DirHash/src/args"
	"github.com/melatonein5/DirHash/src/files"
)

// TestMainLogic tests the main application logic flow
func TestMainLogic(t *testing.T) {
	// Create test environment
	tmpDir, err := os.MkdirTemp("", "dirhash_main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := map[string]string{
		"sample1.txt": "hello world",
		"sample2.exe": "binary content", 
		"document.pdf": "pdf content here",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test the main logic flow
	testArgs := []string{"-i", tmpDir, "-a", "md5", "sha256", "-f", "ioc"}
	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Step 1: Enumerate files (main.go:34)
	fs, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	if len(fs) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(fs))
	}

	// Step 2: Hash files (main.go:42)  
	hashedFiles, err := files.HashFiles(fs, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	if len(hashedFiles) != len(testFiles) {
		t.Errorf("Expected %d hashed files, got %d", len(testFiles), len(hashedFiles))
	}

	// Verify hash completeness
	for _, file := range hashedFiles {
		if len(file.Hashes) != 2 { // md5 and sha256
			t.Errorf("File %s: expected 2 hashes, got %d", file.FileName, len(file.Hashes))
		}
		if _, exists := file.Hashes["md5"]; !exists {
			t.Errorf("File %s missing MD5 hash", file.FileName)
		}
		if _, exists := file.Hashes["sha256"]; !exists {
			t.Errorf("File %s missing SHA256 hash", file.FileName)
		}
	}

	// Step 3: Test output logic (main.go:61-77)
	outputFile := filepath.Join(tmpDir, "output.csv")
	
	// Test format selection switch (main.go:64-71)
	var writeErr error
	switch parsedArgs.OutputFormat {
	case "condensed":
		writeErr = files.WriteOutputCondensed(hashedFiles, outputFile)
	case "ioc":
		writeErr = files.WriteOutputForIOC(hashedFiles, outputFile)
	default: // "standard"
		writeErr = files.WriteOutput(hashedFiles, outputFile)
	}

	if writeErr != nil {
		t.Fatalf("Error writing output file: %v", writeErr)
	}

	// Verify output file
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file should have been created")
	}

	// Verify output content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputContent := string(content)
	
	// Check IOC format headers (lowercase with underscores)
	expectedHeaders := []string{"file_path", "file_name", "file_size", "md5", "sha256"}
	for _, header := range expectedHeaders {
		if !strings.Contains(outputContent, header) {
			t.Errorf("Output should contain header '%s'", header)
		}
	}

	// Check all test files are present
	for filename := range testFiles {
		if !strings.Contains(outputContent, filename) {
			t.Errorf("Output should contain filename '%s'", filename)
		}
	}
}

// TestMainErrorPaths tests error handling in main
func TestMainErrorPaths(t *testing.T) {
	// Test file enumeration error (main.go:35-37)
	_, err := files.EnumerateFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}

	// Test file write error (main.go:73-75)
	tmpDir, err := os.MkdirTemp("", "dirhash_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fs, err := files.EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(fs, []int{0})
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Test invalid output path
	err = files.WriteOutput(hashedFiles, "/invalid/path/output.csv")
	if err == nil {
		t.Error("Expected error for invalid output path")
	}
}

// TestMainOutputFormats tests all output format branches
func TestMainOutputFormats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_format_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.exe")
	err = os.WriteFile(testFile, []byte("executable"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Process file
	fs, err := files.EnumerateFiles(tmpDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(fs, []int{0, 2}) // MD5, SHA256
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Test all format switch cases (main.go:64-71)
	formats := []string{"standard", "condensed", "ioc"}
	
	for _, format := range formats {
		outputFile := filepath.Join(tmpDir, "output_"+format+".csv")
		
		// Test the exact switch logic from main()
		var writeErr error
		switch format {
		case "condensed":
			writeErr = files.WriteOutputCondensed(hashedFiles, outputFile)
		case "ioc":
			writeErr = files.WriteOutputForIOC(hashedFiles, outputFile)
		default: // "standard"
			writeErr = files.WriteOutput(hashedFiles, outputFile)
		}

		if writeErr != nil {
			t.Errorf("Failed to write %s format: %v", format, writeErr)
		}

		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Errorf("Output file not created for %s format", format)
		}
	}
}

// TestMainSecurityWorkflow tests complete security use case
func TestMainSecurityWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_security")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "suspicious")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create security-relevant test files
	securityFiles := map[string]string{
		"malware.exe":           "suspicious executable",
		"document.pdf":          "normal document", 
		"suspicious/trojan.dll": "malicious library",
	}

	for relPath, content := range securityFiles {
		fullPath := filepath.Join(tmpDir, relPath)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", relPath, err)
		}
	}

	// Test complete security workflow
	securityArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha1", "sha256",
		"-f", "ioc",
	}

	parsedArgs, err := args.ParseArgs(securityArgs)
	if err != nil {
		t.Fatalf("Failed to parse security args: %v", err)
	}

	// Execute main workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Enumeration failed: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	// Generate IOC output
	iocFile := filepath.Join(tmpDir, "iocs.csv")
	err = files.WriteOutputForIOC(hashedFiles, iocFile)
	if err != nil {
		t.Fatalf("IOC output failed: %v", err)
	}

	// Verify results
	if len(hashedFiles) != len(securityFiles) {
		t.Errorf("Expected %d files, got %d", len(securityFiles), len(hashedFiles))
	}

	// Verify all files have complete hashes
	for _, file := range hashedFiles {
		if len(file.Hashes) != 3 { // md5, sha1, sha256
			t.Errorf("File %s: expected 3 hashes, got %d", file.FileName, len(file.Hashes))
		}
	}

	// Verify IOC file content
	iocContent, err := os.ReadFile(iocFile)
	if err != nil {
		t.Fatalf("Failed to read IOC file: %v", err)
	}

	iocString := string(iocContent)
	expectedFiles := []string{"malware.exe", "document.pdf", "trojan.dll"}
	for _, expectedFile := range expectedFiles {
		if !strings.Contains(iocString, expectedFile) {
			t.Errorf("IOC output should contain %s", expectedFile)
		}
	}
}