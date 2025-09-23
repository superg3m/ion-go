package AST

import "ion-go/Token"

type Expression interface {
	Node
	isExpression()
}

type ExpressionInteger struct {
	Value int
}

type ExpressionFloat struct {
	Value float32
}

type ExpressionString struct {
	Value string
}

type ExpressionBoolean struct {
	Value bool
}

type ExpressionIdentifier struct {
	Name string
}

type ExpressionGrouping struct {
	Expr Expression
}

type ExpressionUnary struct {
	Operator Token.Token
	Operand  Expression
}

type ExpressionBinary struct {
	Operator Token.Token
	Left     Expression
	Right    Expression
}

type ExpressionFunctionCall struct {
	Name      string
	Arguments []Expression
}

type ExpressionArray struct {
	Elements []Expression
	DeclType DataType
}

type ExpressionLen struct {
	Array Expression
}

type ExpressionArrayAccess struct {
	Name    string
	Indices []Expression
}

type PseudoBehavior int

const (
	BREAK PseudoBehavior = iota
	RETURN
	CONTINUE
)

type ExpressionPseudo struct {
	Expr     Expression
	Behavior PseudoBehavior
}

// Identifier
// IntegerExpr
// FloatExpr
// StringExpr
// BooleanExpr

// Unary
// Binary
