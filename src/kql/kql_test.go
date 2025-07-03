package kql

import (
	"strings"
	"testing"
	"time"

	"github.com/melatonein5/DirHash/src/files"
)

// TestGenerateKQLQuery tests basic KQL query generation functionality
func TestGenerateKQLQuery(t *testing.T) {
	// Create test files
	testFiles := []*files.File{
		{
			FileName: "malware.exe",
			Path:     "/tmp/malware.exe",
			Size:     1024,
			Hashes: map[string]string{
				"md5":    "abc123def456",
				"sha256": "def456abc123ghi789",
			},
		},
		{
			FileName: "trojan.dll",
			Path:     "/tmp/trojan.dll",
			Size:     2048,
			Hashes: map[string]string{
				"md5":    "ghi789jkl012",
				"sha256": "jkl012mno345pqr678",
			},
		},
	}

	// Test basic query generation
	query, err := GenerateKQLQuery(testFiles, "test_query", []string{"md5", "sha256"})
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Verify query structure
	if query.Name != "test_query" {
		t.Errorf("Expected query name 'test_query', got '%s'", query.Name)
	}

	if query.Author != "DirHash" {
		t.Errorf("Expected author 'DirHash', got '%s'", query.Author)
	}

	if len(query.HashTypes) != 2 {
		t.Errorf("Expected 2 hash types, got %d", len(query.HashTypes))
	}

	if len(query.FilenameList) != 2 {
		t.Errorf("Expected 2 filenames, got %d", len(query.FilenameList))
	}

	// Verify query content
	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "DeviceFileEvents") {
		t.Error("Query should contain DeviceFileEvents table")
	}

	if !strings.Contains(kqlContent, "abc123def456") {
		t.Error("Query should contain MD5 hash")
	}

	if !strings.Contains(kqlContent, "def456abc123ghi789") {
		t.Error("Query should contain SHA256 hash")
	}

	if !strings.Contains(kqlContent, "malware.exe") {
		t.Error("Query should contain filename")
	}

	if !strings.Contains(kqlContent, "trojan.dll") {
		t.Error("Query should contain filename")
	}
}

// TestGenerateKQLQueryHashOnly tests hash-only query generation
func TestGenerateKQLQueryHashOnly(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "suspicious.exe",
			Path:     "/tmp/suspicious.exe",
			Size:     1024,
			Hashes: map[string]string{
				"md5":    "hash123",
				"sha256": "hash456",
			},
		},
	}

	query, err := GenerateKQLQueryHashOnly(testFiles, "hash_only_query", []string{"md5", "sha256"})
	if err != nil {
		t.Fatalf("Failed to generate hash-only KQL query: %v", err)
	}

	// Verify query structure
	if query.Name != "hash_only_query" {
		t.Errorf("Expected query name 'hash_only_query', got '%s'", query.Name)
	}

	if len(query.FilenameList) != 0 {
		t.Errorf("Hash-only query should not contain filenames, got %d", len(query.FilenameList))
	}

	// Verify query content
	kqlContent := query.ToKQLFormat()
	if strings.Contains(kqlContent, "suspicious.exe") {
		t.Error("Hash-only query should not contain filenames")
	}

	if !strings.Contains(kqlContent, "hash123") {
		t.Error("Hash-only query should contain MD5 hash")
	}

	if !strings.Contains(kqlContent, "hash456") {
		t.Error("Hash-only query should contain SHA256 hash")
	}
}

