# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DirHash is a command-line tool that recursively traverses directories and generates cryptographic hash values for all files using various algorithms (MD5, SHA1, SHA256, SHA512). It outputs results either to terminal (formatted table) or CSV file format.

## Development Commands

### Testing
```bash
# Run all tests with coverage (as used in CI)
go test -coverprofile=coverage.out ./src/...

# Run all tests except problematic main integration tests
go test ./src/...

# Run specific package tests
go test ./src/args
go test ./src/files
go test ./src/cmdline
go test ./src/yara

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

**Note**: The main integration tests (`main_integration_test.go`) are excluded from CI testing due to test flag compatibility issues. The GitHub Actions workflow uses `go test ./src/...` which focuses on package-specific tests.

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

**Note**: The installers embed the `dirhash` binary using `//go:embed dirhash`. The installer directories are excluded from testing in the main project as they require the binary to be built first and copied to their respective directories.

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

## Features

### ✅ Implemented
- **Multi-algorithm hashing**: MD5, SHA1, SHA256, SHA512 with concurrent processing
- **Flexible output formats**: Standard, Condensed, IOC-friendly for different use cases
- **YARA rule generation**: Automatic rule creation for malware detection with hash and filename patterns
- **KQL query generation**: Microsoft security platform integration (Microsoft 365 Defender, Azure Sentinel, Log Analytics)
- **Comprehensive testing**: 91.7% code coverage with extensive unit and integration tests
- **Cross-platform support**: Linux and Windows binaries with automated installers
- **Security-focused architecture**: Designed for threat hunting and security analysis workflows
- **Modular codebase**: Clean separation of concerns with well-documented APIs

### 🔮 Roadmap
Future enhancements planned:
1. **Sigma rule generation**: Support for SIEM platform rule creation
2. **JSON output format**: API-friendly structured output for integration
3. **Recursive subdirectory exclusion**: Pattern-based directory filtering
4. **Hash verification and integrity checking**: Validate file integrity over time
5. **Performance optimizations**: Enhanced processing for large file sets
6. **Custom hash algorithms**: Support for additional cryptographic functions
7. **Database export**: Direct export to security databases and threat intel platforms
8. **REST API mode**: HTTP API for integration with security orchestration platforms