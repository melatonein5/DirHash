// Package args provides command-line argument parsing and validation for DirHash.
//
// This package handles the parsing, validation, and management of all command-line
// arguments and options for the DirHash application. It supports a comprehensive
// set of options for controlling file hashing, output formatting, and YARA rule generation.
//
// # Supported Arguments
//
// Input/Output Options:
//   - Input Directory: Specifies the directory to process for file hashing
//   - Output File: Defines the file path for saving hash results
//   - Terminal Output: Controls whether results are displayed on screen
//
// Hash Algorithm Options:
//   - Algorithm Selection: Choose from MD5, SHA1, SHA256, SHA512
//   - Multiple Algorithms: Support for computing multiple hash types simultaneously
//
// Output Formatting Options:
//   - Standard Format: Traditional row-per-hash output
//   - Condensed Format: All hashes on single row per file
//   - IOC Format: Security tool-friendly format for analysis platforms
//
// YARA Rule Generation Options:
//   - YARA Output: Enable generation of YARA rules
//   - Rule Naming: Custom names for generated rules
//   - Hash-Only Mode: Generate hash-based rules without filenames
//
// # Usage Example
//
//	args := []string{"-i", "/path/to/files", "-a", "sha256", "-o", "output.csv"}
//	parsedArgs, err := args.ParseArgs(args)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Validation
//
// The package performs comprehensive validation of arguments including:
//   - Directory existence checking
//   - Hash algorithm name validation
//   - Output format validation
//   - File path accessibility
//   - Argument combination consistency
//
// # Default Values
//
// When arguments are not specified, the package provides sensible defaults:
//   - Input Directory: Current working directory (".")
//   - Hash Algorithm: MD5
//   - Output Format: Standard
//   - Terminal Output: Enabled when no output file specified
package args

// Args represents the complete set of parsed command-line arguments.
//
// This structure contains all configuration options that control DirHash behavior,
// from input/output settings to hash algorithm selection and YARA rule generation.
// Fields are populated by ParseArgs() and used throughout the application.
type Args struct {
	// Input Directory Configuration
	StrInputDir string // Path to directory containing files to hash (default: ".")

	// Output File Configuration  
	StrOutputFile    string // Path to output file for hash results (empty = no file output)
	OutputToTerminal bool   // Whether to display results on terminal (default: true when no output file)
	WriteToFile      bool   // Whether to write results to file (set automatically when StrOutputFile provided)

	// Hash Algorithm Configuration
	StrHashAlgorithms []string // Human-readable algorithm names (e.g., ["md5", "sha256"])
	HashAlgorithmId   []int    // Internal algorithm IDs corresponding to StrHashAlgorithms

	// Output Format Configuration
	OutputFormat string // Output format: "standard", "condensed", or "ioc" (default: "standard")

	// YARA Rule Generation Configuration
	YaraOutput   bool   // Whether to generate YARA rules (default: false)
	YaraFile     string // Path to output YARA rule file (required when YaraOutput=true)
	YaraRuleName string // Custom name for generated YARA rule (default: auto-generated)
	YaraHashOnly bool   // Generate hash-only rules without filenames (default: false)

	// Application Control Flags
	Help bool // Whether help was requested (causes immediate exit after help display)
}
