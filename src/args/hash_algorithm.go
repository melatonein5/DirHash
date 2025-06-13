package args

import "errors"

// Translates the string hash algorithm to an int
func StrHashAlgorithmToId(strHashAlgorithm string) int {
	switch strHashAlgorithm {
	case "md5":
		return 0
	case "sha1":
		return 1
	case "sha256":
		return 2
	case "sha512":
		return 3
	default:
		return -1 // Invalid hash algorithm
	}
}

// HashAlgorithmValidation will return an error if the hash algorithm is not valid
func HashAlgorithmValidation(id int) error {
	// Consider changing this to a check for -1, although this could be corrupted by a bit flip (unlikely)
	if id < 0 || id > 3 {
		return errors.New("invalid hash algorithm. argument must be one of: md5, sha1, sha256, sha512")
	}
	return nil
}
