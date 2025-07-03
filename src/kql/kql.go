// Package kql provides KQL (Kusto Query Language) query generation from file hash data.
//
// This package enables the generation of KQL queries suitable for Microsoft Sentinel,
// Azure Log Analytics, Microsoft 365 Defender, and other platforms that support
// Kusto Query Language. The generated queries can be used for threat hunting,
// security analysis, and incident response workflows.
//
// # KQL Query Types
//
// The package supports generating different types of KQL queries:
//   - Hash-based queries: Search for specific file hashes in security logs
//   - Filename-based queries: Search for specific filenames in security logs
//   - Combined queries: Search for both hashes and filenames with logical operators
//   - Multi-table queries: Generate queries that search across multiple log tables
//
// # Supported Log Sources
//
// The generated KQL queries are designed to work with common security log sources:
//   - DeviceFileEvents (Microsoft 365 Defender)
//   - SecurityEvents (Azure Security Center)
//   - CommonSecurityLog (Azure Sentinel)
//   - FileHashAlgorithm tables (custom security logs)
//
// # Usage Example
//
//	files := []*files.File{...} // File data with hashes
//	query, err := kql.GenerateKQLQuery(files, "ThreatHunt_Malware", []string{"md5", "sha256"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(query.ToKQLFormat())
//
// # Query Structure
//
// Generated KQL queries follow security industry best practices:
//   - Use proper KQL syntax and operators
//   - Include metadata comments for documentation
//   - Support time range filtering
//   - Include result sorting and limiting
//   - Provide clear field selection
//
// # Performance Considerations
//
// The package optimizes query generation for performance:
//   - Uses efficient KQL operators (in, contains, has)
//   - Generates queries with proper indexing considerations
//   - Supports batch processing for large file sets
//   - Provides options for query complexity management
package kql

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/melatonein5/DirHash/src/files"
)

// KQLQuery represents a generated KQL query with metadata.
//
// This structure contains the complete KQL query along with associated metadata
// that can be used for documentation, automation, and integration with security platforms.
type KQLQuery struct {
	// Query metadata
	Name        string    // Human-readable name for the query
	Description string    // Description of the query purpose
	Author      string    // Author information (default: "DirHash")
	Generated   time.Time // Timestamp when query was generated
	Tags        []string  // Tags for categorization and search

	// Query configuration
	Tables     []string // Log tables to search (e.g., "DeviceFileEvents", "SecurityEvents")
	HashTypes  []string // Hash algorithms used in the query (e.g., "md5", "sha256")
	TimeRange  string   // Time range for the query (e.g., "7d", "30d")
	MaxResults int      // Maximum number of results to return

	// Query components
	HashList     []string // List of hashes to search for
	FilenameList []string // List of filenames to search for
	QueryBody    string   // The main KQL query body
	Comments     []string // Additional comments for the query
}

// KQLQueryOptions configures KQL query generation.
//
// This structure provides fine-grained control over how KQL queries are generated,
// allowing customization for different security platforms and use cases.
type KQLQueryOptions struct {
	// Query targeting options
	Tables    []string // Log tables to search (default: ["DeviceFileEvents"])
	HashTypes []string // Hash types to include (default: all available)
	TimeRange string   // Time range for query (default: "7d")

	// Query behavior options
	MaxResults      int  // Maximum results to return (default: 1000)
	IncludeHashes   bool // Include hash-based searches (default: true)
	IncludeFilenames bool // Include filename-based searches (default: true)
	CaseSensitive   bool // Case sensitive filename matching (default: false)

	// Output formatting options
	IncludeMetadata bool // Include metadata comments (default: true)
	IncludeComments bool // Include explanatory comments (default: true)
	CompactFormat   bool // Generate compact query format (default: false)
}

// DefaultKQLQueryOptions returns default options for KQL query generation.
//
// These defaults are optimized for common security analysis scenarios and
// provide a good starting point for most use cases.
func DefaultKQLQueryOptions() KQLQueryOptions {
	return KQLQueryOptions{
		Tables:           []string{"DeviceFileEvents"},
		TimeRange:        "7d",
		MaxResults:       1000,
		IncludeHashes:    true,
		IncludeFilenames: true,
		CaseSensitive:    false,
		IncludeMetadata:  true,
		IncludeComments:  true,
		CompactFormat:    false,
	}
}

