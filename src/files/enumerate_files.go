package files

import (
	"os"
	"path/filepath"
)

// EnumnerateFiles enumerates all files in a directory and its subdirectories, and returns a slice of file paths.
func EnumerateFiles(dir string) ([]File, error) {
	var files []File

	//Traverse the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Check if the path is a file (not a directory)
		if !info.IsDir() {
			file := File{
				Path:     path,
				FileName: info.Name(),
				Hash:     "", // Placeholder for the file hash
				HashType: "", // Placeholder for the hash type
			}
			files = append(files, file)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
