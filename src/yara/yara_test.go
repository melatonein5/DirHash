package yara

import (
	"strings"
	"testing"
	"time"

	"github.com/melatonein5/DirHash/src/files"
)

func TestGenerateYaraRule(t *testing.T) {
	// Create test files
	testFiles := []*files.File{
		{
			FileName: "malware.exe",
			Path:     "/tmp/malware.exe",
			Size:     1024,
			ModTime:  time.Now(),
			Hashes: map[string]string{
				"md5":    "d41d8cd98f00b204e9800998ecf8427e",
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
		{
			FileName: "trojan.dll",
			Path:     "/tmp/trojan.dll",
			Size:     2048,
			ModTime:  time.Now(),
			Hashes: map[string]string{
				"md5":  "5d41402abc4b2a76b9719d911017c592",
				"sha1": "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			},
		},
	}

	rule, err := GenerateYaraRule(testFiles, "test_rule")
	if err != nil {
		t.Fatalf("GenerateYaraRule failed: %v", err)
	}

	// Verify rule properties
	if rule.Name != "test_rule" {
		t.Errorf("Expected rule name 'test_rule', got '%s'", rule.Name)
	}

	if rule.Author != "DirHash" {
		t.Errorf("Expected author 'DirHash', got '%s'", rule.Author)
	}

	if len(rule.Tags) == 0 {
		t.Error("Expected rule to have tags")
	}

	if len(rule.Strings) == 0 {
		t.Error("Expected rule to have strings")
	}

	if rule.Condition == "" {
		t.Error("Expected rule to have a condition")
	}

	// Verify YARA format output
	yaraOutput := rule.ToYaraFormat()
	if !strings.Contains(yaraOutput, "rule test_rule") {
		t.Error("YARA output should contain rule name")
	}

	if !strings.Contains(yaraOutput, "meta:") {
		t.Error("YARA output should contain meta section")
	}

	if !strings.Contains(yaraOutput, "strings:") {
		t.Error("YARA output should contain strings section")
	}

	if !strings.Contains(yaraOutput, "condition:") {
		t.Error("YARA output should contain condition section")
	}
}

func TestGenerateYaraRuleFromHashes(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "sample.exe",
			Path:     "/tmp/sample.exe",
			Size:     1024,
			ModTime:  time.Now(),
			Hashes: map[string]string{
				"md5":    "d41d8cd98f00b204e9800998ecf8427e",
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"sha512": "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
			},
		},
	}

	hashTypes := []string{"md5", "sha256"}
	rule, err := GenerateYaraRuleFromHashes(testFiles, "hash_rule", hashTypes)
	if err != nil {
		t.Fatalf("GenerateYaraRuleFromHashes failed: %v", err)
	}

	if rule.Name != "hash_rule" {
		t.Errorf("Expected rule name 'hash_rule', got '%s'", rule.Name)
	}

	// Verify that only hash strings are included
	foundMD5 := false
	foundSHA256 := false
	foundSHA512 := false
	foundFilename := false

	for _, str := range rule.Strings {
		if strings.Contains(str.Name, "md5") {
			foundMD5 = true
		}
		if strings.Contains(str.Name, "sha256") {
			foundSHA256 = true
		}
		if strings.Contains(str.Name, "sha512") {
			foundSHA512 = true
		}
		if strings.Contains(str.Name, "filename") {
			foundFilename = true
		}
	}

	if !foundMD5 {
		t.Error("Expected MD5 hash string")
	}
	if !foundSHA256 {
		t.Error("Expected SHA256 hash string")
	}
	if foundSHA512 {
		t.Error("Should not include SHA512 hash (not in hashTypes)")
	}
	if foundFilename {
		t.Error("Should not include filename strings in hash-only mode")
	}
}

func TestGenerateYaraRule_EmptyFiles(t *testing.T) {
	_, err := GenerateYaraRule([]*files.File{}, "test")
	if err == nil {
		t.Error("Expected error for empty files list")
	}
}

func TestGenerateYaraRuleFromHashes_EmptyFiles(t *testing.T) {
	_, err := GenerateYaraRuleFromHashes([]*files.File{}, "test", []string{"md5"})
	if err == nil {
		t.Error("Expected error for empty files list")
	}
}

func TestGenerateYaraRuleFromHashes_NoValidHashes(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "sample.exe",
			Hashes:   map[string]string{"md5": "somehash"},
		},
	}

	_, err := GenerateYaraRuleFromHashes(testFiles, "test", []string{"sha256"})
	if err == nil {
		t.Error("Expected error when no valid hashes found for specified types")
	}
}

func TestYaraRule_ToYaraFormat(t *testing.T) {
	rule := &YaraRule{
		Name:        "test_rule",
		Description: "Test rule description",
		Author:      "Test Author",
		Date:        "2023-01-01",
		Tags:        []string{"test", "generated"},
		Strings: []YaraString{
			{Name: "$md5_hash", Value: "D4 1D 8C D9 8F 00 B2 04", Type: "hex"},
			{Name: "$filename", Value: "malware.exe", Type: "text"},
		},
		Condition: "$md5_hash or $filename",
	}

	output := rule.ToYaraFormat()

	expectedParts := []string{
		"rule test_rule",
		"meta:",
		"description = \"Test rule description\"",
		"author = \"Test Author\"",
		"date = \"2023-01-01\"",
		"tags = \"test, generated\"",
		"strings:",
		"$md5_hash = { D4 1D 8C D9 8F 00 B2 04 }",
		"$filename = \"malware.exe\"",
		"condition:",
		"$md5_hash or $filename",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("YARA output should contain '%s'\nActual output:\n%s", part, output)
		}
	}
}

