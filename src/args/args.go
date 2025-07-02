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
	//YARA Options
	YaraOutput   bool
	YaraFile     string
	YaraRuleName string
	YaraHashOnly bool
	//Flags
	Help bool
}
