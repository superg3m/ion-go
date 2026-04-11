package Codegen

import (
	"fmt"
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
	GlobalObject(objectName string, alignment, size int) string
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

func (g *GASDirective) GlobalObject(objectName string, alignment, size int) string {
	// .align 32
	// .type   test, @object
	// .size   test, 88
	ret := fmt.Sprintf(".type %s, @object\n", objectName) +
		fmt.Sprintf(".align %d\n", alignment) +
		fmt.Sprintf(".global %s\n", objectName) +
		fmt.Sprintf(".size %s, %d\n", objectName, size) +
		fmt.Sprintf("%s: .zero %d", objectName, size)

	return ret
}

func NewGASDirective() Directive {
	return &GASDirective{}
}
