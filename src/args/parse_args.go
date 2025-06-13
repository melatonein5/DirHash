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
	args.StrHashAlgorithm = "md5"
	args.HashAlgorithmId = 0
	args.OutputToTerminal = false
	args.WriteToFile = true
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
				nextArg := i + 1
				if nextArg >= rawArgsLen {
					return args, errors.New("missing value for -a | --algorithm flag")
				}
				args.StrHashAlgorithm = rawArgs[nextArg]
				args.HashAlgorithmId = StrHashAlgorithmToId(args.StrHashAlgorithm)
				if err := HashAlgorithmValidation(args.HashAlgorithmId); err != nil {
					return args, err
				}
				i += 2 // Skip the next argument since it's the value for the flag
			case "-t", "--terminal":
				//Output to terminal
				args.OutputToTerminal = true
				// Move to the next argument
				i++
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

	return args, nil
}
