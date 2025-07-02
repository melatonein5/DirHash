# DirHash
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
DirHash supports three output formats to suit different use cases:

### Standard Format (default)
Traditional format with separate rows for each hash type per file. Compatible with existing workflows.

### Condensed Format
All hash values for each file on a single row with dynamic columns based on requested algorithms. More compact for analysis.

### IOC Format
Security-optimized format with standardized columns perfect for threat intelligence tools, YARA rules, KQL queries, and Microsoft Sentinel IOC imports.

## Example Outputs
Terminal output shows results in formatted columns. File output saves to CSV format with the selected format option.

##### Terminal Output
```
[jmn@two DirHash]$ go run dirhash.go -a sha1 sha256
File Name                              Path                                                   Hash                                                             Hash Type
COMMIT_EDITMSG                         .git/COMMIT_EDITMSG                                    ac1f40bb9f2231060330d9c7a4ec9ab2355e300f                         sha1
FETCH_HEAD                             .git/FETCH_HEAD                                        5d037ee061137d7c6088c8a8b1f08b88d4899843                         sha1
HEAD                                   .git/HEAD                                              2aa05cb189709905d22504077e79b9d7ed74722a                         sha1
ORIG_HEAD                              .git/ORIG_HEAD                                         039435964f90dfd504265f1feb58e0eaff896bd2                         sha1
config                                 .git/config                                            5312c39011db3903bc8595f37f9ede2b761ec7b0                         sha1
description                            .git/description                                       9eca422d6263200fdf78796cab2765eb3fdc37e5                         sha1
applypatch-msg.sample                  .git/hooks/applypatch-msg.sample                       285716cd0b1f4e75e00d9854480fce108d1b6654                         sha1
commit-msg.sample                      .git/hooks/commit-msg.sample                           aa71dc4a856cae39c4da91d26e4cb1c7f0a8a92f                         sha1
fsmonitor-watchman.sample              .git/hooks/fsmonitor-watchman.sample                   c55bcc85dbbd3785080ff7d236a0102c47dfc5ba                         sha1
post-update.sample                     .git/hooks/post-update.sample                          f2ea063306cc8441634367a514d8fef76da8f45e                         sha1
main                                   .git/refs/heads/main                                   5d342e4e576ac445f7afc2e6c49a3e1986557019a5b1cb4acab447e5f3dbb73d sha256
HEAD                                   .git/refs/remotes/origin/HEAD                          6832d31adc6cd2a5714637a4395ad820983e5482535931d603d0261cae88b837 sha256
main                                   .git/refs/remotes/origin/main                          1378d3f729fbd8f5ec3d570cc94343c069ce13d0cbfb4d237e50f3f550980a68 sha256
LICENSE                                LICENSE                                                51b84b197025b72457944418eac94b71233e080bd0c45245198c7acad0e26aa9 sha256
README.md                              README.md                                              92f61da76c367c58263ace4d954095f6c86239d5ed0efb633ba203f082a689ae sha256
dirhash                                dirhash                                                eb567f761ed92ddb6a990498352bae2e8e878b99e638cc756e3a5ef5a7670a38 sha256
dirhash.go                             dirhash.go                                             8f2096ddeff9543ab56325dd63e9d689e61cd4b48ce074ea663c5a48cf58f5cd sha256
go.mod                                 go.mod                                                 40719d6dbfafc34704d3efd7c4bcea644755b48be7b1589e38cf86a4c7ffdf3b sha256
args.go                                src/args/args.go                                       86da89f3d83ed7fe84a62c946c72401df10274d73c5543a0590cfa0569b1c38e sha256
hash_algorithm.go                      src/args/hash_algorithm.go                             64b723c3cc795a25397cd529cc9148e8b53bc5de83d4750a31fedafe7e9a2d05 sha256
parse_args.go                          src/args/parse_args.go                                 92f5a1bd99e8659cdca3689958861018425611d72691471a050b88385da22281 sha256
help.go                                src/cmdline/help.go                                    b17700c422f02e6de8b2f7f6bac4208600efba96c48e444af7e38c440c6ece10 sha256
output_files.go                        src/cmdline/output_files.go                            e0fc4f16df0961806b235c17bb363f33ab985be7b5e37a48075e95d56dc0186a sha256
enumerate_files.go                     src/files/enumerate_files.go                           c121e905cb787652d97423af50e5f68b67fd62b5d0ea93315a0ca79334830e58 sha256
file.go                                src/files/file.go                                      b60cff7dcd0a219b435ef39e5905796885f51ed3bc0febc0e41befab336374aa sha256
hash_files.go                          src/files/hash_files.go                                30f098d18815d96bc3a10b114486cebcab0e8725e98681d22baced6a508ff752 sha256
md5_files.go                           src/files/md5_files.go                                 43f7e1c3b48076af3ffc6ca10572528e697db175cdf24e1533bf73ba0a1241b0 sha256
sha1_files.go                          src/files/sha1_files.go                                8c54814b5b021ae6085ae91f7bf2c923dd8fe9ff590ee13946f7086c14d5f87c sha256
sha256_files.go                        src/files/sha256_files.go                              2e14eea5233669ffa6c32e8d91e5ce1baaef6aa4cac48bda82607264c77bbf07 sha256
sha512_files.go                        src/files/sha512_files.go                              bb427b881944e600d521532200fb5106b240f940b1d899067ff30832ab667abb sha256
write_output.go                        src/files/write_output.go                              976f2684e4b610428bce50ece072c76351f4cf476f237825ebb8f01b02310a26 sha256
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
# Generate IOC data for suspicious files
go run dirhash.go -i /suspicious/samples -o malware_hashes.csv -a md5 sha1 sha256 sha512 -f ioc
```

### Microsoft Sentinel Integration
IOC format is optimized for Microsoft Sentinel IOC imports and KQL queries:
```bash
# Create IOC feed for Sentinel
go run dirhash.go -i /threat/samples -o sentinel_iocs.csv -f ioc -a md5 sha1 sha256
```

### Threat Hunting
Condensed format provides efficient bulk hash analysis for threat hunting workflows:
```bash
# Analyze large file sets efficiently
go run dirhash.go -i /system/files -o analysis.csv -f condensed -a sha256 md5
```

## Roadmap
Future enhancements planned:
1. Direct YARA rule generation output format
2. KQL query template generation

