package files

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

// HashFiles will hash all files concurrently using the specified hashing algorithms
func HashFiles(files []*File, hashAlgos []int) ([]*File, error) {
	if len(files) == 0 {
		return files, nil
	}

	// Create algorithm name lookup
	algoNames := make(map[int]string)
	for _, algo := range GetSupportedAlgorithms() {
		algoNames[algo.ID] = algo.Name
	}

	// Filter valid algorithms
	var validAlgos []int
	var algoNamesList []string
	for _, algo := range hashAlgos {
		if name, exists := algoNames[algo]; exists {
			validAlgos = append(validAlgos, algo)
			algoNamesList = append(algoNamesList, name)
		} else {
			log.Printf("Unsupported hash algorithm ID: %d", algo)
		}
	}

	if len(validAlgos) == 0 {
		return files, fmt.Errorf("no valid hash algorithms provided")
	}

	// Use worker pool for concurrent processing
	numWorkers := runtime.NumCPU()
	if numWorkers > len(files) {
		numWorkers = len(files)
	}

	// Create channels for work distribution
	fileChan := make(chan *File, len(files))
	resultChan := make(chan *File, len(files))
	errorChan := make(chan error, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go hashWorker(fileChan, resultChan, errorChan, validAlgos, algoNamesList, &wg)
	}

	// Send files to workers
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Process results
	var result []*File
	var errors []error

	for {
		select {
		case file, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				result = append(result, file)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else if err != nil {
				errors = append(errors, err)
			}
		}

		if resultChan == nil && errorChan == nil {
			break
		}
	}

	// Return first error if any occurred
	if len(errors) > 0 {
		return result, errors[0]
	}

	return result, nil
}

// hashWorker processes files from the channel and calculates hashes
func hashWorker(fileChan <-chan *File, resultChan chan<- *File, errorChan chan<- error, algorithms []int, algoNames []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range fileChan {
		if err := calculateAllHashes(file, algorithms, algoNames); err != nil {
			log.Printf("Error hashing file %s: %v", file.Path, err)
			errorChan <- err
			continue
		}
		resultChan <- file
	}
}

// calculateAllHashes reads the file once and calculates all required hashes
func calculateAllHashes(file *File, algorithms []int, algoNames []string) error {
	// Open file
	f, err := os.Open(file.Path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file.Path, err)
	}
	defer f.Close()

	// Create hashers for all algorithms
	hashers := make([]hash.Hash, len(algorithms))
	writers := make([]io.Writer, len(algorithms))

	for i, algo := range algorithms {
		switch algo {
		case 0: // MD5
			hashers[i] = md5.New()
		case 1: // SHA1
			hashers[i] = sha1.New()
		case 2: // SHA256
			hashers[i] = sha256.New()
		case 3: // SHA512
			hashers[i] = sha512.New()
		default:
			return fmt.Errorf("unsupported algorithm: %d", algo)
		}
		writers[i] = hashers[i]
	}

	// Create multi-writer to write to all hashers simultaneously
	multiWriter := io.MultiWriter(writers...)

	// Copy file content to all hashers at once
	if _, err := io.Copy(multiWriter, f); err != nil {
		return fmt.Errorf("failed to read file %s: %w", file.Path, err)
	}

	// Extract hash values
	for i, hasher := range hashers {
		file.Hashes[algoNames[i]] = fmt.Sprintf("%x", hasher.Sum(nil))
	}

	return nil
}