// GenerateKQLQuery creates a KQL query from file hash data.
//
// This function generates a comprehensive KQL query that can be used to search
// for the provided files across various security log sources. The query includes
// both hash-based and filename-based detection logic.
//
// Parameters:
//   - files: Slice of File structures containing hash and filename data
//   - queryName: Human-readable name for the generated query
//   - hashTypes: Hash algorithms to include in the query (empty = all available)
//
// Returns:
//   - *KQLQuery: Generated query structure with metadata
//   - error: Error if query generation fails
//
// Example:
//
//	files := []*files.File{
//		{FileName: "malware.exe", Hashes: map[string]string{"md5": "abc123", "sha256": "def456"}},
//	}
//	query, err := GenerateKQLQuery(files, "MalwareHunt", []string{"md5", "sha256"})
func GenerateKQLQuery(files []*files.File, queryName string, hashTypes []string) (*KQLQuery, error) {
	return GenerateKQLQueryWithOptions(files, queryName, hashTypes, DefaultKQLQueryOptions())
}

// GenerateKQLQueryWithOptions creates a KQL query with custom options.
//
// This function provides full control over KQL query generation, allowing
// customization of all aspects of the query including target tables, time ranges,
// result limits, and formatting preferences.
//
// Parameters:
//   - files: Slice of File structures containing hash and filename data
//   - queryName: Human-readable name for the generated query
//   - hashTypes: Hash algorithms to include in the query (empty = all available)
//   - options: Configuration options for query generation
//
// Returns:
//   - *KQLQuery: Generated query structure with metadata
//   - error: Error if query generation fails
//
// Example:
//
//	options := KQLQueryOptions{
//		Tables:    []string{"DeviceFileEvents", "SecurityEvents"},
//		TimeRange: "30d",
//		MaxResults: 5000,
//	}
//	query, err := GenerateKQLQueryWithOptions(files, "ThreatHunt", []string{"sha256"}, options)
func GenerateKQLQueryWithOptions(files []*files.File, queryName string, hashTypes []string, options KQLQueryOptions) (*KQLQuery, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for KQL query generation")
	}

	if queryName == "" {
		queryName = "dirhash_generated_query"
	}

	// Sanitize query name for KQL compatibility
	queryName = sanitizeKQLName(queryName)

	// Collect hashes and filenames from files
	hashMap := make(map[string][]string)
	var filenames []string

	for _, file := range files {
		// Collect filenames if enabled
		if options.IncludeFilenames && file.FileName != "" {
			filenames = append(filenames, file.FileName)
		}

		// Collect hashes if enabled
		if options.IncludeHashes {
			for hashType, hashValue := range file.Hashes {
				// Filter by requested hash types
				if len(hashTypes) > 0 && !contains(hashTypes, hashType) {
					continue
				}
				if len(options.HashTypes) > 0 && !contains(options.HashTypes, hashType) {
					continue
				}
				hashMap[hashType] = append(hashMap[hashType], hashValue)
			}
		}
	}

	// Remove duplicates and sort
	filenames = removeDuplicatesAndSort(filenames)
	for hashType := range hashMap {
		hashMap[hashType] = removeDuplicatesAndSort(hashMap[hashType])
	}

	// Create query structure
	query := &KQLQuery{
		Name:        queryName,
		Description: fmt.Sprintf("KQL query to detect files based on hashes and filenames - Generated from %d files", len(files)),
		Author:      "DirHash",
		Generated:   time.Now(),
		Tags:        []string{"threat-hunting", "file-detection", "security", "dirhash"},
		Tables:      options.Tables,
		HashTypes:   getHashTypesFromMap(hashMap),
		TimeRange:   options.TimeRange,
		MaxResults:  options.MaxResults,
		FilenameList: filenames,
	}

	// Collect all hashes into a single list for metadata
	var allHashes []string
	for _, hashes := range hashMap {
		allHashes = append(allHashes, hashes...)
	}
	query.HashList = allHashes

	// Generate query body
	queryBody, err := buildKQLQueryBody(hashMap, filenames, options)
	if err != nil {
		return nil, fmt.Errorf("failed to build KQL query body: %v", err)
	}
	query.QueryBody = queryBody

	// Generate comments
	query.Comments = generateKQLComments(query, options)

	return query, nil
}

