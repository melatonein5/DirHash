// Package files provides core file processing functionality for DirHash.
//
// This package handles file discovery, hash computation, and output generation
// for the DirHash application. It provides the fundamental building blocks for
// processing files and generating cryptographic hashes at scale.
//
// # Core Functionality
//
//   - File Discovery: Recursively enumerate files in directory trees
//   - Hash Computation: Generate MD5, SHA1, SHA256, SHA512 hashes concurrently
//   - Output Generation: Export results in multiple formats (standard, condensed, IOC)
//   - Metadata Collection: Capture file size, modification time, and path information
//
// # File Processing Pipeline
//
// The package implements a three-stage processing pipeline:
//
//  1. Enumeration: Discover all files in the specified directory tree
//  2. Hashing: Compute cryptographic hashes using specified algorithms
//  3. Output: Format and export results to terminal and/or files
//
// # Concurrency and Performance
//
// File processing is optimized for performance through:
//   - Concurrent hash computation across multiple files
//   - Efficient memory usage for large file sets
//   - Streaming I/O to handle files of any size
//   - Minimal memory allocation during processing
//
// # Supported Hash Algorithms
//
// The package supports industry-standard cryptographic hash functions:
//   - MD5: Fast legacy algorithm for file identification
//   - SHA1: Legacy algorithm still used in some contexts
//   - SHA256: Modern standard for cryptographic hashing
//   - SHA512: Extended version with larger digest size
//
// # Output Formats
//
//   - Standard: Traditional format with separate rows per hash type
//   - Condensed: All hashes for a file on a single row
//   - IOC: Indicator of Compromise format for security analysis tools
//
// # Usage Example
//
//	// Enumerate files
//	files, err := files.EnumerateFiles("/path/to/directory")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Generate hashes
//	hashAlgorithms := []int{0, 2} // MD5 and SHA256
//	hashedFiles, err := files.HashFiles(files, hashAlgorithms)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Write output
//	err = files.WriteOutput(hashedFiles, "output.csv")
//	if err != nil {
//		log.Fatal(err)
//	}
package files

import (
	"os"
	"time"
)

// File represents a file with its metadata and computed cryptographic hash values.
//
// This structure serves as the primary data container for file information
// throughout the DirHash processing pipeline. It stores both filesystem
// metadata and computed hash values for each processed file.
//
// The File struct supports JSON serialization for potential future API
// or configuration file use cases.
type File struct {
	FileName string            `json:"filename"` // Base filename without directory path
	Path     string            `json:"path"`     // Full filesystem path to the file
	Size     int64             `json:"size"`     // File size in bytes
	ModTime  time.Time         `json:"mod_time"` // Last modification timestamp
	Hashes   map[string]string `json:"hashes"`   // Hash values keyed by algorithm name (e.g., "md5", "sha256")
}

// NewFile creates a new File struct with initialized fields from filesystem information.
//
// This constructor function creates a File instance with metadata populated from
// os.FileInfo and an initialized (empty) hash map ready for hash computation.
//
// Parameters:
//   - path: Full filesystem path to the file
//   - fileName: Base filename (typically from filepath.Base)
//   - fileInfo: os.FileInfo containing filesystem metadata
//
// Returns:
//   - *File: Initialized File struct ready for hash computation
//
// The returned File will have an empty Hashes map that should be populated
// by calling HashFiles() or individual hash computation functions.
func NewFile(path, fileName string, fileInfo os.FileInfo) *File {
	return &File{
		FileName: fileName,
		Path:     path,
		Size:     fileInfo.Size(),
		ModTime:  fileInfo.ModTime(),
		Hashes:   make(map[string]string),
	}
}

// HashAlgorithm represents a supported cryptographic hash algorithm.
//
// This structure defines the mapping between human-readable algorithm names
// and internal numeric identifiers used throughout the application.
type HashAlgorithm struct {
	ID   int    // Internal numeric identifier for the algorithm
	Name string // Human-readable algorithm name (e.g., "md5", "sha256")
}

// GetSupportedAlgorithms returns all cryptographic hash algorithms supported by DirHash.
//
// This function provides the canonical list of supported hash algorithms with their
// corresponding internal IDs. The IDs are used throughout the application for
// efficient algorithm identification and processing.
//
// Returns:
//   - []HashAlgorithm: Slice containing all supported algorithms with their IDs and names
//
// Supported algorithms:
//   - ID 0: MD5 (fast, legacy, suitable for file identification)
//   - ID 1: SHA1 (legacy, still used in some security contexts)
//   - ID 2: SHA256 (modern standard, recommended for most use cases)
//   - ID 3: SHA512 (extended version with larger digest size)
func GetSupportedAlgorithms() []HashAlgorithm {
	return []HashAlgorithm{
		{ID: 0, Name: "md5"},
		{ID: 1, Name: "sha1"},
		{ID: 2, Name: "sha256"},
		{ID: 3, Name: "sha512"},
	}
}
