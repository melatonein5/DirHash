// Package main implements DirHash, a command-line tool for generating file hashes and YARA rules.
//
// DirHash is a comprehensive file hashing utility designed for security analysis,
// forensics, and malware research. It can process directories of files to generate
// cryptographic hashes and export detection rules in YARA format.
//
// # Features
//
//   - Multiple hash algorithms: MD5, SHA1, SHA256, SHA512
//   - Multiple output formats: standard, condensed, IOC-friendly
//   - YARA rule generation for malware detection
//   - Concurrent file processing for performance
//   - Terminal and file output support
//
// # Usage
//
// Basic file hashing:
//
//	dirhash -i /path/to/files -a sha256 -o hashes.csv
//
// Generate YARA rules:
//
//	dirhash -i /suspicious/files -y malware.yar --yara-rule-name threat_detection
//
// Multiple algorithms with IOC format:
//
//	dirhash -i /files -a md5 sha1 sha256 -f ioc -o iocs.csv
//
// # Command Line Options
//
//   - -i, --input-dir: Input directory to process (default: current directory)
//   - -o, --output: Output file for hash results
//   - -a, --algorithm: Hash algorithms to use (md5, sha1, sha256, sha512)
//   - -f, --format: Output format (standard, condensed, ioc)
//   - -t, --terminal: Output to terminal
//   - -y, --yara: Generate YARA rule file
//   - --yara-rule-name: Custom name for YARA rule
//   - --yara-hash-only: Generate hash-only YARA rules
//   - -h, --help: Show help message
//
// # Output Formats
//
//   - standard: Traditional format with separate rows per hash type
//   - condensed: All hashes on single row per file
//   - ioc: IOC-friendly format for security tools (YARA, KQL, Sentinel)
//
// # YARA Rule Generation
//
// DirHash can generate YARA rules for malware detection:
//
//   - Standard rules: Include both file hashes and filenames
//   - Hash-only rules: Include only cryptographic hashes
//   - Support for all hash algorithms
//   - Automatic rule name sanitization
//   - Metadata includes author, date, and tags
//
// # Performance
//
// DirHash uses concurrent processing to maximize performance when hashing
// large numbers of files. The tool is optimized for security analysis
// workflows where processing speed is critical.
package main

import (
	"log"
	"os"

	"github.com/melatonein5/DirHash/src/args"
	"github.com/melatonein5/DirHash/src/cmdline"
	"github.com/melatonein5/DirHash/src/files"
	"github.com/melatonein5/DirHash/src/kql"
	"github.com/melatonein5/DirHash/src/yara"
)

// arguments holds the parsed command-line arguments for the application.
// This global variable is populated during initialization and used throughout
// the main execution flow to control program behavior.
var arguments args.Args

// init parses command-line arguments and handles early program setup.
//
// This function runs automatically before main() and performs the following:
//   - Parses command-line arguments using the args package
//   - Handles help requests by printing usage and exiting
//   - Validates argument consistency
//   - Terminates the program on argument parsing errors
//
// The function will call os.Exit(0) if help is requested or log.Fatal()
// if argument parsing fails, preventing main() from executing.
func init() {
	// First, grab the args from the command line, ignoring the first argument which is the program name, then parse them
	rawArgs := os.Args[1:]
	var err error
	arguments, err = args.ParseArgs(rawArgs)
	if err != nil {
		log.Fatal(err)
	}

	if arguments.Help {
		cmdline.PrintHelp()
		// Exit the program after printing help
		os.Exit(0)
	}
}

