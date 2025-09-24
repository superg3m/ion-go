package AST

type Node interface {
	isNode()
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
