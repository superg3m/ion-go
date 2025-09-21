package AST

import (
	"fmt"
)

type Node interface {
	isNode()
}

type TypeModifier string

const (
	NO_MODIFIER TypeModifier = ""
	ARRAY                    = "[]"
)

type DataType struct {
	name     string
	modifier TypeModifier
}

func CreateDataType(name string, modifier TypeModifier) DataType {
	return DataType{
		name:     name,
		modifier: modifier,
	}
}

func (t DataType) String() string {
	return fmt.Sprintf("%s%s", t.name, t.modifier)
}

type Program struct {
	Declarations []Declaration
}

func (p *Program) isNode() {}
