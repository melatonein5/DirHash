// Package cmdline provides command-line interface functionality for DirHash.
//
// This package handles terminal output formatting, help text generation, and
// user interface elements for the DirHash command-line application. It provides
// functions for displaying file processing results and usage information.
//
// # Terminal Output Functions
//
//   - Standard output formatting for hash results
//   - Condensed format for compact display
//   - IOC format for security analysis workflows
//   - Help text and usage information
//
// # Output Formatting
//
// The package supports multiple output formats optimized for different use cases:
//
//   - Standard: Traditional tabular format with clear column separation
//   - Condensed: Space-efficient format with all hashes on one row
//   - IOC: Security tool-friendly format for integration with analysis platforms
//
// # Usage
//
//	// Display results in standard format
//	cmdline.OutputFiles(hashedFiles)
//
//	// Display results in condensed format
//	cmdline.OutputFilesCondensed(hashedFiles)
//
//	// Display results in IOC format
//	cmdline.OutputFilesIOC(hashedFiles)
//
//	// Show help information
//	cmdline.PrintHelp()
package cmdline

import "fmt"

// PrintHelp displays comprehensive usage information for the DirHash application.
//
// This function outputs detailed help text including all supported command-line
// options, arguments, examples, and format descriptions. The help text is
// designed to provide users with complete information needed to effectively
// use DirHash for their file hashing and analysis needs.
//
// The help output includes:
//   - Command-line syntax and option descriptions
//   - Supported hash algorithms and their characteristics
//   - Available output formats and their use cases
//   - Practical usage examples for common scenarios
//   - YARA rule generation options and examples
func PrintHelp() {
	helpText := `
Usage: dirhash [options]

File Processing Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format for both terminal and file output (default: standard)
  -t, --terminal           Output to terminal (default: false)

YARA Rule Generation Options:
  -y, --yara <file>        Generate YARA rule and save to specified file
  --yara-rule-name <name>  Specify custom name for generated YARA rule
  --yara-hash-only         Generate hash-only rules without filenames

KQL Query Generation Options:
  -q, --kql <file>         Generate KQL query and save to specified file
  --kql-name <name>        Specify custom name for generated KQL query
  --kql-hash-only          Generate hash-only queries without filenames
  --kql-tables <tables>    Specify target tables (default: DeviceFileEvents), can take more than 1 argument

General Options:
  -h, --help               Show this help message and exit

Supported algorithms:
  md5, sha1, sha256, sha512

Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)

Supported KQL tables:
  DeviceFileEvents    - Microsoft 365 Defender file events (default)
  SecurityEvents      - Windows security events
  CommonSecurityLog   - Common security log format

Examples:
  Basic file hashing:
    dirhash -i /path/to/dir -o output.csv -a sha256
    dirhash --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
    dirhash -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc

  YARA rule generation:
    dirhash -i /malware/samples -y detection.yar --yara-rule-name malware_detection
    dirhash -i /files -a sha256 sha512 -y hashes.yar --yara-hash-only
    dirhash -i /suspicious -o results.csv -y rules.yar --yara-rule-name threat_hunt

  KQL query generation:
    dirhash -i /malware/samples -q detection.kql --kql-name malware_hunt
    dirhash -i /files -a sha256 sha512 -q hashes.kql --kql-hash-only
    dirhash -i /suspicious -q security.kql --kql-tables DeviceFileEvents SecurityEvents
    dirhash -i /threats -o iocs.csv -q hunt.kql --kql-name threat_detection

  Combined YARA and KQL generation:
    dirhash -i /malware -o results.csv -y rules.yar -q queries.kql -a sha256 sha512
    dirhash -i /samples -y detection.yar -q hunting.kql --yara-rule-name malware --kql-name threats

  Terminal output:
    dirhash -t
    dirhash -i /files -t -a sha256

  Help:
    dirhash --help
`
	fmt.Print(helpText)
}
