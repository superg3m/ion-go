package AST

type Declaration interface {
	Node
	isDeclaration()
}

type DeclarationVariable struct {
	Name     string
	DeclType DataType
	RHS      Expression
}

type Parameter struct {
	Name string
}

type DeclarationFunction struct {
	Name       string
	Parameters []Parameter
	ReturnType DataType
	Block      *StatementBlock
}
