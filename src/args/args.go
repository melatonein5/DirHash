package args

type Args struct {
	//InputDir
	StrInputDir string
	//OutputFile
	StrOutputFile    string
	OutputToTerminal bool
	WriteToFile      bool
	//HashAlgorithm
	StrHashAlgorithm string
	HashAlgorithmId  int
	//Flags
	Help bool
}
