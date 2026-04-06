package AST

type Node interface {
	isNode()
}

type Deferrable interface {
	Node
	isDeferrable()
}

type Iterable interface {
	Node
	isIterable()
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