// GenerateKQLQueryHashOnly creates a KQL query using only hash values.
//
// This function generates a KQL query that focuses exclusively on hash-based
// detection, excluding filename-based searches. This is useful for scenarios
// where filenames may change but hash values remain constant.
//
// Parameters:
//   - files: Slice of File structures containing hash data
//   - queryName: Human-readable name for the generated query
//   - hashTypes: Hash algorithms to include in the query
//
// Returns:
//   - *KQLQuery: Generated hash-only query structure
//   - error: Error if query generation fails
//
// Example:
//
//	query, err := GenerateKQLQueryHashOnly(files, "HashOnlyDetection", []string{"sha256"})
func GenerateKQLQueryHashOnly(files []*files.File, queryName string, hashTypes []string) (*KQLQuery, error) {
	options := DefaultKQLQueryOptions()
	options.IncludeFilenames = false
	options.IncludeHashes = true
	options.HashTypes = hashTypes

	return GenerateKQLQueryWithOptions(files, queryName, hashTypes, options)
}

// ToKQLFormat returns the complete KQL query as a formatted string.
//
// This method generates the final KQL query string that can be executed
// in KQL-enabled platforms like Azure Sentinel, Microsoft 365 Defender,
// or Azure Log Analytics.
//
// The returned string includes:
//   - Metadata comments (if enabled)
//   - Query description and documentation
//   - The executable KQL query
//   - Result formatting and limits
//
// Returns:
//   - string: Complete formatted KQL query ready for execution
//
// Example:
//
//	query, _ := GenerateKQLQuery(files, "ThreatHunt", []string{"sha256"})
//	kqlString := query.ToKQLFormat()
//	fmt.Println(kqlString)
func (q *KQLQuery) ToKQLFormat() string {
	var parts []string

	// Add metadata comments
	if len(q.Comments) > 0 {
		parts = append(parts, strings.Join(q.Comments, "\n"))
		parts = append(parts, "")
	}

	// Add the main query
	parts = append(parts, q.QueryBody)

	return strings.Join(parts, "\n")
}

// buildKQLQueryBody constructs the main KQL query body.
func buildKQLQueryBody(hashMap map[string][]string, filenames []string, options KQLQueryOptions) (string, error) {
	var queryParts []string
	var unionParts []string

	// Build query for each table
	for _, table := range options.Tables {
		tableParts := []string{table}

		// Add time range filter
		if options.TimeRange != "" {
			tableParts = append(tableParts, fmt.Sprintf("| where TimeGenerated >= ago(%s)", options.TimeRange))
		}

		// Build conditions
		var conditions []string

		// Add hash conditions
		if options.IncludeHashes && len(hashMap) > 0 {
			var hashConditions []string
			for hashType, hashes := range hashMap {
				if len(hashes) > 0 {
					hashField := getHashFieldName(hashType, table)
					hashList := strings.Join(quoteStrings(hashes), ", ")
					hashConditions = append(hashConditions, fmt.Sprintf("(%s in (%s))", hashField, hashList))
				}
			}
			if len(hashConditions) > 0 {
				conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(hashConditions, " or ")))
			}
		}

		// Add filename conditions
		if options.IncludeFilenames && len(filenames) > 0 {
			filenameField := getFilenameFieldName(table)
			filenameList := strings.Join(quoteStrings(filenames), ", ")
			if options.CaseSensitive {
				conditions = append(conditions, fmt.Sprintf("(%s in (%s))", filenameField, filenameList))
			} else {
				conditions = append(conditions, fmt.Sprintf("(%s in~ (%s))", filenameField, filenameList))
			}
		}

		// Combine conditions
		if len(conditions) > 0 {
			tableParts = append(tableParts, fmt.Sprintf("| where %s", strings.Join(conditions, " or ")))
		}

		// Add field selection
		tableParts = append(tableParts, fmt.Sprintf("| project TimeGenerated, %s", getProjectFields(table)))

		// Add table identifier
		tableParts = append(tableParts, fmt.Sprintf("| extend SourceTable = \"%s\"", table))

		unionParts = append(unionParts, strings.Join(tableParts, "\n"))
	}

	// Combine all table queries
	if len(unionParts) > 1 {
		queryParts = append(queryParts, fmt.Sprintf("union (\n%s\n)", strings.Join(unionParts, "\n),\n(")))
	} else {
		queryParts = append(queryParts, unionParts[0])
	}

	// Add sorting and limiting
	queryParts = append(queryParts, "| sort by TimeGenerated desc")
	if options.MaxResults > 0 {
		queryParts = append(queryParts, fmt.Sprintf("| take %d", options.MaxResults))
	}

	return strings.Join(queryParts, "\n"), nil
}

