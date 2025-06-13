package files

import "errors"

// HashFiles will hash all files in a directory using the specified hashing algorithm
func HashFiles(files []File, algorithm int) ([]File, error) {
	switch algorithm {
	case 0: // MD5
		return MD5Files(files)
	case 1: // SHA1
		return SHA1Files(files)
	case 2: // SHA256
		return SHA256Files(files)
	case 3: // SHA512
		return SHA512Files(files)
	default:
		return nil, errors.New("invalid hash algorithm. argument must be one of: md5, sha1, sha256, sha512")
	}
}
