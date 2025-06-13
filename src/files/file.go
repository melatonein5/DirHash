package files

// File is a struct that represents a file with its path and hash
type File struct {
	//FileName is the name of the file
	FileName string `json:"filename"`
	//Path is the path to the file
	Path string `json:"path"`
	//Hash is the hash of the file
	Hash string `json:"hash"`
	//Get the HashType of the file
	HashType string `json:"hash_type"`
}
