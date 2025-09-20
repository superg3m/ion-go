package AST

type Statement interface {
	Node
	isStatement()
}

type StatementAssignment struct {
	Name string
	RHS  Expression
}

type StatementPrint struct {
	Expr Expression
}

type StatementBlock struct {
	Body []Node
}

type StatementReturn struct {
	Expr Expression
}

type StatementFor struct {
	Initializer *DeclarationVariable
	Condition   Expression
	Increment   *StatementAssignment
	Body        *StatementBlock
}

// Block
// Assignment
// VariableDeclaration
