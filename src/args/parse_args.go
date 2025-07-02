package args

import (
	"errors"
)

// ParseArgs will parse the arguments provided in the command line
func ParseArgs(rawArgs []string) (Args, error) {
	//Return Value
	var args Args

	//Default values
	args.StrInputDir = "."
	args.StrOutputFile = ""
	args.StrHashAlgorithms = []string{}
	args.HashAlgorithmId = []int{}
	args.OutputToTerminal = false
	args.WriteToFile = true
	args.OutputFormat = "standard"
	args.Help = false

	//Get the length of the raw arguments for later use
	rawArgsLen := len(rawArgs)

	//Parse the arguments
	for i := 0; i < rawArgsLen; {
		//Check if the argument is a flag
		if rawArgs[i][0] == '-' {
			switch rawArgs[i] {
			case "-i", "--input-dir":
				nextArg := i + 1
				//Input directory
				if nextArg >= rawArgsLen {
					return args, errors.New("missing value for -i | --input-dir flag")
				}
				args.StrInputDir = rawArgs[nextArg]
				// Skip the next argument since it's the value for the flag
				i += 2
			case "-o", "--output":
				//Output file
				nextArg := i + 1
				if nextArg >= rawArgsLen {
					return args, errors.New("missing value for -o | --output flag")
				}
				args.StrOutputFile = rawArgs[nextArg]
				// Skip the next argument since it's the value for the flag
				i += 2
			case "-a", "--algorithm":
				//Hash algorithm
				//There can be multiple algorithms specified, so we need to loop until we hit a flag or run out of arguments
				for j := i + 1; j < rawArgsLen && rawArgs[j][0] != '-'; j++ {
					args.StrHashAlgorithms = append(args.StrHashAlgorithms, rawArgs[j])
					// Convert the string hash algorithm to an int and append it to the HashAlgorithmId slice
					id := StrHashAlgorithmToId(rawArgs[j])
					if id == -1 {
						return args, errors.New("invalid hash algorithm: " + rawArgs[j])
					}
					args.HashAlgorithmId = append(args.HashAlgorithmId, id)
					// Move to the next argument
					i = j
				}
				// Move to the next argument
				i++
			case "-t", "--terminal":
				//Output to terminal
				args.OutputToTerminal = true
				// Move to the next argument
				i++
			case "-f", "--format":
				//Output format
				nextArg := i + 1
				if nextArg >= rawArgsLen {
					return args, errors.New("missing value for -f | --format flag")
				}
				format := rawArgs[nextArg]
				if format != "standard" && format != "condensed" && format != "ioc" {
					return args, errors.New("invalid output format: " + format + ". Valid options: standard, condensed, ioc")
				}
				args.OutputFormat = format
				// Skip the next argument since it's the value for the flag
				i += 2
			case "-h", "--help":
				//Help flag
				args.Help = true
				// Move to the next argument
				i++
			default:
				return args, errors.New("unknown flag: " + rawArgs[i])
			}
		} else {
			return args, errors.New("unexpected argument: " + rawArgs[i])
		}
	}

	// If there is no output file specified, output to terminal by default
	if args.StrOutputFile == "" {
		args.OutputToTerminal = true
		args.WriteToFile = false
	}

	//If no hash algorithms were specified, default to md5
	if len(args.StrHashAlgorithms) == 0 {
		args.StrHashAlgorithms = []string{"md5"}
		args.HashAlgorithmId = []int{0} // MD5
	}

	return args, nil
}
