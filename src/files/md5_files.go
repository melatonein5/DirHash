package files

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
)

// MD5Files will hash all files in a directory using the MD5 algorithm
func MD5Files(files []File) ([]File, error) {
	//Declare a hasher
	hasher := md5.New()
	//Iterate over the files
	for i, file := range files {
		//Read the file as bytes to hash it
		data, err := os.ReadFile(file.Path)
		if err != nil {
			log.Println("Error reading file:", file.Path, "-", err)
			continue // Skip this file and continue with the next one
		}
		//Write the data to the hasher
		hasher.Write(data)
		//Set the hash and hash type in the file struct
		files[i].Hash = fmt.Sprintf("%x", hasher.Sum(nil))
		//Set the hash type to md5
		files[i].HashType = "md5"
	}

	return files, nil
}
