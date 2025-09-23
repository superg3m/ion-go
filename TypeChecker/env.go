package TypeChecker

import "ion-go/AST"

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

func (t *TypeEnv) has(key string) bool {
	current := t
	for current != nil {
		_, ok := current.variables[key]
		if ok {
			return true
		}
		current = current.parent
	}

	return false
}

func (t *TypeEnv) get(key string) *AST.DeclarationVariable {
	current := t
	for current != nil {
		value, ok := current.variables[key]
		if ok {
			return value
		}
		current = current.parent
	}

	panic("Undeclared Identifier: " + key)
	return &AST.DeclarationVariable{}
}

func (t *TypeEnv) set(key string, value *AST.DeclarationVariable) {
	if t.has(key) {
		panic("Variable " + key + " already defined")
	}

	t.variables[key] = value
}
