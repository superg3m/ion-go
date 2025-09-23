package Interpreter

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

type Scope struct {
	parent    *Scope
	variables map[string]AST.Expression
}

func CreateScope(parent *Scope) Scope {
	return Scope{
		parent:    parent,
		variables: make(map[string]AST.Expression),
	}
}

func (s *Scope) has(key Token.Token) bool {
	current := s
	for current != nil {
		_, ok := current.variables[key.Lexeme]
		if ok {
			return true
		}
		current = current.parent
	}

	return false
}

func (s *Scope) get(key Token.Token) AST.Expression {
	current := s
	for current != nil {
		value, ok := current.variables[key.Lexeme]
		if ok {
			return value
		}
		current = current.parent
	}

	panic(fmt.Sprintf("Line: %d | Undeclared Identifier: %s", key.Line, key.Lexeme))
	return nil
}

func (s *Scope) set(key Token.Token, value AST.Expression) {
	current := s
	for current != nil {
		_, ok := current.variables[key.Lexeme]
		if ok {
			current.variables[key.Lexeme] = value
			return
		}
		current = current.parent
	}

	s.variables[key.Lexeme] = value
}
