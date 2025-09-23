package TypeChecker

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

type Status int

const (
	NORMAL Status = iota
	IN_LOOP
)

type TypeEnv struct {
	parent        *TypeEnv
	variables     map[string]*AST.DeclarationVariable
	CurrentStatus Status
}

func NewTypeEnv(parent *TypeEnv) *TypeEnv {
	return &TypeEnv{
		parent:        parent,
		variables:     make(map[string]*AST.DeclarationVariable),
		CurrentStatus: NORMAL,
	}
}

func (t *TypeEnv) has(key Token.Token) bool {
	current := t
	for current != nil {
		_, ok := current.variables[key.Lexeme]
		if ok {
			return true
		}
		current = current.parent
	}

	return false
}

func (t *TypeEnv) get(key Token.Token) *AST.DeclarationVariable {
	current := t
	for current != nil {
		value, ok := current.variables[key.Lexeme]
		if ok {
			return value
		}
		current = current.parent
	}

	panic(fmt.Sprintf("Line %d | Undeclared Identifier: %s", key.Line, key.Lexeme))
	return &AST.DeclarationVariable{}
}

func (t *TypeEnv) set(key Token.Token, value *AST.DeclarationVariable) {
	if t.has(key) {
		panic(fmt.Sprintf("Line: %d | Variable %s already defined", key.Line, key.Lexeme))
	}

	t.variables[key.Lexeme] = value
}
