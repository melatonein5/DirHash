package args

type Args struct {
	//InputDir
	StrInputDir string
	//OutputFile
	StrOutputFile    string
	OutputToTerminal bool
	WriteToFile      bool
	//HashAlgorithm
	StrHashAlgorithms []string
	HashAlgorithmId   []int
	//Output Format
	OutputFormat string // "standard", "condensed", "ioc"
	//Flags
	Help bool
}
