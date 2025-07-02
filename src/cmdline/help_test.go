package cmdline

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call PrintHelp
	PrintHelp()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Test that help contains expected sections
	expectedSections := []string{
		"Usage: dirhash [options]",
		"Options:",
		"-i, --input-dir",
		"-o, --output",
		"-a, --algorithm",
		"-f, --format",
		"-t, --terminal",
		"-h, --help",
		"Supported algorithms:",
		"md5, sha1, sha256, sha512",
		"Supported output formats:",
		"standard",
		"condensed",
		"ioc",
		"Examples:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Help output should contain '%s'", section)
		}
	}
}

func TestPrintHelp_AlgorithmDescriptions(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that all supported algorithms are mentioned
	algorithms := []string{"md5", "sha1", "sha256", "sha512"}
	for _, algo := range algorithms {
		if !strings.Contains(output, algo) {
			t.Errorf("Help should mention algorithm '%s'", algo)
		}
	}
}

func TestPrintHelp_FormatDescriptions(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check format descriptions
	formatDescriptions := []string{
		"Traditional format with separate rows per hash type",
		"All hashes on single row per file",
		"IOC-friendly format for security tools",
		"YARA, KQL, Sentinel",
	}

	for _, desc := range formatDescriptions {
		if !strings.Contains(output, desc) {
			t.Errorf("Help should contain format description '%s'", desc)
		}
	}
}

func TestPrintHelp_Examples(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that examples contain different format options
	examples := []string{
		"dirhash -i",
		"--format condensed",
		"-f ioc",
		"dirhash -t",
		"dirhash --help",
	}

	for _, example := range examples {
		if !strings.Contains(output, example) {
			t.Errorf("Help should contain example '%s'", example)
		}
	}
}

func TestPrintHelp_OptionDescriptions(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check detailed option descriptions
	optionDescriptions := []string{
		"Specify the input directory",
		"Specify the output file",
		"Specify the hash algorithms",
		"Specify the output format for both terminal and file output",
		"Output to terminal",
		"Show this help message and exit",
	}

	for _, desc := range optionDescriptions {
		if !strings.Contains(output, desc) {
			t.Errorf("Help should contain option description '%s'", desc)
		}
	}
}

func TestPrintHelp_DefaultValues(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that default values are mentioned
	defaults := []string{
		"default: current directory",
		"default: no output file",
		"default: md5",
		"default: standard",
		"default: false",
	}

	for _, defaultVal := range defaults {
		if !strings.Contains(output, defaultVal) {
			t.Errorf("Help should mention default value '%s'", defaultVal)
		}
	}
}

func TestPrintHelp_SecurityUseCase(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that security/IOC use case is mentioned
	securityTerms := []string{
		"security tools",
		"YARA",
		"KQL",
		"Sentinel",
	}

	securityMentioned := false
	for _, term := range securityTerms {
		if strings.Contains(output, term) {
			securityMentioned = true
			break
		}
	}

	if !securityMentioned {
		t.Error("Help should mention security tools or IOC use cases")
	}
}

func TestPrintHelp_LongFormOptions(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that both short and long form options are shown
	optionPairs := [][]string{
		{"-i", "--input-dir"},
		{"-o", "--output"},
		{"-a", "--algorithm"},
		{"-f", "--format"},
		{"-t", "--terminal"},
		{"-h", "--help"},
	}

	for _, pair := range optionPairs {
		shortForm, longForm := pair[0], pair[1]
		if !strings.Contains(output, shortForm) {
			t.Errorf("Help should contain short form option '%s'", shortForm)
		}
		if !strings.Contains(output, longForm) {
			t.Errorf("Help should contain long form option '%s'", longForm)
		}
	}
}

func TestPrintHelp_MultipleAlgorithmExample(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that examples show multiple algorithms can be specified
	if !strings.Contains(output, "sha512 sha1") || !strings.Contains(output, "md5 sha1 sha256 sha512") {
		t.Error("Help should show examples with multiple algorithms")
	}
}

func TestPrintHelp_OutputNotEmpty(t *testing.T) {
	// Capture stdout
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintHelp()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Basic check that help produces output
	if len(output) == 0 {
		t.Error("PrintHelp should produce non-empty output")
	}

	// Check that output has reasonable length (should be substantial help text)
	if len(output) < 500 {
		t.Errorf("Help output seems too short (%d characters), expected substantial help text", len(output))
	}
}