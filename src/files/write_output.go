package files

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

// WriteOutput will write the output to a file as a CSV with expanded format
func WriteOutput(files []*File, outputPath string) error {
	// Open the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Path", "FileName", "Size", "Hash", "HashType"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write the file data - one row per hash type per file
	for _, f := range files {
		if len(f.Hashes) == 0 {
			// No hashes available
			record := []string{f.Path, f.FileName, strconv.FormatInt(f.Size, 10), "N/A", "N/A"}
			if err := writer.Write(record); err != nil {
				return err
			}
		} else {
			// Sort hash types for consistent output
			var hashTypes []string
			for hashType := range f.Hashes {
				hashTypes = append(hashTypes, hashType)
			}
			sort.Strings(hashTypes)

			// Write a row for each hash type
			for _, hashType := range hashTypes {
				record := []string{
					f.Path,
					f.FileName,
					strconv.FormatInt(f.Size, 10),
					f.Hashes[hashType],
					hashType,
				}
				if err := writer.Write(record); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// WriteOutputCondensed writes all hashes for each file on a single row
func WriteOutputCondensed(files []*File, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Determine all unique hash types across all files
	hashTypeSet := make(map[string]bool)
	for _, f := range files {
		for hashType := range f.Hashes {
			hashTypeSet[hashType] = true
		}
	}

	// Sort hash types for consistent column order
	var hashTypes []string
	for hashType := range hashTypeSet {
		hashTypes = append(hashTypes, hashType)
	}
	sort.Strings(hashTypes)

	// Create header
	header := []string{"Path", "FileName", "Size"}
	for _, hashType := range hashTypes {
		header = append(header, hashType)
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write file data
	for _, f := range files {
		record := []string{f.Path, f.FileName, strconv.FormatInt(f.Size, 10)}
		
		// Add hash values in order
		for _, hashType := range hashTypes {
			if hash, exists := f.Hashes[hashType]; exists {
				record = append(record, hash)
			} else {
				record = append(record, "")
			}
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// WriteOutputForIOC writes output in a format suitable for IOC/YARA generation
func WriteOutputForIOC(files []*File, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header optimized for IOC tools
	header := []string{"file_path", "file_name", "file_size", "md5", "sha1", "sha256", "sha512"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write file data in IOC format
	for _, f := range files {
		record := []string{
			f.Path,
			f.FileName,
			strconv.FormatInt(f.Size, 10),
			getHashOrEmpty(f.Hashes, "md5"),
			getHashOrEmpty(f.Hashes, "sha1"),
			getHashOrEmpty(f.Hashes, "sha256"),
			getHashOrEmpty(f.Hashes, "sha512"),
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// getHashOrEmpty returns the hash value or empty string if not present
func getHashOrEmpty(hashes map[string]string, hashType string) string {
	if hash, exists := hashes[hashType]; exists {
		return hash
	}
	return ""
}
