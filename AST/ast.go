package AST

type Node interface {
	isNode()
}

type Deferrable interface {
	Node
	isDeferrable()
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
