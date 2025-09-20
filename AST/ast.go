package AST

type Node interface {
	isNode()
}

type DataType struct {
	Name string
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
