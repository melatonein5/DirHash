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