// TestGenerateKQLQueryWithOptions tests custom options
func TestGenerateKQLQueryWithOptions(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Path:     "/tmp/test.exe",
			Size:     1024,
			Hashes: map[string]string{
				"md5":    "test123",
				"sha256": "test456",
			},
		},
	}

	// Test custom options
	options := KQLQueryOptions{
		Tables:           []string{"DeviceFileEvents", "SecurityEvents"},
		TimeRange:        "30d",
		MaxResults:       5000,
		IncludeHashes:    true,
		IncludeFilenames: true,
		CaseSensitive:    true,
		IncludeMetadata:  true,
		IncludeComments:  true,
		CompactFormat:    false,
	}

	query, err := GenerateKQLQueryWithOptions(testFiles, "custom_query", []string{"md5", "sha256"}, options)
	if err != nil {
		t.Fatalf("Failed to generate KQL query with options: %v", err)
	}

	// Verify options are applied
	if len(query.Tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(query.Tables))
	}

	if query.TimeRange != "30d" {
		t.Errorf("Expected time range '30d', got '%s'", query.TimeRange)
	}

	if query.MaxResults != 5000 {
		t.Errorf("Expected max results 5000, got %d", query.MaxResults)
	}

	// Verify query content
	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "DeviceFileEvents") {
		t.Error("Query should contain DeviceFileEvents table")
	}

	if !strings.Contains(kqlContent, "SecurityEvents") {
		t.Error("Query should contain SecurityEvents table")
	}

	if !strings.Contains(kqlContent, "30d") {
		t.Error("Query should contain time range")
	}

	if !strings.Contains(kqlContent, "take 5000") {
		t.Error("Query should contain result limit")
	}

	if !strings.Contains(kqlContent, "union") {
		t.Error("Query should contain union for multiple tables")
	}
}

// TestKQLQuerySanitization tests query name sanitization
func TestKQLQuerySanitization(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Hashes: map[string]string{
				"md5": "test123",
			},
		},
	}

	// Test with invalid characters
	query, err := GenerateKQLQuery(testFiles, "test-query!@#$%^&*()", []string{"md5"})
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Verify sanitization
	if strings.Contains(query.Name, "-") || strings.Contains(query.Name, "!") {
		t.Errorf("Query name should be sanitized, got '%s'", query.Name)
	}

	if !strings.Contains(query.Name, "_") {
		t.Errorf("Sanitized query name should contain underscores, got '%s'", query.Name)
	}
}

// TestKQLQueryErrorHandling tests error conditions
func TestKQLQueryErrorHandling(t *testing.T) {
	// Test empty files list
	_, err := GenerateKQLQuery([]*files.File{}, "test", []string{"md5"})
	if err == nil {
		t.Error("Expected error for empty files list")
	}

	// Test files with no hashes
	filesWithoutHashes := []*files.File{
		{
			FileName: "test.exe",
			Hashes:   map[string]string{},
		},
	}

	query, err := GenerateKQLQuery(filesWithoutHashes, "test", []string{"md5"})
	if err != nil {
		t.Fatalf("Should not error for files without hashes: %v", err)
	}

	// Should still generate query with filename
	if len(query.FilenameList) == 0 {
		t.Error("Query should contain filename even without hashes")
	}
}

// TestKQLQueryMetadata tests metadata generation
func TestKQLQueryMetadata(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Hashes: map[string]string{
				"md5": "test123",
			},
		},
	}

	query, err := GenerateKQLQuery(testFiles, "metadata_test", []string{"md5"})
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Verify metadata
	if len(query.Tags) == 0 {
		t.Error("Query should have tags")
	}

	if query.Generated.IsZero() {
		t.Error("Query should have generation timestamp")
	}

	if time.Since(query.Generated) > time.Minute {
		t.Error("Query generation time should be recent")
	}

	// Verify metadata comments
	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "// KQL Query:") {
		t.Error("Query should contain metadata comments")
	}

	if !strings.Contains(kqlContent, "// Author:") {
		t.Error("Query should contain author metadata")
	}

	if !strings.Contains(kqlContent, "// Generated:") {
		t.Error("Query should contain generation timestamp")
	}

	if !strings.Contains(kqlContent, "// Tags:") {
		t.Error("Query should contain tags metadata")
	}
}

