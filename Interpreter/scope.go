package Interpreter

import "ion-go/AST"

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

func (s *Scope) has(key string) bool {
	current := s
	for current != nil {
		_, ok := current.variables[key]
		if ok {
			return true
		}
		current = current.parent
	}

	return false
}

func (s *Scope) get(key string) AST.Expression {
	current := s
	for current != nil {
		value, ok := current.variables[key]
		if ok {
			return value
		}
		current = current.parent
	}

	panic("Undeclared Identifier: " + key)
	return nil
}

func (s *Scope) set(key string, value AST.Expression) {
	current := s
	for current != nil {
		_, ok := current.variables[key]
		if ok {
			current.variables[key] = value
			return
		}
		current = current.parent
	}

	s.variables[key] = value
}
