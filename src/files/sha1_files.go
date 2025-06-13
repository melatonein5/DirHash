package files

import (
	"crypto/sha1"
	"fmt"
	"os"
)

// SHA1Files will hash all files in a directory using the SHA1 algorithm
func SHA1Files(files []File) ([]File, error) {
	// Declare a hasher
	hasher := sha1.New()
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
		// Set the hash type to sha1
		files[i].HashType = "sha1"
	}

	return files, nil
}
