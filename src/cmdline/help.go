package cmdline

import "fmt"

func PrintHelp() {
	helpText := `
Usage: dirhash [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithm (default: md5)
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Examples:
  dirhash -i /path/to/dir -o output.txt -a sha256
  dirhash --input-dir /path/to/dir --output output.txt --algorithm sha512
  dirhash -t
  dirhash --help
`
	fmt.Print(helpText)
}
