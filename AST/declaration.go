package AST

import "ion-go/Token"

type Declaration interface {
	Node
	isDeclaration()
}

type DeclarationVariable struct {
	Tok      Token.Token
	DeclType DataType
	RHS      Expression
}

type Parameter struct {
	Tok      Token.Token
	DeclType DataType
}

type DeclarationFunction struct {
	Tok        Token.Token
	Parameters []Parameter
	ReturnType DataType
	Block      *StatementBlock
}
