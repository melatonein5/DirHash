// installer/main.go
package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed dirhash
var dirhashBinary []byte

const installPath = "/usr/local/bin/dirhash"

func main() {
	//This is to shut the linter up, it is not used in this file.
	garbage := embed.FS{}
	_ = garbage

	// Print a message to indicate the start of the installation process.
	fmt.Println("Starting the installation process for dirhash...")

	// 1. Check for superuser privileges.
	// On Linux, a user ID of 0 indicates the root user.
	if os.Geteuid() != 0 {
		fmt.Println("Error: This installer must be run with superuser privileges. Please use sudo.")
		os.Exit(1)
	}

	fmt.Println("Superuser privileges detected.")

	// 2. Write the embedded binary to the installation path.
	// The permissions 0755 mean:
	// - The owner can read, write, and execute.
	// - Group members and others can only read and execute.
	err := os.WriteFile(installPath, dirhashBinary, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to install binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed dirhash to %s\n", installPath)
	fmt.Println("You can now run 'dirhash' from your terminal.")
}
