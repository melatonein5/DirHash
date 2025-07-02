package files

import (
	"os"
	"time"
)

// File is a struct that represents a file with its metadata and multiple hash values
type File struct {
	// FileName is the name of the file
	FileName string `json:"filename"`
	// Path is the path to the file
	Path string `json:"path"`
	// Size is the file size in bytes
	Size int64 `json:"size"`
	// ModTime is the file modification time
	ModTime time.Time `json:"mod_time"`
	// Hashes contains all hash values keyed by algorithm name
	Hashes map[string]string `json:"hashes"`
}

// NewFile creates a new File struct with initialized fields
func NewFile(path, fileName string, fileInfo os.FileInfo) *File {
	return &File{
		FileName: fileName,
		Path:     path,
		Size:     fileInfo.Size(),
		ModTime:  fileInfo.ModTime(),
		Hashes:   make(map[string]string),
	}
}

// HashAlgorithm represents a hash algorithm configuration
type HashAlgorithm struct {
	ID   int
	Name string
}

// GetSupportedAlgorithms returns all supported hash algorithms
func GetSupportedAlgorithms() []HashAlgorithm {
	return []HashAlgorithm{
		{ID: 0, Name: "md5"},
		{ID: 1, Name: "sha1"},
		{ID: 2, Name: "sha256"},
		{ID: 3, Name: "sha512"},
	}
}
