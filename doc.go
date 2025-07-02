// Package dirhash provides a comprehensive file hashing and YARA rule generation utility.
//
// DirHash is designed for security analysts, malware researchers, and forensic investigators
// who need to generate cryptographic hashes for files and create YARA rules for detection.
//
// # Overview
//
// DirHash processes directories of files to:
//   - Generate cryptographic hashes (MD5, SHA1, SHA256, SHA512)
//   - Export results in multiple formats (standard, condensed, IOC)
//   - Create YARA rules for malware detection and file identification
//   - Support high-performance concurrent processing
//
// # Architecture
//
// The application is structured into several specialized packages:
//
//   - main: Application entry point and workflow orchestration
//   - args: Command-line argument parsing and validation
//   - files: Core file processing, hashing, and output generation
//   - cmdline: Terminal output formatting and help text
//   - yara: YARA rule generation and formatting
//
// # Quick Start
//
// Install and build DirHash:
//
//	git clone https://github.com/melatonein5/DirHash
//	cd DirHash
//	go build .
//
// Basic usage examples:
//
//	# Hash files with SHA256
//	./dirhash -i /path/to/files -a sha256 -o hashes.csv
//
//	# Generate YARA rules
//	./dirhash -i /suspicious/files -y malware.yar --yara-rule-name threat_detection
//
//	# Multiple algorithms with IOC format
//	./dirhash -i /files -a md5 sha1 sha256 -f ioc -o iocs.csv
//
// # Use Cases
//
// Security Analysis:
//   - Incident response file identification
//   - Malware sample cataloging
//   - IOC generation for threat hunting
//   - Forensic evidence processing
//
// Development and Testing:
//   - Build artifact verification
//   - Software integrity checking
//   - File change detection
//   - Release validation
//
// Research and Investigation:
//   - Malware family clustering
//   - Sample similarity analysis
//   - Detection rule creation
//   - Threat intelligence gathering
//
// # Performance
//
// DirHash is optimized for processing large file sets through:
//   - Concurrent hash computation across multiple files
//   - Streaming I/O for memory efficiency
//   - Minimal memory allocation during processing
//   - Support for files of any size
//
// # Output Formats
//
// Standard Format: Traditional CSV with separate rows per hash type
//
//	File Name,Path,Size,Hash,Hash Type
//	malware.exe,/tmp/malware.exe,1024,d41d8cd98f00b204e9800998ecf8427e,md5
//	malware.exe,/tmp/malware.exe,1024,e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855,sha256
//
// Condensed Format: All hashes on single row per file
//
//	File Name,Path,Size,MD5,SHA256
//	malware.exe,/tmp/malware.exe,1024,d41d8cd98f00b204e9800998ecf8427e,e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
//
// IOC Format: Security tool-friendly format with underscored headers
//
//	file_name,file_path,file_size,md5,sha256
//	malware.exe,/tmp/malware.exe,1024,d41d8cd98f00b204e9800998ecf8427e,e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
//
// # YARA Rule Generation
//
// DirHash can generate YARA rules for automated detection:
//
// Standard Rules (hashes + filenames):
//
//	rule malware_detection {
//	    meta:
//	        description = "Generated rule based on 2 files"
//	        author = "DirHash"
//	        date = "2023-12-01"
//	        tags = "generated, dirhash"
//	    strings:
//	        $md5_malware = { D4 1D 8C D9 8F 00 B2 04 }
//	        $filename_malware = "malware.exe"
//	    condition:
//	        $md5_malware or $filename_malware
//	}
//
// Hash-Only Rules (cryptographic hashes only):
//
//	rule hash_detection {
//	    meta:
//	        description = "Hash-based rule for 1 files"
//	        author = "DirHash"
//	        date = "2023-12-01"
//	        tags = "hash, generated, dirhash"
//	    strings:
//	        $md5_sample = { AB CD EF 12 34 56 78 90 }
//	        $sha256_sample = { DE AD BE EF CA FE BA BE }
//	    condition:
//	        any of ($md5_sample, $sha256_sample)
//	}
//
// # Security Considerations
//
// Hash Algorithm Selection:
//   - MD5: Fast but cryptographically broken, suitable only for file identification
//   - SHA1: Legacy algorithm with known vulnerabilities, avoid for security applications
//   - SHA256: Current standard, recommended for most security use cases
//   - SHA512: Extended version with larger digest size, suitable for high-security requirements
//
// YARA Rule Security:
//   - All rule names and identifiers are sanitized for YARA compliance
//   - Generated rules include metadata for attribution and tracking
//   - Hash values are properly formatted as hex patterns
//   - Condition logic is optimized for performance and accuracy
//
// # Integration
//
// DirHash integrates well with security tools and workflows:
//   - YARA command-line scanner
//   - VirusTotal intelligence platform
//   - SIEM and SOAR platforms
//   - Custom analysis pipelines
//   - Incident response automation
//
// # Documentation
//
// Complete package documentation is available via godoc:
//
//	godoc -http=:8080
//	# Visit http://localhost:8080/pkg/github.com/melatonein5/DirHash/
//
// For the latest information, visit:
// https://github.com/melatonein5/DirHash
package main