// main executes the core DirHash workflow: file enumeration, hashing, and output generation.
//
// The main function orchestrates the complete file processing pipeline:
//
//  1. File Discovery: Enumerates all files in the specified input directory
//  2. Hash Generation: Computes cryptographic hashes using specified algorithms
//  3. Output Generation: Writes results to terminal and/or files in requested format
//  4. YARA Export: Optionally generates YARA rules for malware detection
//
// The function handles errors at each stage and will terminate the program
// with log.Fatalf() if any critical step fails. Progress information is
// logged to help users track processing status.
//
// Output behavior is controlled by command-line arguments:
//   - Terminal output is displayed if OutputToTerminal is true
//   - File output is written if WriteToFile is true
//   - YARA rules are generated if YaraOutput is true
//
// The function supports multiple output formats (standard, condensed, IOC)
// and automatically selects the appropriate formatting function based on
// the user's choice.
func main() {
	// Enumerate the files in the input directory
	fs, err := files.EnumerateFiles(arguments.StrInputDir)
	if err != nil {
		log.Fatalf("Error enumerating files: %v", err)
	}

	log.Printf("Found %d files to process", len(fs))

	// Hash the files using the specified algorithms concurrently
	hashedFiles, err := files.HashFiles(fs, arguments.HashAlgorithmId)
	if err != nil {
		log.Fatalf("Error hashing files: %v", err)
	}

	log.Printf("Successfully processed %d files", len(hashedFiles))

	// Check if the output should be written to a file or printed to the terminal
	if arguments.OutputToTerminal {
		switch arguments.OutputFormat {
		case "condensed":
			cmdline.OutputFilesCondensed(hashedFiles)
		case "ioc":
			cmdline.OutputFilesIOC(hashedFiles)
		default: // "standard"
			cmdline.OutputFiles(hashedFiles)
		}
	}

	if arguments.WriteToFile {
		// Write the files to the output file using the specified format
		var err error
		switch arguments.OutputFormat {
		case "condensed":
			err = files.WriteOutputCondensed(hashedFiles, arguments.StrOutputFile)
		case "ioc":
			err = files.WriteOutputForIOC(hashedFiles, arguments.StrOutputFile)
		default: // "standard"
			err = files.WriteOutput(hashedFiles, arguments.StrOutputFile)
		}

		if err != nil {
			log.Fatalf("Error writing files to output file: %v", err)
		}
		log.Printf("Output written to: %s (format: %s)", arguments.StrOutputFile, arguments.OutputFormat)
	}

	// Generate YARA rule if requested
	if arguments.YaraOutput {
		err := generateYaraRule(hashedFiles)
		if err != nil {
			log.Fatalf("Error generating YARA rule: %v", err)
		}
	}

	// Generate KQL query if requested
	if arguments.KQLOutput {
		err := generateKQLQuery(hashedFiles)
		if err != nil {
			log.Fatalf("Error generating KQL query: %v", err)
		}
	}
}

// generateYaraRule creates and writes a YARA rule based on the processed files.
//
// This function generates YARA rules for malware detection and file identification
// based on the cryptographic hashes and metadata of the processed files.
//
// Parameters:
//   - hashedFiles: Slice of File structs containing hash data and metadata
//
// Returns:
//   - error: Any error that occurred during rule generation or file writing
//
// The function supports two YARA rule generation modes:
//
//  1. Standard Mode (default): Generates rules with both hash strings and filename strings,
//     providing multiple detection vectors for comprehensive coverage.
//
//  2. Hash-Only Mode: Generates rules containing only cryptographic hash patterns,
//     useful for situations where filename-based detection might produce false positives.
//
// The generated YARA rule includes:
//   - Metadata section with author, description, date, and categorization tags
//   - String definitions for file hashes and/or filenames
//   - Logical conditions using AND/OR operators for flexible matching
//   - Proper YARA syntax compliance and identifier sanitization
//
// Rule names are automatically sanitized to ensure YARA compliance by replacing
// invalid characters and ensuring proper identifier structure. If no rule name
// is specified, a default name "dirhash_generated_rule" is used.
//
// The function writes the generated rule to the file path specified in the
// global arguments.YaraFile and logs the operation result.
func generateYaraRule(hashedFiles []*files.File) error {
	var rule *yara.YaraRule
	var err error

	// Determine rule name
	ruleName := arguments.YaraRuleName
	if ruleName == "" {
		ruleName = "dirhash_generated_rule"
	}

	// Generate rule based on mode
	if arguments.YaraHashOnly {
		// Hash-only mode: only include hash-based conditions
		hashTypes := append([]string{}, arguments.StrHashAlgorithms...)
		rule, err = yara.GenerateYaraRuleFromHashes(hashedFiles, ruleName, hashTypes)
	} else {
		// Standard mode: include both hashes and filenames
		rule, err = yara.GenerateYaraRule(hashedFiles, ruleName)
	}

	if err != nil {
		return err
	}

	// Write YARA rule to file
	yaraContent := rule.ToYaraFormat()
	err = os.WriteFile(arguments.YaraFile, []byte(yaraContent), 0644)
	if err != nil {
		return err
	}

	log.Printf("YARA rule written to: %s (rule name: %s)", arguments.YaraFile, rule.Name)
	return nil
}

