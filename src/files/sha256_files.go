package files

import (
	"crypto/sha256"
	"fmt"
	"os"
)

// SHA256Files will hash all files in a directory using the SHA256 algorithm
func SHA256Files(files []File) ([]File, error) {
	// Declare a hasher
	hasher := sha256.New()
	// Iterate over the files
	for i, file := range files {
		// Read the file as bytes to hash it
		data, err := os.ReadFile(file.Path)
		if err != nil {
			return nil, err
		}
		// Write the data to the hasher
		hasher.Write(data)
		// Set the hash and hash type in the file struct
		files[i].Hash = fmt.Sprintf("%x", hasher.Sum(nil))
		// Set the hash type to sha256
		files[i].HashType = "sha256"
	}

	return files, nil
}
