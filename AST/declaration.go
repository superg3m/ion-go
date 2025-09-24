package AST

import (
	"ion-go/TS"
	"ion-go/Token"
)

type Declaration interface {
	Node
	isDeclaration()
}

type DeclarationVariable struct {
	Tok      Token.Token
	DeclType *TS.Type
	RHS      Expression
}

type DeclarationFunction struct {
	Tok      Token.Token
	DeclType *TS.Type
	Block    *StatementBlock
}