// TestKQLQueryDuplicateHandling tests duplicate hash and filename handling
func TestKQLQueryDuplicateHandling(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "duplicate.exe",
			Hashes: map[string]string{
				"md5": "duplicate123",
			},
		},
		{
			FileName: "duplicate.exe", // Same filename
			Hashes: map[string]string{
				"md5": "duplicate123", // Same hash
			},
		},
		{
			FileName: "different.exe",
			Hashes: map[string]string{
				"md5": "duplicate123", // Same hash, different filename
			},
		},
	}

	query, err := GenerateKQLQuery(testFiles, "duplicate_test", []string{"md5"})
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Verify deduplication
	if len(query.FilenameList) != 2 {
		t.Errorf("Expected 2 unique filenames, got %d", len(query.FilenameList))
	}

	if len(query.HashList) != 1 {
		t.Errorf("Expected 1 unique hash, got %d", len(query.HashList))
	}

	// Verify content
	kqlContent := query.ToKQLFormat()
	duplicateCount := strings.Count(kqlContent, "duplicate123")
	if duplicateCount != 1 {
		t.Errorf("Hash should appear only once in query, found %d times", duplicateCount)
	}
}

// TestKQLQueryMultipleHashTypes tests queries with multiple hash algorithms
func TestKQLQueryMultipleHashTypes(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "multi.exe",
			Hashes: map[string]string{
				"md5":    "md5hash123",
				"sha1":   "sha1hash456",
				"sha256": "sha256hash789",
				"sha512": "sha512hash012",
			},
		},
	}

	query, err := GenerateKQLQuery(testFiles, "multi_hash_test", []string{"md5", "sha1", "sha256", "sha512"})
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	// Verify all hash types are included
	if len(query.HashTypes) != 4 {
		t.Errorf("Expected 4 hash types, got %d", len(query.HashTypes))
	}

	// Verify query content
	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "md5hash123") {
		t.Error("Query should contain MD5 hash")
	}

	if !strings.Contains(kqlContent, "sha1hash456") {
		t.Error("Query should contain SHA1 hash")
	}

	if !strings.Contains(kqlContent, "sha256hash789") {
		t.Error("Query should contain SHA256 hash")
	}

	if !strings.Contains(kqlContent, "sha512hash012") {
		t.Error("Query should contain SHA512 hash")
	}
}

// TestKQLQueryTableSpecificFields tests table-specific field handling
func TestKQLQueryTableSpecificFields(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Hashes: map[string]string{
				"md5":    "test123",
				"sha256": "test456",
			},
		},
	}

	// Test DeviceFileEvents table
	options := DefaultKQLQueryOptions()
	options.Tables = []string{"DeviceFileEvents"}

	query, err := GenerateKQLQueryWithOptions(testFiles, "device_events_test", []string{"md5", "sha256"}, options)
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "MD5") {
		t.Error("DeviceFileEvents query should use MD5 field")
	}

	if !strings.Contains(kqlContent, "SHA256") {
		t.Error("DeviceFileEvents query should use SHA256 field")
	}

	if !strings.Contains(kqlContent, "FileName") {
		t.Error("DeviceFileEvents query should use FileName field")
	}

	// Test SecurityEvents table
	options.Tables = []string{"SecurityEvents"}
	query, err = GenerateKQLQueryWithOptions(testFiles, "security_events_test", []string{"md5", "sha256"}, options)
	if err != nil {
		t.Fatalf("Failed to generate KQL query: %v", err)
	}

	kqlContent = query.ToKQLFormat()
	if !strings.Contains(kqlContent, "FileHash") {
		t.Error("SecurityEvents query should use FileHash field")
	}
}

// TestKQLQueryCompactFormat tests compact format generation
func TestKQLQueryCompactFormat(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Hashes: map[string]string{
				"md5": "test123",
			},
		},
	}

	// Test compact format
	options := DefaultKQLQueryOptions()
	options.CompactFormat = true
	options.IncludeMetadata = false
	options.IncludeComments = false

	query, err := GenerateKQLQueryWithOptions(testFiles, "compact_test", []string{"md5"}, options)
	if err != nil {
		t.Fatalf("Failed to generate compact KQL query: %v", err)
	}

	kqlContent := query.ToKQLFormat()
	
	// Compact format should have minimal comments
	if len(query.Comments) != 0 {
		t.Errorf("Compact format should have no comments, got %d", len(query.Comments))
	}

	// Should not contain metadata comments
	if strings.Contains(kqlContent, "// KQL Query:") {
		t.Error("Compact format should not contain metadata comments")
	}
}

