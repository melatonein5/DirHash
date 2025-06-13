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
	//Flags
	Help bool
}
