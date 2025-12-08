package models

// Component Variable type
type VariableType int

const (
	Container VariableType = iota
	Number
	Flag
	Text
	Media
)

// Component value
type VariableValue []byte

// Component style
type Style string

// Acts as building block of each sheet
type Component struct {
	Id          ComponentID // Component Id
	UId         UserID // User Id
	SId         SheetID // Parent Sheet Id
	Name        string
	Description string
	Template    bool
	VarName     string
	VarType     VariableType
	VarValue    VariableValue
	Style       Style
}