// TestKQLQueryDefaultOptions tests default options behavior
func TestKQLQueryDefaultOptions(t *testing.T) {
	options := DefaultKQLQueryOptions()

	// Verify default values
	if len(options.Tables) != 1 || options.Tables[0] != "DeviceFileEvents" {
		t.Errorf("Default tables should be [DeviceFileEvents], got %v", options.Tables)
	}

	if options.TimeRange != "7d" {
		t.Errorf("Default time range should be '7d', got '%s'", options.TimeRange)
	}

	if options.MaxResults != 1000 {
		t.Errorf("Default max results should be 1000, got %d", options.MaxResults)
	}

	if !options.IncludeHashes {
		t.Error("Default should include hashes")
	}

	if !options.IncludeFilenames {
		t.Error("Default should include filenames")
	}

	if options.CaseSensitive {
		t.Error("Default should not be case sensitive")
	}

	if !options.IncludeMetadata {
		t.Error("Default should include metadata")
	}

	if !options.IncludeComments {
		t.Error("Default should include comments")
	}

	if options.CompactFormat {
		t.Error("Default should not be compact format")
	}
}

// TestKQLQueryHelperFunctions tests internal helper functions
func TestKQLQueryHelperFunctions(t *testing.T) {
	// Test sanitizeKQLName
	testCases := []struct {
		input    string
		expected string
	}{
		{"valid_name", "valid_name"},
		{"invalid-name", "invalid_name"},
		{"name with spaces", "name_with_spaces"},
		{"name!@#$%^&*()", "name__________"},
		{"123name", "_123name"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := sanitizeKQLName(tc.input)
		if result != tc.expected {
			t.Errorf("sanitizeKQLName(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}

	// Test contains function
	slice := []string{"a", "b", "c"}
	if !contains(slice, "b") {
		t.Error("contains should return true for existing item")
	}

	if contains(slice, "d") {
		t.Error("contains should return false for non-existing item")
	}

	// Test removeDuplicatesAndSort
	input := []string{"c", "a", "b", "a", "c"}
	result := removeDuplicatesAndSort(input)
	expected := []string{"a", "b", "c"}

	if len(result) != len(expected) {
		t.Errorf("removeDuplicatesAndSort length mismatch: got %d, expected %d", len(result), len(expected))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("removeDuplicatesAndSort[%d] = %q, expected %q", i, result[i], v)
		}
	}

	// Test quoteStrings
	input2 := []string{"a", "b", "c"}
	result2 := quoteStrings(input2)
	expected2 := []string{`"a"`, `"b"`, `"c"`}

	if len(result2) != len(expected2) {
		t.Errorf("quoteStrings length mismatch: got %d, expected %d", len(result2), len(expected2))
	}

	for i, v := range expected2 {
		if result2[i] != v {
			t.Errorf("quoteStrings[%d] = %q, expected %q", i, result2[i], v)
		}
	}
}

// TestKQLQueryCaseSensitivity tests case sensitivity options
func TestKQLQueryCaseSensitivity(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "Test.EXE",
			Hashes: map[string]string{
				"md5": "test123",
			},
		},
	}

	// Test case-insensitive (default)
	options := DefaultKQLQueryOptions()
	options.CaseSensitive = false

	query, err := GenerateKQLQueryWithOptions(testFiles, "case_insensitive_test", []string{"md5"}, options)
	if err != nil {
		t.Fatalf("Failed to generate case-insensitive KQL query: %v", err)
	}

	kqlContent := query.ToKQLFormat()
	if !strings.Contains(kqlContent, "in~") {
		t.Error("Case-insensitive query should use 'in~' operator")
	}

	// Test case-sensitive
	options.CaseSensitive = true
	query, err = GenerateKQLQueryWithOptions(testFiles, "case_sensitive_test", []string{"md5"}, options)
	if err != nil {
		t.Fatalf("Failed to generate case-sensitive KQL query: %v", err)
	}

	kqlContent = query.ToKQLFormat()
	if strings.Contains(kqlContent, "in~") {
		t.Error("Case-sensitive query should not use 'in~' operator")
	}

	if !strings.Contains(kqlContent, " in (") {
		t.Error("Case-sensitive query should use 'in' operator")
	}
}