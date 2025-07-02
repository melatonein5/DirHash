package args

import (
	"testing"
)

func TestStrHashAlgorithmToId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "md5 algorithm",
			input:    "md5",
			expected: 0,
		},
		{
			name:     "sha1 algorithm",
			input:    "sha1",
			expected: 1,
		},
		{
			name:     "sha256 algorithm",
			input:    "sha256",
			expected: 2,
		},
		{
			name:     "sha512 algorithm",
			input:    "sha512",
			expected: 3,
		},
		{
			name:     "uppercase MD5",
			input:    "MD5",
			expected: -1, // Should be case sensitive
		},
		{
			name:     "invalid algorithm",
			input:    "invalid",
			expected: -1,
		},
		{
			name:     "empty string",
			input:    "",
			expected: -1,
		},
		{
			name:     "sha3 (unsupported)",
			input:    "sha3",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrHashAlgorithmToId(tt.input)
			if result != tt.expected {
				t.Errorf("StrHashAlgorithmToId(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStrHashAlgorithmToId_AllValidAlgorithms(t *testing.T) {
	validAlgorithms := map[string]int{
		"md5":    0,
		"sha1":   1,
		"sha256": 2,
		"sha512": 3,
	}

	for algo, expectedId := range validAlgorithms {
		t.Run("valid_"+algo, func(t *testing.T) {
			result := StrHashAlgorithmToId(algo)
			if result != expectedId {
				t.Errorf("StrHashAlgorithmToId(%s) = %d, expected %d", algo, result, expectedId)
			}
		})
	}
}

func TestStrHashAlgorithmToId_ConsistencyWithParseArgs(t *testing.T) {
	// Test that the algorithm mapping is consistent with ParseArgs
	algorithms := []string{"md5", "sha1", "sha256", "sha512"}
	
	for i, algo := range algorithms {
		expectedId := i
		actualId := StrHashAlgorithmToId(algo)
		
		if actualId != expectedId {
			t.Errorf("Algorithm %s should map to ID %d, but got %d", algo, expectedId, actualId)
		}
	}
}

func TestStrHashAlgorithmToId_BoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"whitespace", " md5 "},
		{"tab", "\tsha1\t"},
		{"newline", "sha256\n"},
		{"mixed_case", "Sha512"},
		{"partial_match", "md"},
		{"extra_chars", "md5_extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrHashAlgorithmToId(tt.input)
			if result != -1 {
				t.Errorf("StrHashAlgorithmToId(%q) should return -1 for invalid input, got %d", tt.input, result)
			}
		})
	}
}