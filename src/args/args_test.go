package args

import (
	"reflect"
	"testing"
)

func TestParseArgs_BasicUsage(t *testing.T) {
	tests := []struct {
		name     string
		rawArgs  []string
		expected Args
		hasError bool
	}{
		{
			name:    "default args",
			rawArgs: []string{},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "input directory only",
			rawArgs: []string{"-i", "/test/path"},
			expected: Args{
				StrInputDir:       "/test/path",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "output file specified",
			rawArgs: []string{"-o", "output.csv"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "output.csv",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  false,
				WriteToFile:       true,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "multiple algorithms",
			rawArgs: []string{"-a", "md5", "sha256", "sha512"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5", "sha256", "sha512"},
				HashAlgorithmId:   []int{0, 2, 3},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "terminal output flag",
			rawArgs: []string{"-t"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "format option",
			rawArgs: []string{"-f", "condensed"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "condensed",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "help flag",
			rawArgs: []string{"-h"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        false,
				YaraFile:          "",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              true,
			},
			hasError: false,
		},
		{
			name:    "yara output flag",
			rawArgs: []string{"-y", "rules.yar"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        true,
				YaraFile:          "rules.yar",
				YaraRuleName:      "",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "yara with rule name",
			rawArgs: []string{"-y", "rules.yar", "--yara-rule-name", "malware_detection"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        true,
				YaraFile:          "rules.yar",
				YaraRuleName:      "malware_detection",
				YaraHashOnly:      false,
				Help:              false,
			},
			hasError: false,
		},
		{
			name:    "yara hash only mode",
			rawArgs: []string{"-y", "rules.yar", "--yara-hash-only"},
			expected: Args{
				StrInputDir:       ".",
				StrOutputFile:     "",
				StrHashAlgorithms: []string{"md5"},
				HashAlgorithmId:   []int{0},
				OutputToTerminal:  true,
				WriteToFile:       false,
				OutputFormat:      "standard",
				YaraOutput:        true,
				YaraFile:          "rules.yar",
				YaraRuleName:      "",
				YaraHashOnly:      true,
				Help:              false,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseArgs(tt.rawArgs)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.hasError && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseArgs() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

func TestParseArgs_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		rawArgs []string
	}{
		{
			name:    "missing input directory value",
			rawArgs: []string{"-i"},
		},
		{
			name:    "missing output file value",
			rawArgs: []string{"-o"},
		},
		// Note: -a without values is handled by defaulting to md5, not an error
		{
			name:    "missing format value",
			rawArgs: []string{"-f"},
		},
		{
			name:    "invalid algorithm",
			rawArgs: []string{"-a", "invalid"},
		},
		{
			name:    "invalid format",
			rawArgs: []string{"-f", "invalid"},
		},
		{
			name:    "unknown flag",
			rawArgs: []string{"-x"},
		},
		{
			name:    "unexpected argument",
			rawArgs: []string{"unexpected"},
		},
		{
			name:    "missing yara file value",
			rawArgs: []string{"-y"},
		},
		{
			name:    "missing yara rule name value",
			rawArgs: []string{"--yara-rule-name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseArgs(tt.rawArgs)
			if err == nil {
				t.Errorf("Expected error but got none for args: %v", tt.rawArgs)
			}
		})
	}
}

func TestParseArgs_LongFlags(t *testing.T) {
	rawArgs := []string{
		"--input-dir", "/test/dir",
		"--output", "test.csv",
		"--algorithm", "sha256", "md5",
		"--format", "ioc",
		"--terminal",
		"--yara", "rules.yar",
		"--yara-rule-name", "test_rule",
		"--yara-hash-only",
		"--help",
	}

	expected := Args{
		StrInputDir:       "/test/dir",
		StrOutputFile:     "test.csv",
		StrHashAlgorithms: []string{"sha256", "md5"},
		HashAlgorithmId:   []int{2, 0},
		OutputToTerminal:  true,
		WriteToFile:       true,
		OutputFormat:      "ioc",
		YaraOutput:        true,
		YaraFile:          "rules.yar",
		YaraRuleName:      "test_rule",
		YaraHashOnly:      true,
		Help:              true,
	}

	result, err := ParseArgs(rawArgs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseArgs() = %+v, expected %+v", result, expected)
	}
}

func TestParseArgs_ComplexCase(t *testing.T) {
	rawArgs := []string{
		"-i", "/suspicious/files",
		"-o", "iocs.csv",
		"-a", "md5", "sha1", "sha256", "sha512",
		"-f", "ioc",
		"-y", "malware.yar",
		"--yara-rule-name", "suspicious_files",
	}

	result, err := ParseArgs(rawArgs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.StrInputDir != "/suspicious/files" {
		t.Errorf("Expected input dir '/suspicious/files', got '%s'", result.StrInputDir)
	}
	if result.StrOutputFile != "iocs.csv" {
		t.Errorf("Expected output file 'iocs.csv', got '%s'", result.StrOutputFile)
	}
	if result.OutputFormat != "ioc" {
		t.Errorf("Expected format 'ioc', got '%s'", result.OutputFormat)
	}
	if !result.YaraOutput {
		t.Error("Expected YARA output to be enabled")
	}
	if result.YaraFile != "malware.yar" {
		t.Errorf("Expected YARA file 'malware.yar', got '%s'", result.YaraFile)
	}
	if result.YaraRuleName != "suspicious_files" {
		t.Errorf("Expected YARA rule name 'suspicious_files', got '%s'", result.YaraRuleName)
	}
	if len(result.HashAlgorithmId) != 4 {
		t.Errorf("Expected 4 algorithms, got %d", len(result.HashAlgorithmId))
	}
	expectedIds := []int{0, 1, 2, 3} // md5, sha1, sha256, sha512
	if !reflect.DeepEqual(result.HashAlgorithmId, expectedIds) {
		t.Errorf("Expected algorithm IDs %v, got %v", expectedIds, result.HashAlgorithmId)
	}
}