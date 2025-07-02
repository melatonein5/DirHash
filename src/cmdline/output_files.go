package cmdline

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/melatonein5/DirHash/src/files"
)

// OutputFiles takes a slice of File structs and outputs them to the terminal
func OutputFiles(fileList []*files.File) {
	if len(fileList) == 0 {
		fmt.Println("No files to display")
		return
	}

	// Create a tab writer to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Print the header
	fmt.Fprintln(w, "File Name\tPath\tSize\tHash\tHash Type")

	// For each file, create a row for each hash type
	for _, f := range fileList {
		if len(f.Hashes) == 0 {
			// No hashes available
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", f.FileName, f.Path, f.Size, "N/A", "N/A")
		} else {
			// Sort hash types for consistent output
			var hashTypes []string
			for hashType := range f.Hashes {
				hashTypes = append(hashTypes, hashType)
			}
			sort.Strings(hashTypes)

			// Print a row for each hash type
			for _, hashType := range hashTypes {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", f.FileName, f.Path, f.Size, f.Hashes[hashType], hashType)
			}
		}
	}

	// Flush the writer to ensure all output is printed
	w.Flush()
}

// OutputFilesCondensed provides a more compact view with all hashes on one line
func OutputFilesCondensed(fileList []*files.File) {
	if len(fileList) == 0 {
		fmt.Println("No files to display")
		return
	}

	// Create a tab writer to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Print the header
	fmt.Fprintln(w, "File Name\tPath\tSize\tHashes")

	for _, f := range fileList {
		if len(f.Hashes) == 0 {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", f.FileName, f.Path, f.Size, "N/A")
		} else {
			// Sort hash types for consistent output
			var hashTypes []string
			for hashType := range f.Hashes {
				hashTypes = append(hashTypes, hashType)
			}
			sort.Strings(hashTypes)

			// Create hash string
			var hashStrings []string
			for _, hashType := range hashTypes {
				hashStrings = append(hashStrings, fmt.Sprintf("%s:%s", hashType, f.Hashes[hashType]))
			}
			
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", f.FileName, f.Path, f.Size, strings.Join(hashStrings, " | "))
		}
	}

	// Flush the writer to ensure all output is printed
	w.Flush()
}

// OutputFilesIOC provides IOC-friendly terminal output format
func OutputFilesIOC(fileList []*files.File) {
	if len(fileList) == 0 {
		fmt.Println("No files to display")
		return
	}

	// Create a tab writer to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	// Print the header with IOC-friendly column names
	fmt.Fprintln(w, "File Path\tFile Name\tSize\tMD5\tSHA1\tSHA256\tSHA512")

	for _, f := range fileList {
		// Extract hash values or use "N/A" if not available
		md5Hash := getHashOrNA(f.Hashes, "md5")
		sha1Hash := getHashOrNA(f.Hashes, "sha1")
		sha256Hash := getHashOrNA(f.Hashes, "sha256")
		sha512Hash := getHashOrNA(f.Hashes, "sha512")

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n", 
			f.Path, f.FileName, f.Size, md5Hash, sha1Hash, sha256Hash, sha512Hash)
	}

	// Flush the writer to ensure all output is printed
	w.Flush()
}

// getHashOrNA returns the hash value or "N/A" if not present
func getHashOrNA(hashes map[string]string, hashType string) string {
	if hash, exists := hashes[hashType]; exists {
		return hash
	}
	return "N/A"
}
