# DirHash

[![Test Coverage](https://github.com/melatonein5/DirHash/actions/workflows/test-coverage.yml/badge.svg)](https://github.com/melatonein5/DirHash/actions/workflows/test-coverage.yml)
![Coverage](https://img.shields.io/badge/Coverage-91.7%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/melatonein5/DirHash)](https://goreportcard.com/report/github.com/melatonein5/DirHash)

DirHash is a command line tool that generates cryptographic file hashes and exports detection rules for security analysis. It supports multiple hash algorithms, flexible output formats, and can generate both YARA rules and KQL queries for threat hunting and malware detection.

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

File Processing Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format for both terminal and file output (default: standard)
  -t, --terminal           Output to terminal (default: false)

YARA Rule Generation Options:
  -y, --yara <file>        Generate YARA rule and save to specified file
  --yara-rule-name <name>  Specify custom name for generated YARA rule
  --yara-hash-only         Generate hash-only rules without filenames

KQL Query Generation Options:
  -q, --kql <file>         Generate KQL query and save to specified file
  --kql-name <name>        Specify custom name for generated KQL query
  --kql-hash-only          Generate hash-only queries without filenames
  --kql-tables <tables>    Specify target tables (default: DeviceFileEvents), can take more than 1 argument

General Options:
  -h, --help               Show this help message and exit

Supported algorithms:
  md5, sha1, sha256, sha512

Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)

Supported KQL tables:
  DeviceFileEvents    - Microsoft 365 Defender file events (default)
  SecurityEvents      - Windows security events
  CommonSecurityLog   - Common security log format

Examples:
  Basic file hashing:
    ./dirhash -i /path/to/dir -o output.csv -a sha256
    ./dirhash --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
    ./dirhash -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc

  YARA rule generation:
    ./dirhash -i /malware/samples -y detection.yar --yara-rule-name malware_detection
    ./dirhash -i /files -a sha256 sha512 -y hashes.yar --yara-hash-only
    ./dirhash -i /suspicious -o results.csv -y rules.yar --yara-rule-name threat_hunt

  KQL query generation:
    ./dirhash -i /malware/samples -q detection.kql --kql-name malware_hunt
    ./dirhash -i /files -a sha256 sha512 -q hashes.kql --kql-hash-only
    ./dirhash -i /suspicious -q security.kql --kql-tables DeviceFileEvents SecurityEvents
    ./dirhash -i /threats -o iocs.csv -q hunt.kql --kql-name threat_detection

  Combined YARA and KQL generation:
    ./dirhash -i /malware -o results.csv -y rules.yar -q queries.kql -a sha256 sha512
    ./dirhash -i /samples -y detection.yar -q hunting.kql --yara-rule-name malware --kql-name threats

  Terminal output:
    ./dirhash -t
    ./dirhash -i /files -t -a sha256

  Help:
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

File Processing Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format for both terminal and file output (default: standard)
  -t, --terminal           Output to terminal (default: false)

YARA Rule Generation Options:
  -y, --yara <file>        Generate YARA rule and save to specified file
  --yara-rule-name <name>  Specify custom name for generated YARA rule
  --yara-hash-only         Generate hash-only rules without filenames

KQL Query Generation Options:
  -q, --kql <file>         Generate KQL query and save to specified file
  --kql-name <name>        Specify custom name for generated KQL query
  --kql-hash-only          Generate hash-only queries without filenames
  --kql-tables <tables>    Specify target tables (default: DeviceFileEvents), can take more than 1 argument

General Options:
  -h, --help               Show this help message and exit

Supported algorithms:
  md5, sha1, sha256, sha512

Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)

Supported KQL tables:
  DeviceFileEvents    - Microsoft 365 Defender file events (default)
  SecurityEvents      - Windows security events
  CommonSecurityLog   - Common security log format

Examples:
  Basic file hashing:
    go run dirhash.go -i /path/to/dir -o output.csv -a sha256
    go run dirhash.go --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
    go run dirhash.go -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc

  YARA rule generation:
    go run dirhash.go -i /malware/samples -y detection.yar --yara-rule-name malware_detection
    go run dirhash.go -i /files -a sha256 sha512 -y hashes.yar --yara-hash-only
    go run dirhash.go -i /suspicious -o results.csv -y rules.yar --yara-rule-name threat_hunt

  KQL query generation:
    go run dirhash.go -i /malware/samples -q detection.kql --kql-name malware_hunt
    go run dirhash.go -i /files -a sha256 sha512 -q hashes.kql --kql-hash-only
    go run dirhash.go -i /suspicious -q security.kql --kql-tables DeviceFileEvents SecurityEvents
    go run dirhash.go -i /threats -o iocs.csv -q hunt.kql --kql-name threat_detection

  Combined YARA and KQL generation:
    go run dirhash.go -i /malware -o results.csv -y rules.yar -q queries.kql -a sha256 sha512
    go run dirhash.go -i /samples -y detection.yar -q hunting.kql --yara-rule-name malware --kql-name threats

  Terminal output:
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
DirHash can automatically generate YARA rules for malware detection:
```bash
# Generate YARA rule from suspicious files
go run dirhash.go -i /malware/samples -y malware_detection.yar --yara-rule-name APT_Malware_Family

# Generate hash-only YARA rules (no filenames)
go run dirhash.go -i /suspicious/files -y hashes.yar --yara-hash-only -a sha256 sha512

# Combined hash analysis and YARA rule generation
go run dirhash.go -i /threat/samples -o analysis.csv -y detection.yar -f ioc -a md5 sha1 sha256
```

### KQL Query Generation for Microsoft Security
Generate KQL queries for threat hunting in Microsoft 365 Defender, Azure Sentinel, and Log Analytics:
```bash
# Generate KQL query for Microsoft 365 Defender
go run dirhash.go -i /malware/samples -q threat_hunt.kql --kql-name Advanced_Threat_Detection

# Multi-table KQL queries for comprehensive hunting
go run dirhash.go -i /suspicious/files -q security_hunt.kql --kql-tables DeviceFileEvents SecurityEvents

# Hash-only KQL queries for broader detection
go run dirhash.go -i /iocs -q hash_detection.kql --kql-hash-only -a sha256 sha512

# Combined analysis with both YARA and KQL generation
go run dirhash.go -i /threat/samples -o iocs.csv -y rules.yar -q queries.kql --kql-name ThreatHunt -f ioc
```

### Microsoft Sentinel Integration
IOC format is optimized for Microsoft Sentinel IOC imports and works seamlessly with generated KQL queries:
```bash
# Create complete Sentinel threat hunting package
go run dirhash.go -i /threat/samples -o sentinel_iocs.csv -q sentinel_hunt.kql -f ioc -a md5 sha1 sha256

# Generate IOC feed and KQL queries for Sentinel
go run dirhash.go -i /malware/collection -o iocs.csv -q hunting.kql --kql-tables DeviceFileEvents CommonSecurityLog
```

### Threat Hunting Workflows
DirHash supports complete threat hunting workflows from hash generation to detection rule creation:
```bash
# Complete threat hunting workflow
go run dirhash.go -i /suspicious/files -o analysis.csv -y detection.yar -q hunt.kql -f ioc -a sha256 md5

# Bulk analysis with condensed output for large datasets
go run dirhash.go -i /system/files -o bulk_analysis.csv -f condensed -a sha256

# Multi-platform detection rule generation
go run dirhash.go -i /malware/samples -y yara_rules.yar -q kql_queries.kql --yara-rule-name MalwareFamily --kql-name ThreatDetection
```

### Detection Rule Examples

**Generated YARA Rule Example:**
```yara
rule malware_detection {
    meta:
        author = "DirHash"
        description = "Auto-generated rule for malware detection"
        date = "2025-07-03"
        
    strings:
        $hash_md5_0 = "d6eb32081c822ed572b70567826d9d9d"
        $hash_sha256_0 = "a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae"
        $filename_0 = "malware.exe"
        
    condition:
        any of ($hash_*) or any of ($filename_*)
}
```

**Generated KQL Query Example:**
```kql
// KQL Query: threat_detection
// Description: KQL query to detect files based on hashes and filenames
// Author: DirHash
// Generated: 2025-07-03 19:30:48 UTC

DeviceFileEvents
| where TimeGenerated >= ago(7d)
| where ((MD5 in ("d6eb32081c822ed572b70567826d9d9d")) or (SHA256 in ("a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae"))) or (FileName in~ ("malware.exe"))
| project TimeGenerated, DeviceName, FileName, FolderPath, MD5, SHA1, SHA256, ProcessCommandLine, InitiatingProcessFileName
| extend SourceTable = "DeviceFileEvents"
| sort by TimeGenerated desc
| take 1000
```

## KQL Query Output Examples

DirHash can generate KQL queries for Microsoft security platforms (Microsoft 365 Defender, Azure Sentinel, Log Analytics):

### Standard KQL Query Generation
```bash
# Generate KQL query for threat hunting
go run dirhash.go -i /suspicious/files -q threat_hunt.kql --kql-name malware_detection -a md5 sha256
```

**Generated KQL Query Output:**
```kql
// KQL Query: malware_detection
// Description: KQL query to detect files based on hashes and filenames - Generated from 3 files
// Author: DirHash
// Generated: 2025-07-03 19:30:48 UTC
// Tags: threat-hunting, file-detection, security, dirhash
//
// Hash Count: 6
// Hash Types: md5, sha256
// Filename Count: 3
// Tables: DeviceFileEvents
// Time Range: 7d
// Max Results: 1000
//
// This query searches for files based on cryptographic hashes and filenames.
// It can be used for threat hunting, incident response, and security analysis.
// Modify the time range and result limits as needed for your environment.

DeviceFileEvents
| where TimeGenerated >= ago(7d)
| where ((MD5 in ("d6eb32081c822ed572b70567826d9d9d", "abc123def456ghi789", "xyz789abc012def345")) or (SHA256 in ("a1fff0ffefb9eace7230c24e50731f0a91c62f9cefdfe77121c2f607125dffae", "def456abc123ghi789jkl012", "mno345pqr678stu901vwx234"))) or (FileName in~ ("malware.exe", "trojan.dll", "suspicious.bin"))
| project TimeGenerated, DeviceName, FileName, FolderPath, MD5, SHA1, SHA256, ProcessCommandLine, InitiatingProcessFileName
| extend SourceTable = "DeviceFileEvents"
| sort by TimeGenerated desc
| take 1000
```

### Hash-Only KQL Query
```bash
# Generate hash-only KQL query (no filenames)
go run dirhash.go -i /iocs -q hash_detection.kql --kql-hash-only -a sha256 sha512
```

**Generated Hash-Only KQL Output:**
```kql
// KQL Query: dirhash_generated_query
// Description: KQL query to detect files based on hashes and filenames - Generated from 2 files
// Author: DirHash
// Generated: 2025-07-03 19:35:12 UTC

DeviceFileEvents
| where TimeGenerated >= ago(7d)
| where ((SHA256 in ("def456abc123ghi789jkl012mno345", "pqr678stu901vwx234yzabc012def456")) or (SHA512 in ("ghi789jkl012mno345pqr678stu901vwx234yzabc012def456", "789abc012def456ghi789jkl012mno345pqr678stu901vwx234")))
| project TimeGenerated, DeviceName, FileName, FolderPath, MD5, SHA1, SHA256, ProcessCommandLine, InitiatingProcessFileName
| extend SourceTable = "DeviceFileEvents"
| sort by TimeGenerated desc
| take 1000
```

### Multi-Table KQL Query
```bash
# Generate KQL query for multiple security log sources
go run dirhash.go -i /threats -q comprehensive_hunt.kql --kql-tables DeviceFileEvents SecurityEvents CommonSecurityLog
```

**Generated Multi-Table KQL Output:**
```kql
// KQL Query: dirhash_generated_query
// Description: KQL query to detect files based on hashes and filenames - Generated from 2 files
// Author: DirHash
// Generated: 2025-07-03 19:40:25 UTC

union (
DeviceFileEvents
| where TimeGenerated >= ago(7d)
| where ((MD5 in ("abc123def456")) or (SHA256 in ("def456abc123ghi789"))) or (FileName in~ ("threat.exe"))
| project TimeGenerated, DeviceName, FileName, FolderPath, MD5, SHA1, SHA256, ProcessCommandLine, InitiatingProcessFileName
| extend SourceTable = "DeviceFileEvents"
),
(
SecurityEvents
| where TimeGenerated >= ago(7d)
| where ((FileHash in ("abc123def456", "def456abc123ghi789"))) or (FileName in~ ("threat.exe"))
| project TimeGenerated, Computer, FileName, FilePath, FileHash, ProcessName, CommandLine
| extend SourceTable = "SecurityEvents"
),
(
CommonSecurityLog
| where TimeGenerated >= ago(7d)
| where ((FileHash in ("abc123def456", "def456abc123ghi789"))) or (FileName in~ ("threat.exe"))
| project TimeGenerated, Computer, FileName, FilePath, FileHash, ProcessName, CommandLine
| extend SourceTable = "CommonSecurityLog"
)
| sort by TimeGenerated desc
| take 1000
```

