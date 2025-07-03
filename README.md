# DirHash

[![Test Coverage](https://github.com/melatonein5/DirHash/actions/workflows/test-coverage.yml/badge.svg)](https://github.com/melatonein5/DirHash/actions/workflows/test-coverage.yml)
![Coverage](https://img.shields.io/badge/Coverage-91.7%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/melatonein5/DirHash)](https://goreportcard.com/report/github.com/melatonein5/DirHash)

DirHash is a command line tool to take a directory and return the file hashes.

## Installation
At the time of writing, an installer is only available for amd64 Linux. Download the installer from the [releases page](https://github.com/melatonein5/DirHash/releases/tag/latest), and run `sudo ./dirhash_amd64_linux_installer'`.

Alternatively, you can download the binary and use DirHash the following way:

#### Running through the pre-compiled binary
At the time of writing, installer and precomiled for AMD64 platforms running Linux and Windows. The following are ways to run DirHash on different systems:

```
Windows Binary:       dirhash.exe
Windows Installation: dirhash
Linux Binary:         ./dirhash
Linux Installation    dirhash
```

The following is generic usage using a Linux binary in the current directory:

```
Usage: ./dirhash [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format (default: standard)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Supported algorithms:
  md5, sha1, sha256, sha512
Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)
Examples:
  ./dirhash -i /path/to/dir -o output.csv -a sha256
  ./dirhash --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
  ./dirhash -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc
  ./dirhash -t
  ./dirhash --help
```

## Usage
Download DirFetch through the following command, and once DirHash is downloaded, move to the target directory using:
```
git clone https://github.com/melatonein5/DirHash.git
cd DirHash
```

#### Running through Go Runtime Environment (requires Go installation)
If a binary for your system is not available, you can run DirHash through the Go runtime environment with `go run dirhash.go`.
```
Usage: go run dirhash.go [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format (default: standard)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Supported algorithms:
  md5, sha1, sha256, sha512
Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)
Examples:
  go run dirhash.go -i /path/to/dir -o output.csv -a sha256
  go run dirhash.go --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
  go run dirhash.go -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc
  go run dirhash.go -t
  go run dirhash.go --help
```

#### Build from Source (requires Go installation)
To build from source, you can compile dirhash using `go build dirhash.go`, to compile an executable. The Go compiler selects your system architecture and operating system automatically. Note that if you have compiled to Windows, the `./dirhash` command becomes `dirhash.exe`.

## Output Formats
DirHash supports three output formats that apply to both terminal and file output:

### Standard Format (default)
Traditional format with separate rows for each hash type per file. Compatible with existing workflows.

### Condensed Format
All hash values for each file on a single row with dynamic columns based on requested algorithms. More compact for analysis.

### IOC Format
Security-optimized format with standardized columns perfect for threat intelligence tools, YARA rules, KQL queries, and Microsoft Sentinel IOC imports. Shows all hash types (MD5, SHA1, SHA256, SHA512) with "N/A" for missing values.

## Example Outputs
The format option (`-f`) affects both terminal display and file output, providing consistent formatting across all output methods.

##### Terminal Output Examples

**Standard Format** (`-f standard` or default):
```
$ go run dirhash.go -i src/args -t -a md5 sha256 -f standard
File Name         Path                       Size Hash                                                             Hash Type
args.go           src/args/args.go           314  4b140b31128a0eabb082555a1e7bed4c                                 md5
args.go           src/args/args.go           314  5864498dab678a250fd8dd86902ac26ed2ae763395718313d09272e130cb670a sha256
parse_args.go     src/args/parse_args.go     3184 ec0ca60f7e27a31e74a35f55cc224e10                                 md5
parse_args.go     src/args/parse_args.go     3184 d7d4d340d49e4d9f0096538929cb4c9153f96df23dea4450825df99a34abfa46 sha256
```

**Condensed Format** (`-f condensed`):
```
$ go run dirhash.go -i src/args -t -a md5 sha256 -f condensed
File Name         Path                       Size Hashes
args.go           src/args/args.go           314  md5:4b140b31128a0eabb082555a1e7bed4c | sha256:5864498dab678a250fd8dd86902ac26ed2ae763395718313d09272e130cb670a
parse_args.go     src/args/parse_args.go     3184 md5:ec0ca60f7e27a31e74a35f55cc224e10 | sha256:d7d4d340d49e4d9f0096538929cb4c9153f96df23dea4450825df99a34abfa46
```

**IOC Format** (`-f ioc`):
```
$ go run dirhash.go -i src/args -t -a md5 sha1 sha256 -f ioc
File Path              File Name     Size MD5                              SHA1                                     SHA256
src/args/args.go       args.go       314  4b140b31128a0eabb082555a1e7bed4c 2be014706c9b91ca8e30352c40f906dcbbdbbcec 5864498dab678a250fd8dd86902ac26ed2ae763395718313d09272e130cb670a
src/args/parse_args.go parse_args.go 3184 ec0ca60f7e27a31e74a35f55cc224e10 ea57b1ee0e60e87f826d40f4ce9706d7bad0f938 d7d4d340d49e4d9f0096538929cb4c9153f96df23dea4450825df99a34abfa46
```

#### Standard Format CSV Output
```
Path,FileName,Size,Hash,HashType
src/files/hash_files.go,hash_files.go,3762,507e44f05e69d0057101e7d0b14cb9d8,md5
src/files/hash_files.go,hash_files.go,3762,70ed2662d7e772e9cda622d053e3982c95d341646f3865403a1a470b2be44bb2,sha256
```

#### Condensed Format CSV Output (`-f condensed`)
```
Path,FileName,Size,md5,sha256
src/files/file.go,file.go,1167,f28df1fb9d3a80512f34b9ce2031fffd,2366f1379994276a372dfca9481bdee6376fd23051522fee9b7555a1cc1f4628
src/files/hash_files.go,hash_files.go,3762,507e44f05e69d0057101e7d0b14cb9d8,70ed2662d7e772e9cda622d053e3982c95d341646f3865403a1a470b2be44bb2
```

#### IOC Format CSV Output (`-f ioc`)
```
file_path,file_name,file_size,md5,sha1,sha256,sha512
src/files/file.go,file.go,1167,f28df1fb9d3a80512f34b9ce2031fffd,d19a6dcdce7ec9c707b8ce6ec0003c8293144678,2366f1379994276a372dfca9481bdee6376fd23051522fee9b7555a1cc1f4628,a7bbfc6ded135cd67db441a0900d160c2be0aeeffcb5a347c002ee3fe140eddaf2868cdc9176232961f1228095824ada880a41c92508da70fd8e8e4253da0b19
src/files/hash_files.go,hash_files.go,3762,507e44f05e69d0057101e7d0b14cb9d8,359abf9272d8636f97ae66e1fcda4c05420dc406,70ed2662d7e772e9cda622d053e3982c95d341646f3865403a1a470b2be44bb2,588b500917b62c5bb045de4df50b5110d8f8b62da7ac44994fe9a6042877889117231f86b169e0d00cea4f30b63b2dd6db0bc776402728915caf1798d2d28b0c
```

## Security & Threat Intelligence Use Cases

### YARA Rule Generation
IOC format provides all necessary data for YARA rule creation:
```bash
# Generate IOC data for suspicious files (file output)
go run dirhash.go -i /suspicious/samples -o malware_hashes.csv -a md5 sha1 sha256 sha512 -f ioc

# Quick terminal review of IOC data
go run dirhash.go -i /suspicious/samples -t -a md5 sha1 sha256 sha512 -f ioc
```

### Microsoft Sentinel Integration
IOC format is optimized for Microsoft Sentinel IOC imports and KQL queries:
```bash
# Create IOC feed for Sentinel (file output)
go run dirhash.go -i /threat/samples -o sentinel_iocs.csv -f ioc -a md5 sha1 sha256

# Preview IOC data in terminal
go run dirhash.go -i /threat/samples -t -f ioc -a md5 sha1 sha256
```

### Threat Hunting
Condensed format provides efficient bulk hash analysis for threat hunting workflows:
```bash
# Analyze large file sets efficiently (file output)
go run dirhash.go -i /system/files -o analysis.csv -f condensed -a sha256 md5

# Quick condensed terminal view for analysis
go run dirhash.go -i /system/files -t -f condensed -a sha256 md5
```

## Roadmap
Future enhancements planned:
1. KQL query template generation

