package main

import (
	"log"
	"os"

	"github.com/melatonein5/DirHash/src/args"
	"github.com/melatonein5/DirHash/src/cmdline"
	"github.com/melatonein5/DirHash/src/files"
)

// arguments is a global struct that holds the arguments provided in the command line
var arguments args.Args

// init parses the args provided
func init() {
	// First, grab the args from the command line, ignoring the first argument which is the program name, then parse them
	rawArgs := os.Args[1:]
	var err error
	arguments, err = args.ParseArgs(rawArgs)
	if err != nil {
		log.Fatal(err)
	}

	if arguments.Help {
		cmdline.PrintHelp()
		// Exit the program after printing help
		os.Exit(0)
	}
}

func main() {
	// Enumerate the files in the input directory
	fs, err := files.EnumerateFiles(arguments.StrInputDir)
	if err != nil {
		log.Fatalf("Error enumerating files: %v", err)
	}

	log.Printf("Found %d files to process", len(fs))

	// Hash the files using the specified algorithms concurrently
	hashedFiles, err := files.HashFiles(fs, arguments.HashAlgorithmId)
	if err != nil {
		log.Fatalf("Error hashing files: %v", err)
	}

	log.Printf("Successfully processed %d files", len(hashedFiles))

	// Check if the output should be written to a file or printed to the terminal
	if arguments.OutputToTerminal {
		switch arguments.OutputFormat {
		case "condensed":
			cmdline.OutputFilesCondensed(hashedFiles)
		case "ioc":
			cmdline.OutputFilesIOC(hashedFiles)
		default: // "standard"
			cmdline.OutputFiles(hashedFiles)
		}
	}

	if arguments.WriteToFile {
		// Write the files to the output file using the specified format
		var err error
		switch arguments.OutputFormat {
		case "condensed":
			err = files.WriteOutputCondensed(hashedFiles, arguments.StrOutputFile)
		case "ioc":
			err = files.WriteOutputForIOC(hashedFiles, arguments.StrOutputFile)
		default: // "standard"
			err = files.WriteOutput(hashedFiles, arguments.StrOutputFile)
		}

		if err != nil {
			log.Fatalf("Error writing files to output file: %v", err)
		}
		log.Printf("Output written to: %s (format: %s)", arguments.StrOutputFile, arguments.OutputFormat)
	}
}
