package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/melatonein5/DirHash/src/args"
	"github.com/melatonein5/DirHash/src/files"
	"github.com/melatonein5/DirHash/src/kql"
	"github.com/melatonein5/DirHash/src/yara"
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
		"sample1.txt":  "hello world",
		"sample2.exe":  "binary content",
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

// TestMainYaraIntegration tests YARA rule generation functionality
func TestMainYaraIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_yara_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := map[string]string{
		"malware.exe": "malicious content here",
		"trojan.dll":  "another malicious file",
		"spyware.bin": "spyware payload",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test YARA generation with all hash types
	testArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha256", "sha512",
		"-y", filepath.Join(tmpDir, "malware.yar"),
		"--yara-rule-name", "malware_detection",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Execute main workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Test YARA rule generation (mimicking main.go logic)
	var rule *yara.YaraRule
	ruleName := parsedArgs.YaraRuleName
	if ruleName == "" {
		ruleName = "dirhash_generated_rule"
	}

	if parsedArgs.YaraHashOnly {
		hashTypes := append([]string{}, parsedArgs.StrHashAlgorithms...)
		rule, err = yara.GenerateYaraRuleFromHashes(hashedFiles, ruleName, hashTypes)
	} else {
		rule, err = yara.GenerateYaraRule(hashedFiles, ruleName)
	}

	if err != nil {
		t.Fatalf("Failed to generate YARA rule: %v", err)
	}

	// Write YARA rule
	yaraContent := rule.ToYaraFormat()
	err = os.WriteFile(parsedArgs.YaraFile, []byte(yaraContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write YARA file: %v", err)
	}

	// Verify YARA rule was created
	if _, err := os.Stat(parsedArgs.YaraFile); os.IsNotExist(err) {
		t.Error("YARA rule file should have been created")
	}

	// Read and verify YARA content
	yaraFileContent, err := os.ReadFile(parsedArgs.YaraFile)
	if err != nil {
		t.Fatalf("Failed to read YARA file: %v", err)
	}

	yaraString := string(yaraFileContent)

	// Verify YARA rule structure
	expectedElements := []string{
		"rule malware_detection",
		"meta:",
		"author = \"DirHash\"",
		"strings:",
		"condition:",
	}

	for _, element := range expectedElements {
		if !strings.Contains(yaraString, element) {
			t.Errorf("YARA rule should contain '%s'", element)
		}
	}

	// Verify hash strings are present
	if !strings.Contains(yaraString, "md5") {
		t.Error("YARA rule should contain MD5 hash strings")
	}
	if !strings.Contains(yaraString, "sha256") {
		t.Error("YARA rule should contain SHA256 hash strings")
	}
	if !strings.Contains(yaraString, "sha512") {
		t.Error("YARA rule should contain SHA512 hash strings")
	}

	// Verify filename strings are present
	for filename := range testFiles {
		if !strings.Contains(yaraString, filename) {
			t.Errorf("YARA rule should contain filename '%s'", filename)
		}
	}
}

// TestMainYaraHashOnlyMode tests YARA hash-only mode
func TestMainYaraHashOnlyMode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_yara_hash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "suspicious.exe")
	err = os.WriteFile(testFile, []byte("suspicious content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test YARA hash-only generation
	testArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha256",
		"-y", filepath.Join(tmpDir, "hashes.yar"),
		"--yara-rule-name", "hash_detection",
		"--yara-hash-only",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Execute workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Generate hash-only YARA rule
	hashTypes := append([]string{}, parsedArgs.StrHashAlgorithms...)
	rule, err := yara.GenerateYaraRuleFromHashes(hashedFiles, parsedArgs.YaraRuleName, hashTypes)
	if err != nil {
		t.Fatalf("Failed to generate hash-only YARA rule: %v", err)
	}

	// Write and verify
	yaraContent := rule.ToYaraFormat()
	err = os.WriteFile(parsedArgs.YaraFile, []byte(yaraContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write YARA file: %v", err)
	}

	// Read and verify content
	yaraFileContent, err := os.ReadFile(parsedArgs.YaraFile)
	if err != nil {
		t.Fatalf("Failed to read YARA file: %v", err)
	}

	yaraString := string(yaraFileContent)

	// Verify hash-only content
	if !strings.Contains(yaraString, "rule hash_detection") {
		t.Error("YARA rule should contain correct rule name")
	}

	if !strings.Contains(yaraString, "md5") {
		t.Error("Hash-only rule should contain MD5 hashes")
	}
	if !strings.Contains(yaraString, "sha256") {
		t.Error("Hash-only rule should contain SHA256 hashes")
	}

	// Verify NO filename strings in hash-only mode
	if strings.Contains(yaraString, "suspicious.exe") {
		t.Error("Hash-only rule should NOT contain filename strings")
	}
	if strings.Contains(yaraString, "$filename") {
		t.Error("Hash-only rule should NOT contain filename variables")
	}
}

// TestMainYaraErrorHandling tests YARA error scenarios
func TestMainYaraErrorHandling(t *testing.T) {
	// Test empty files list
	_, err := yara.GenerateYaraRule([]*files.File{}, "test")
	if err == nil {
		t.Error("Should return error for empty files list")
	}

	// Test invalid hash types
	testFile := &files.File{
		FileName: "test.exe",
		Hashes:   map[string]string{"md5": "somehash"},
	}
	_, err = yara.GenerateYaraRuleFromHashes([]*files.File{testFile}, "test", []string{"nonexistent"})
	if err == nil {
		t.Error("Should return error for invalid hash types")
	}
}

// TestMainKQLIntegration tests KQL query generation functionality
func TestMainKQLIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_kql_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := map[string]string{
		"malware.exe": "malicious content here",
		"trojan.dll":  "another malicious file",
		"spyware.bin": "spyware payload",
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test KQL generation with all hash types
	testArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha256", "sha512",
		"-q", filepath.Join(tmpDir, "malware.kql"),
		"--kql-name", "malware_detection",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Execute main workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Test KQL query generation (mimicking main.go logic)
	var query *kql.KQLQuery
	queryName := parsedArgs.KQLName
	if queryName == "" {
		queryName = "dirhash_generated_query"
	}

	options := kql.DefaultKQLQueryOptions()
	options.Tables = parsedArgs.KQLTables
	options.IncludeHashes = true
	options.IncludeFilenames = !parsedArgs.KQLHashOnly

	if parsedArgs.KQLHashOnly {
		query, err = kql.GenerateKQLQueryHashOnly(hashedFiles, queryName, parsedArgs.StrHashAlgorithms)
	} else {
		query, err = kql.GenerateKQLQueryWithOptions(hashedFiles, queryName, parsedArgs.StrHashAlgorithms, options)
	}

	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Write KQL query
	kqlContent := query.ToKQLFormat()
	err = os.WriteFile(parsedArgs.KQLFile, []byte(kqlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write KQL file: %v", err)
	}

	// Verify KQL query was created
	if _, err := os.Stat(parsedArgs.KQLFile); os.IsNotExist(err) {
		t.Error("KQL query file should have been created")
	}

	// Read and verify KQL content
	kqlFileContent, err := os.ReadFile(parsedArgs.KQLFile)
	if err != nil {
		t.Fatalf("Failed to read KQL file: %v", err)
	}

	kqlString := string(kqlFileContent)

	// Verify KQL query structure
	expectedElements := []string{
		"// KQL Query: malware_detection",
		"// Author: DirHash",
		"DeviceFileEvents",
		"TimeGenerated",
		"sort by TimeGenerated desc",
		"take 1000",
	}

	for _, element := range expectedElements {
		if !strings.Contains(kqlString, element) {
			t.Errorf("KQL query should contain '%s'", element)
		}
	}

	// Verify hash strings are present
	if !strings.Contains(kqlString, "MD5") {
		t.Error("KQL query should contain MD5 hash references")
	}
	if !strings.Contains(kqlString, "SHA256") {
		t.Error("KQL query should contain SHA256 hash references")
	}

	// Verify filename strings are present
	for filename := range testFiles {
		if !strings.Contains(kqlString, filename) {
			t.Errorf("KQL query should contain filename '%s'", filename)
		}
	}
}

// TestMainKQLHashOnlyMode tests KQL hash-only mode
func TestMainKQLHashOnlyMode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_kql_hash_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "suspicious.exe")
	err = os.WriteFile(testFile, []byte("suspicious content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test KQL hash-only generation
	testArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha256",
		"-q", filepath.Join(tmpDir, "hashes.kql"),
		"--kql-name", "hash_detection",
		"--kql-hash-only",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Execute workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Generate hash-only KQL query
	query, err := kql.GenerateKQLQueryHashOnly(hashedFiles, parsedArgs.KQLName, parsedArgs.StrHashAlgorithms)
	if err != nil {
		t.Fatalf("Failed to generate hash-only KQL query: %v", err)
	}

	// Write and verify
	kqlContent := query.ToKQLFormat()
	err = os.WriteFile(parsedArgs.KQLFile, []byte(kqlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write KQL file: %v", err)
	}

	// Read and verify content
	kqlFileContent, err := os.ReadFile(parsedArgs.KQLFile)
	if err != nil {
		t.Fatalf("Failed to read KQL file: %v", err)
	}

	kqlString := string(kqlFileContent)

	// Verify hash-only content
	if !strings.Contains(kqlString, "// KQL Query: hash_detection") {
		t.Error("KQL query should contain correct query name")
	}

	if !strings.Contains(kqlString, "MD5") {
		t.Error("Hash-only query should contain MD5 hash references")
	}
	if !strings.Contains(kqlString, "SHA256") {
		t.Error("Hash-only query should contain SHA256 hash references")
	}

	// Verify NO filename strings in hash-only mode
	if strings.Contains(kqlString, "suspicious.exe") {
		t.Error("Hash-only query should NOT contain filename strings")
	}
	if strings.Contains(kqlString, "FileName") && strings.Contains(kqlString, "suspicious.exe") {
		t.Error("Hash-only query should NOT contain filename-based conditions")
	}
}

// TestMainKQLMultipleTables tests KQL generation with multiple tables
func TestMainKQLMultipleTables(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dirhash_kql_tables_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "security.exe")
	err = os.WriteFile(testFile, []byte("security test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test KQL generation with multiple tables
	testArgs := []string{
		"-i", tmpDir,
		"-a", "md5", "sha256",
		"-q", filepath.Join(tmpDir, "security.kql"),
		"--kql-name", "security_detection",
		"--kql-tables", "DeviceFileEvents", "SecurityEvents", "CommonSecurityLog",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse args: %v", err)
	}

	// Verify KQL tables are parsed correctly
	expectedTables := []string{"DeviceFileEvents", "SecurityEvents", "CommonSecurityLog"}
	if len(parsedArgs.KQLTables) != len(expectedTables) {
		t.Fatalf("Expected %d KQL tables, got %d", len(expectedTables), len(parsedArgs.KQLTables))
	}

	for i, expected := range expectedTables {
		if parsedArgs.KQLTables[i] != expected {
			t.Errorf("Expected KQL table '%s', got '%s'", expected, parsedArgs.KQLTables[i])
		}
	}

	// Execute workflow
	enumFiles, err := files.EnumerateFiles(parsedArgs.StrInputDir)
	if err != nil {
		t.Fatalf("Error enumerating files: %v", err)
	}

	hashedFiles, err := files.HashFiles(enumFiles, parsedArgs.HashAlgorithmId)
	if err != nil {
		t.Fatalf("Error hashing files: %v", err)
	}

	// Generate KQL query with multiple tables
	options := kql.DefaultKQLQueryOptions()
	options.Tables = parsedArgs.KQLTables

	query, err := kql.GenerateKQLQueryWithOptions(hashedFiles, parsedArgs.KQLName, parsedArgs.StrHashAlgorithms, options)
	if err != nil {
		t.Fatalf("Failed to generate multi-table KQL query: %v", err)
	}

	// Write and verify
	kqlContent := query.ToKQLFormat()
	err = os.WriteFile(parsedArgs.KQLFile, []byte(kqlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write KQL file: %v", err)
	}

	// Read and verify content
	kqlFileContent, err := os.ReadFile(parsedArgs.KQLFile)
	if err != nil {
		t.Fatalf("Failed to read KQL file: %v", err)
	}

	kqlString := string(kqlFileContent)

	// Verify all tables are present
	for _, table := range expectedTables {
		if !strings.Contains(kqlString, table) {
			t.Errorf("KQL query should contain table '%s'", table)
		}
	}

	// Verify union is used for multiple tables
	if !strings.Contains(kqlString, "union") {
		t.Error("Multi-table KQL query should contain union statement")
	}

	// Verify table-specific field handling
	if !strings.Contains(kqlString, "MD5") {
		t.Error("DeviceFileEvents query should contain MD5 field")
	}

	if !strings.Contains(kqlString, "FileHash") {
		t.Error("SecurityEvents/CommonSecurityLog query should contain FileHash field")
	}
}

// TestMainKQLErrorHandling tests KQL error scenarios
func TestMainKQLErrorHandling(t *testing.T) {
	// Test empty files list
	_, err := kql.GenerateKQLQuery([]*files.File{}, "test", []string{"md5"})
	if err == nil {
		t.Error("Should return error for empty files list")
	}

	// Test files without hashes
	testFile := &files.File{
		FileName: "test.exe",
		Hashes:   map[string]string{},
	}
	query, err := kql.GenerateKQLQuery([]*files.File{testFile}, "test", []string{"md5"})
	if err != nil {
		t.Fatalf("Should not error for files without hashes: %v", err)
	}

	// Should still generate query with filename
	if len(query.FilenameList) == 0 {
		t.Error("Query should contain filename even without hashes")
	}
}

// TestMainKQLArgumentParsing tests KQL-specific argument parsing
func TestMainKQLArgumentParsing(t *testing.T) {
	// Test basic KQL arguments
	testArgs := []string{
		"-i", "/test/dir",
		"-q", "/output/query.kql",
		"--kql-name", "test_query",
	}

	parsedArgs, err := args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse KQL args: %v", err)
	}

	if !parsedArgs.KQLOutput {
		t.Error("KQLOutput should be true when -q flag is provided")
	}

	if parsedArgs.KQLFile != "/output/query.kql" {
		t.Errorf("Expected KQL file '/output/query.kql', got '%s'", parsedArgs.KQLFile)
	}

	if parsedArgs.KQLName != "test_query" {
		t.Errorf("Expected KQL name 'test_query', got '%s'", parsedArgs.KQLName)
	}

	// Test KQL hash-only flag
	testArgs = []string{
		"-i", "/test/dir",
		"-q", "/output/query.kql",
		"--kql-hash-only",
	}

	parsedArgs, err = args.ParseArgs(testArgs)
	if err != nil {
		t.Fatalf("Failed to parse KQL hash-only args: %v", err)
	}

	if !parsedArgs.KQLHashOnly {
		t.Error("KQLHashOnly should be true when --kql-hash-only flag is provided")
	}

	// Test default KQL tables
	if len(parsedArgs.KQLTables) != 1 || parsedArgs.KQLTables[0] != "DeviceFileEvents" {
		t.Errorf("Default KQL tables should be [DeviceFileEvents], got %v", parsedArgs.KQLTables)
	}
}
