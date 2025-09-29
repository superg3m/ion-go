package AST

import (
	"ion-go/TS"
	"ion-go/Token"
)

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
	Tok Token.Token
}

type ExpressionGrouping struct {
	Expr Expression
}

type ExpressionTypeCast struct {
	Tok      Token.Token
	CastType *TS.Type
	Expr     Expression
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

type ExpressionArray struct {
	Elements []Expression
	DeclType *TS.Type
}

type ExpressionArrayAccess struct {
	Tok     Token.Token
	Indices []Expression
}

type ExpressionStruct struct {
	Tok          Token.Token
	MemberValues map[string]Expression
}

type ExpressionStructMemberAccess struct {
	Tok      Token.Token
	Accesses []Token.Token
}

type ExpressionLen struct {
	Iterable Expression
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