// generateKQLQuery creates and writes a KQL query based on the processed files.
//
// This function generates KQL (Kusto Query Language) queries for threat hunting
// and security analysis in Microsoft Sentinel, Azure Log Analytics, and other
// KQL-enabled security platforms.
//
// Parameters:
//   - hashedFiles: Slice of File structs containing hash data and metadata
//
// Returns:
//   - error: Any error that occurred during query generation or file writing
//
// The function supports two KQL query generation modes:
//
//  1. Standard Mode (default): Generates queries with both hash-based and filename-based
//     search conditions, providing comprehensive detection coverage across multiple
//     security log sources.
//
//  2. Hash-Only Mode: Generates queries containing only cryptographic hash searches,
//     useful for scenarios where filename-based detection might produce false positives
//     or when analyzing files that frequently change names.
//
// The generated KQL query includes:
//   - Metadata comments with author, description, generation date, and tags
//   - Multi-table search capabilities (DeviceFileEvents, SecurityEvents, etc.)
//   - Proper KQL syntax with efficient operators (in, contains, has)
//   - Time range filtering and result limiting for performance optimization
//   - Field selection optimized for security analysis workflows
//
// Query names are automatically sanitized to ensure KQL compliance by replacing
// invalid characters with underscores and ensuring proper identifier structure.
// If no query name is specified, a default name "dirhash_generated_query" is used.
//
// The function supports customizable target tables through the KQLTables argument,
// allowing users to specify which log sources to search (e.g., DeviceFileEvents,
// SecurityEvents, CommonSecurityLog).
//
// The function writes the generated query to the file path specified in the
// global arguments.KQLFile and logs the operation result.
func generateKQLQuery(hashedFiles []*files.File) error {
	var query *kql.KQLQuery
	var err error

	// Determine query name
	queryName := arguments.KQLName
	if queryName == "" {
		queryName = "dirhash_generated_query"
	}

	// Prepare KQL options
	options := kql.DefaultKQLQueryOptions()
	options.Tables = arguments.KQLTables
	options.IncludeHashes = true
	options.IncludeFilenames = !arguments.KQLHashOnly

	// Generate query based on mode
	if arguments.KQLHashOnly {
		// Hash-only mode: only include hash-based conditions
		query, err = kql.GenerateKQLQueryHashOnly(hashedFiles, queryName, arguments.StrHashAlgorithms)
	} else {
		// Standard mode: include both hashes and filenames
		query, err = kql.GenerateKQLQueryWithOptions(hashedFiles, queryName, arguments.StrHashAlgorithms, options)
	}

	if err != nil {
		return err
	}

	// Write KQL query to file
	kqlContent := query.ToKQLFormat()
	err = os.WriteFile(arguments.KQLFile, []byte(kqlContent), 0644)
	if err != nil {
		return err
	}

	log.Printf("KQL query written to: %s (query name: %s)", arguments.KQLFile, query.Name)
	return nil
}
