package Codegen

import (
	"fmt"
	"ion-go/TS"
)

type AssemblerType int

const (
	GAS AssemblerType = iota
	NASM
	MASM
)

type Directive interface {
	ReadOnlyData() string
	Data() string
	BSS() string
	Text() string
	GlobalFunction(functionName string) string
	GlobalObject(objectName string, t TS.Type) string
}

type GASDirective struct{}

func (g *GASDirective) ReadOnlyData() string {
	return ".section .rodata"
}

func (g *GASDirective) Data() string {
	return ".section .data"
}

func (g *GASDirective) BSS() string {
	return ".section .bss"
}

func (g *GASDirective) Text() string {
	return ".section .text"
}

func (g *GASDirective) GlobalFunction(functionName string) string {
	return ".type " + functionName + ",@function\n" + ".global " + functionName
}

func (g *GASDirective) GlobalObject(objectName string, t TS.Type) string {
	ret := fmt.Sprintf(".align %d\n", t.Align()) +
		fmt.Sprintf(".type %s, @object\n", objectName) +
		fmt.Sprintf(".global %s\n", objectName) +
		fmt.Sprintf(".size %s, %d\n", objectName, t.Size()) +
		fmt.Sprintf("%s:\n\t.zero %d", objectName, t.Size())

	return ret
}

func NewGASDirective() Directive {
	return &GASDirective{}
}
