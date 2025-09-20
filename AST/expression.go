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
}

type ExpressionLen struct {
	Array Expression
}

type ExpressionArrayAccess struct {
	Name  string
	Index Expression
}

// Identifier
// IntegerExpr
// FloatExpr
// StringExpr
// BooleanExpr

// Unary
// Binary
