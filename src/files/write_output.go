package files

import (
	"encoding/csv"
	"os"
)

// WriteOutput will write the output to a file as a CSV
func WriteOutput(files []File, outputPath string) error {
	// Open the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Create a CSV writer
	writer := csv.NewWriter(file)
	// Write the header
	header := []string{"Path", "Hash", "HashType"}
	if err := writer.Write(header); err != nil {
		return err
	}
	// Write the file data
	for _, f := range files {
		record := []string{f.Path, f.Hash, f.HashType}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	// Flush the writer to ensure all data is written
	writer.Flush()
	return nil
}
