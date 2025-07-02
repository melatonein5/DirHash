# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DirHash is a command-line tool that recursively traverses directories and generates cryptographic hash values for all files using various algorithms (MD5, SHA1, SHA256, SHA512). It outputs results either to terminal (formatted table) or CSV file format.

## Development Commands

### Running the Application
```bash
# Run with Go runtime (development)
go run dirhash.go [options]

# Build executable
go build dirhash.go

# Run built executable (Linux)
./dirhash [options]

# Run built executable (Windows)  
dirhash.exe [options]
```

### Common Usage Examples
```bash
# Hash current directory with MD5 (default), output to terminal
go run dirhash.go -t

# Hash specific directory with multiple algorithms, save to file
go run dirhash.go -i /path/to/dir -o output.txt -a sha256 sha512

# Show help
go run dirhash.go --help
```

### Building Installers
The project includes platform-specific installers in the `installer/` directory:
- Linux installer: `installer/linux/installer.go`
- Windows installer: `installer/windows/installer.go`

## Architecture Overview

The codebase follows a clean, modular pipeline architecture:

**Data Flow**: Command Line Args → File Enumeration → Hash Computation → Output Generation

### Core Packages

1. **`src/args/`** - Command-line argument parsing and validation
   - `args.go`: Defines Args struct with configuration options
   - `parse_args.go`: Robust argument parsing with error handling
   - `hash_algorithm.go`: Maps algorithm names to IDs (0=MD5, 1=SHA1, 2=SHA256, 3=SHA512)

2. **`src/files/`** - Core file processing functionality  
   - `file.go`: File struct definition with JSON serialization
   - `enumerate_files.go`: Recursive directory traversal
   - `hash_files.go`: Coordinates hash computation across algorithms
   - Individual hash implementations: `md5_files.go`, `sha1_files.go`, `sha256_files.go`, `sha512_files.go`
   - `write_output.go`: CSV output generation

3. **`src/cmdline/`** - User interface and output formatting
   - `help.go`: Comprehensive help text and usage examples
   - `output_files.go`: Terminal output formatting with tabwriter

### Key Design Patterns

- **Pipeline Architecture**: Each stage processes data and passes to next stage
- **Algorithm Abstraction**: Hash algorithms identified by integer IDs with individual implementations
- **Graceful Degradation**: Individual file errors don't stop entire process
- **Centralized Configuration**: All options stored in Args struct
- **Error Recovery**: Comprehensive error handling at each stage

## Code Conventions

- Use `log.Fatalf()` for fatal errors with descriptive messages
- Algorithm implementations follow consistent interface pattern
- File structs include JSON tags for potential serialization
- Error handling: `if err != nil { log.Fatalf("Error description: %v", err) }`
- Global `arguments` variable holds parsed command-line configuration

## Extending the Codebase

### Adding New Hash Algorithms
1. Add algorithm constant to `src/args/hash_algorithm.go`
2. Create new `[algorithm]_files.go` in `src/files/`
3. Update switch statement in `src/files/hash_files.go`
4. Update help text in `src/cmdline/help.go`

### Modifying Output Formats
- Terminal output: Edit `src/cmdline/output_files.go`
- File output: Edit `src/files/write_output.go`

The modular design allows easy modification of individual components without affecting the overall architecture.