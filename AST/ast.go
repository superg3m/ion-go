package AST

type Node interface {
	isNode()
}

type TypeModifier string

const (
	NO_MODIFIER TypeModifier = ""
	ARRAY                    = "[]"
)

type DataType struct {
	name string
}

func CreateDataType(name string) DataType {
	return DataType{
		name: name,
	}
}

func (t DataType) String() string {
	return t.name
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
