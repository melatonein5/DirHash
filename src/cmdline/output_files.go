package cmdline

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/melatonein5/DirHash/src/files"
)

// OutputFiles takes a slice of File structs and outputs them to the terminal or a file
func OutputFiles(files []files.File) {
	//Create a tab writer to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	// Print the header
	fmt.Fprintln(w, "File Name\tPath\tHash\tHash Type")

	for _, f := range files {
		// Print the file name, path, hash, and hash type
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", f.FileName, f.Path, f.Hash, f.HashType)
	}
	// Flush the writer to ensure all output is printed
	w.Flush()
}
