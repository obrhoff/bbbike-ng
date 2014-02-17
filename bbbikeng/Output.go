package bbbikeng

type BBJSON struct {

	Response bool
	Preferences Preferences
	Distance int
	Time int
	Lights int
	Instruction []BBJSONInstruction
	Path [][2]float64

}

type BBJSONInstruction struct {

	PathIndex int
	Name string
	Type string
	Quality string
	Instruction string

}
