package AST

import "ion-go/Token"

type Statement interface {
	Node
	isStatement()
}

type StatementAssignment struct {
	LHS Expression
	RHS Expression
}

type StatementPrint struct {
	IsNewLine bool
	Expr      Expression
}

type StatementBlock struct {
	Body []Node
}

type StatementReturn struct {
	Expr Expression
}

type StatementDefer struct {
	Tok          Token.Token
	DeferredNode Deferrable
}

type StatementBreak struct{}
type StatementContinue struct{}

type StatementFor struct {
	Initializer *DeclarationVariable
	Condition   Expression
	Increment   *StatementAssignment
	Block       *StatementBlock
}

type StatementWhile struct {
	Condition Expression
	Block     *StatementBlock
}

type StatementIfElse struct {
	Condition Expression
	IfBlock   *StatementBlock
	ElseBlock *StatementBlock
}

// Block
// Assignment
// VariableDeclaration