// generateKQLComments creates documentation comments for the query.
func generateKQLComments(query *KQLQuery, options KQLQueryOptions) []string {
	var comments []string

	if !options.IncludeMetadata {
		return comments
	}

	comments = append(comments, fmt.Sprintf("// KQL Query: %s", query.Name))
	comments = append(comments, fmt.Sprintf("// Description: %s", query.Description))
	comments = append(comments, fmt.Sprintf("// Author: %s", query.Author))
	comments = append(comments, fmt.Sprintf("// Generated: %s", query.Generated.Format("2006-01-02 15:04:05 UTC")))
	comments = append(comments, fmt.Sprintf("// Tags: %s", strings.Join(query.Tags, ", ")))
	comments = append(comments, "//")

	if len(query.HashList) > 0 {
		comments = append(comments, fmt.Sprintf("// Hash Count: %d", len(query.HashList)))
		comments = append(comments, fmt.Sprintf("// Hash Types: %s", strings.Join(query.HashTypes, ", ")))
	}

	if len(query.FilenameList) > 0 {
		comments = append(comments, fmt.Sprintf("// Filename Count: %d", len(query.FilenameList)))
	}

	comments = append(comments, fmt.Sprintf("// Tables: %s", strings.Join(query.Tables, ", ")))
	comments = append(comments, fmt.Sprintf("// Time Range: %s", query.TimeRange))
	comments = append(comments, fmt.Sprintf("// Max Results: %d", query.MaxResults))

	if options.IncludeComments {
		comments = append(comments, "//")
		comments = append(comments, "// This query searches for files based on cryptographic hashes and filenames.")
		comments = append(comments, "// It can be used for threat hunting, incident response, and security analysis.")
		comments = append(comments, "// Modify the time range and result limits as needed for your environment.")
	}

	return comments
}

// Helper functions

// sanitizeKQLName sanitizes a string for use as a KQL identifier.
func sanitizeKQLName(name string) string {
	// Replace invalid characters with underscores
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)

	// Ensure it starts with a letter or underscore
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		result = "_" + result
	}

	return result
}

// contains checks if a string slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// removeDuplicatesAndSort removes duplicates from a string slice and sorts it.
func removeDuplicatesAndSort(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

// getHashTypesFromMap extracts and sorts hash types from a hash map.
func getHashTypesFromMap(hashMap map[string][]string) []string {
	var hashTypes []string
	for hashType := range hashMap {
		hashTypes = append(hashTypes, hashType)
	}
	sort.Strings(hashTypes)
	return hashTypes
}

// quoteStrings adds quotes around each string in a slice.
func quoteStrings(strings []string) []string {
	var quoted []string
	for _, s := range strings {
		quoted = append(quoted, fmt.Sprintf(`"%s"`, s))
	}
	return quoted
}

// getHashFieldName returns the appropriate hash field name for a given table.
func getHashFieldName(hashType, table string) string {
	switch table {
	case "DeviceFileEvents":
		switch hashType {
		case "md5":
			return "MD5"
		case "sha1":
			return "SHA1"
		case "sha256":
			return "SHA256"
		default:
			return "SHA256" // Default to SHA256 if unknown
		}
	case "SecurityEvents":
		return "FileHash"
	case "CommonSecurityLog":
		return "FileHash"
	default:
		return fmt.Sprintf("%sHash", strings.ToUpper(hashType))
	}
}

// getFilenameFieldName returns the appropriate filename field name for a given table.
func getFilenameFieldName(table string) string {
	switch table {
	case "DeviceFileEvents":
		return "FileName"
	case "SecurityEvents":
		return "FileName"
	case "CommonSecurityLog":
		return "FileName"
	default:
		return "FileName"
	}
}

// getProjectFields returns the appropriate project fields for a given table.
func getProjectFields(table string) string {
	switch table {
	case "DeviceFileEvents":
		return "DeviceName, FileName, FolderPath, MD5, SHA1, SHA256, ProcessCommandLine, InitiatingProcessFileName"
	case "SecurityEvents":
		return "Computer, FileName, FilePath, FileHash, ProcessName, CommandLine"
	case "CommonSecurityLog":
		return "Computer, FileName, FilePath, FileHash, ProcessName, CommandLine"
	default:
		return "Computer, FileName, FilePath, FileHash"
	}
}