func TestSanitizeRuleName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"valid_rule", "valid_rule"},
		{"Valid-Rule.123", "Valid_Rule_123"},
		{"123invalid", "_123invalid"},
		{"", "generated_rule"},
		{"rule with spaces", "rule_with_spaces"},
		{"rule@#$%", "rule____"},
	}

	for _, test := range tests {
		result := sanitizeRuleName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeRuleName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestSanitizeStringName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"file.exe", "file"},
		{"malware.dll", "malware"},
		{"123file.txt", "_123file"},
		{"", "file"},
		{"file-name.bin", "file_name"},
		{"file@#$.exe", "file___"},
	}

	for _, test := range tests {
		result := sanitizeStringName(test.input)
		if result != test.expected {
			t.Errorf("sanitizeStringName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestFormatHashForYara(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"d41d8cd98f00b204", "D4 1D 8C D9 8F 00 B2 04"},
		{"abc123", "AB C1 23"},
		{"", ""},
		{"a", ""},
	}

	for _, test := range tests {
		result := formatHashForYara(test.input)
		if result != test.expected {
			t.Errorf("formatHashForYara(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestGenerateHashStrings(t *testing.T) {
	testFiles := []*files.File{
		{
			FileName: "test.exe",
			Hashes: map[string]string{
				"md5":    "d41d8cd98f00b204e9800998ecf8427e",
				"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
		},
	}

	strings := generateHashStrings(testFiles)
	if len(strings) != 2 {
		t.Errorf("Expected 2 hash strings, got %d", len(strings))
	}

	foundMD5 := false
	foundSHA256 := false
	for _, str := range strings {
		if str.Name == "$md5_test" && str.Type == "hex" {
			foundMD5 = true
		}
		if str.Name == "$sha256_test" && str.Type == "hex" {
			foundSHA256 = true
		}
	}

	if !foundMD5 {
		t.Error("Expected MD5 hash string")
	}
	if !foundSHA256 {
		t.Error("Expected SHA256 hash string")
	}
}

func TestGenerateFilenameStrings(t *testing.T) {
	testFiles := []*files.File{
		{FileName: "malware.exe"},
		{FileName: "trojan.dll"},
		{FileName: "malware.exe"}, // Duplicate should be ignored
	}

	strings := generateFilenameStrings(testFiles)
	if len(strings) != 2 {
		t.Errorf("Expected 2 filename strings (duplicates removed), got %d", len(strings))
	}

	foundMalware := false
	foundTrojan := false
	for _, str := range strings {
		if str.Value == "malware.exe" && str.Type == "text" {
			foundMalware = true
		}
		if str.Value == "trojan.dll" && str.Type == "text" {
			foundTrojan = true
		}
	}

	if !foundMalware {
		t.Error("Expected malware.exe filename string")
	}
	if !foundTrojan {
		t.Error("Expected trojan.dll filename string")
	}
}

func TestGenerateCondition(t *testing.T) {
	tests := []struct {
		name     string
		strings  []YaraString
		expected string
	}{
		{
			name:     "empty strings",
			strings:  []YaraString{},
			expected: "true",
		},
		{
			name: "single hash",
			strings: []YaraString{
				{Name: "$md5_hash", Type: "hex"},
			},
			expected: "$md5_hash",
		},
		{
			name: "multiple hashes",
			strings: []YaraString{
				{Name: "$md5_hash", Type: "hex"},
				{Name: "$sha256_hash", Type: "hex"},
			},
			expected: "any of ($md5_hash, $sha256_hash)",
		},
		{
			name: "hash and filename",
			strings: []YaraString{
				{Name: "$md5_hash", Type: "hex"},
				{Name: "$filename_malware", Type: "text"},
			},
			expected: "$md5_hash or $filename_malware",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generateCondition(test.strings)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestGenerateHashCondition(t *testing.T) {
	tests := []struct {
		name     string
		strings  []YaraString
		expected string
	}{
		{
			name:     "empty",
			strings:  []YaraString{},
			expected: "true",
		},
		{
			name: "single hash",
			strings: []YaraString{
				{Name: "$md5_hash"},
			},
			expected: "$md5_hash",
		},
		{
			name: "multiple hashes",
			strings: []YaraString{
				{Name: "$md5_hash"},
				{Name: "$sha256_hash"},
			},
			expected: "any of ($md5_hash, $sha256_hash)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generateHashCondition(test.strings)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestGetSupportedHashTypes(t *testing.T) {
	types := GetSupportedHashTypes()
	expectedTypes := []string{"md5", "sha1", "sha256", "sha512"}

	if len(types) != len(expectedTypes) {
		t.Errorf("Expected %d hash types, got %d", len(expectedTypes), len(types))
	}

	for _, expected := range expectedTypes {
		found := false
		for _, actual := range types {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected hash type '%s' not found", expected)
		}
	}
}
