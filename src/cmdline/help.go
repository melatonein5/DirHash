package cmdline

import "fmt"

func PrintHelp() {
	helpText := `
Usage: dirhash [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -f, --format <format>    Specify the output format for both terminal and file output (default: standard)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Supported algorithms:
  md5, sha1, sha256, sha512
Supported output formats:
  standard  - Traditional format with separate rows per hash type
  condensed - All hashes on single row per file  
  ioc       - IOC-friendly format for security tools (YARA, KQL, Sentinel)
Examples:
  dirhash -i /path/to/dir -o output.csv -a sha256
  dirhash --input-dir /path/to/dir --output output.csv --algorithm sha512 sha1 --format condensed
  dirhash -i /suspicious/files -o iocs.csv -a md5 sha1 sha256 sha512 -f ioc
  dirhash -t
  dirhash --help
`
	fmt.Print(helpText)
}
