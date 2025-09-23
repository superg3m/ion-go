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

type StatementBreak struct{}
type StatementContinue struct{}

type StatementFor struct {
	Initializer *DeclarationVariable
	Condition   Expression
	Increment   *StatementAssignment
	Block       *StatementBlock
}

type StatementIfElse struct {
	Condition Expression
	IfBlock   *StatementBlock
	ElseBlock *StatementBlock
}

// Block
// Assignment
// VariableDeclaration
