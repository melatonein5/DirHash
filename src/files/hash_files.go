package files

import "log"

// HashFiles will hash all files in a directory using the specified hashing algorithm
func HashFiles(files []File, hashAlgos []int) ([]File, error) {
	//As there can be multiple hash algorithms, we will iterate over them
	//First, create the result slice
	var result []File
	for _, algo := range hashAlgos {
		var hashedFiles []File
		var err error

		//Switch on the algorithm
		switch algo {
		case 0: // MD5
			hashedFiles, err = MD5Files(files)
		case 1: // SHA1
			hashedFiles, err = SHA1Files(files)
		case 2: // SHA256
			hashedFiles, err = SHA256Files(files)
		case 3: // SHA512
			hashedFiles, err = SHA512Files(files)
		default:
			log.Println("Unsupported hash algorithm")
			//Skip unsupported algorithms
			continue
		}

		if err != nil {
			return nil, err
		}

		result = append(result, hashedFiles...)
	}
	return result, nil
}
