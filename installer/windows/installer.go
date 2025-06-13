// windows_installer/main.go
package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

//go:embed dirhash.exe
var dirhashBinary []byte

// The installation directory for our application.
var installDir = filepath.Join(os.Getenv("ProgramFiles"), "DirHash")

func main() {
	//force import "embed" package
	_ = embed.FS{}

	fmt.Println("Starting the installation process for dirhash...")

	// 1. Check if the program is running with Administrator privileges.
	if !isAdmin() {
		fmt.Println("Error: Please run this installer as an Administrator.")
		// Keep the window open for a moment so the user can read the message.
		fmt.Println("Press Enter to exit.")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("Administrator privileges detected.")

	// 2. Create the installation directory.
	// os.MkdirAll is safe to run even if the directory already exists.
	err := os.MkdirAll(installDir, 0755)
	if err != nil {
		fmt.Printf("Error: Failed to create installation directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Installation directory '%s' is ready.\n", installDir)

	// 3. Write the embedded .exe to the installation directory.
	installPath := filepath.Join(installDir, "dirhash.exe")
	err = os.WriteFile(installPath, dirhashBinary, 0755)
	if err != nil {
		fmt.Printf("Error: Failed to write binary to disk: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully installed dirhash.exe to %s\n", installPath)

	// 4. Add the installation directory to the system PATH environment variable.
	err = addToSystemPath(installDir)
	if err != nil {
		fmt.Printf("Error: Failed to add installation directory to system PATH: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Success! DirHash has been installed.")
	fmt.Println("IMPORTANT: You must open a new terminal window to use the 'dirhash' command.")
	fmt.Println("\nPress Enter to exit.")
	fmt.Scanln()
}

// isAdmin checks if the current process has administrator privileges.
func isAdmin() bool {
	var sid *windows.SID
	// The RID for the local administrators group is 544.
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

// addToSystemPath adds the given directory to the system's PATH environment variable.
func addToSystemPath(newPath string) error {
	// Open the registry key for system environment variables.
	// KEY_READ and KEY_WRITE are combined to request both permissions.
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ|registry.WRITE)
	if err != nil {
		return err
	}
	defer key.Close()

	// Get the current value of the PATH variable.
	currentPath, _, err := key.GetStringValue("Path")
	if err != nil {
		return err
	}

	// Check if our path is already in the system PATH.
	for _, p := range filepath.SplitList(currentPath) {
		if strings.EqualFold(p, newPath) {
			fmt.Println("Installation directory is already in the system PATH.")
			return nil
		}
	}

	// If not present, append our new path.
	newSystemPath := currentPath + ";" + newPath
	err = key.SetStringValue("Path", newSystemPath)
	if err != nil {
		return err
	}

	fmt.Println("Successfully added installation directory to system PATH.")
	return nil
}
