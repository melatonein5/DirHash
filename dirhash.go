package main

import (
	"log"
	"os"

	"github.com/melatonein5/DirHash/src/args"
	"github.com/melatonein5/DirHash/src/cmdline"
	"github.com/melatonein5/DirHash/src/files"
)

// args is a global struct that holds the arguments provided in the command line
var arguments args.Args

// Init parses the args provided
func init() {
	//First, grab the args from the command line, ignoring the first argument which is the program name, then parse them
	rawArgs := os.Args[1:]
	var err error
	arguments, err = args.ParseArgs(rawArgs)
	if err != nil {
		// Handle error

		log.Fatal(err)
	}

	if arguments.Help {
		cmdline.PrintHelp()
		// Exit the program after printing help
		os.Exit(0)
	}
}

func main() {
	//Firstly, enumerate the files in the input directory
	fs, err := files.EnumerateFiles(arguments.StrInputDir)
	if err != nil {
		log.Fatalf("Error enumerating files: %v", err)
	}

	//Then, hash the files using the specified algorithm
	hashedFiles, err := files.HashFiles(fs, arguments.HashAlgorithmId)
	if err != nil {
		log.Fatalf("Error hashing files: %v", err)
	}

	//Check if the output should be written to a file or printed to the terminal
	if arguments.OutputToTerminal {
		cmdline.OutputFiles(hashedFiles)
	}

	if arguments.WriteToFile {
		//Write the files to the output file
		if err := files.WriteOutput(hashedFiles, arguments.StrOutputFile); err != nil {
			log.Fatalf("Error writing files to output file: %v", err)
		}
	}
}
