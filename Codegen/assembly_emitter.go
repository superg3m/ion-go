package Codegen

import (
	"fmt"
	"strconv"
)

type AssemblyEmitter interface {
	EmitInstructions()
	EmitLoadIntegerConstant(integerConstant int) IntegerRegister
	// EmitAssignment()

	// EmitFunctionExit()
	// EmitFunctionExitReturnInteger()
	// EmitFunctionExitReturnFloat()

	// this can be a pointer, a struct?
	// EmitFunctionExitReturnPointer()

	// EmitFunctionCallSite()
	// EmitFunctionBody()
	// EmitAssignment()
	// EmitIfElse()
	// EmitWhileLoop()
	// EmitForLoop()

	// EmitBinaryExpression()
	// EmitUnaryExpression()

	// EmitBitwiseOr()
	// EmitBitwiseAnd()
	// EmitBitwiseNot()

	// EmitComputeVariableAddress()

	// for like strings .ascii
	// CreateDataDeclaration

	// EmitFunctionPrologue()
	// EmitFunctionEpilogue()

	// GetNextStringLabel() string
	// GetNextGenericLabel() string
}

type AMD64AssemblyEmitter struct {
	StringConstantCount int
	LabelCount          int

	instructions      []string
	registerAllocator RegisterAllocator
}

func (e *AMD64AssemblyEmitter) AddZeroArgumentInstruction(instruction, source, destination string) {
	e.instructions = append(e.instructions, fmt.Sprintf("\n%s", instruction))
}

func (e *AMD64AssemblyEmitter) AddTwoArgumentInstruction(instruction, source, destination string) {
	switch _ := e.registerAllocator.(type) {
	case *AMD64SystemVRegisterAllocator:
		e.instructions = append(e.instructions, fmt.Sprintf("\n%s %s, %s", source, destination, instruction))
		//case *ATTSystemVX64RegisterAllocator:
		// e.instructions = append(e.instructions, fmt.Sprintf("\n%s %s, %s", source, destination, instruction))
	}
}

func (e *AMD64AssemblyEmitter) EmitInstructions() {
}

func (e *AMD64AssemblyEmitter) EmitLoadIntegerConstant(integerConstant int) IntegerRegister {
	register := e.registerAllocator.AcquireIntegerRegister()
	registerName := e.registerAllocator.GetInteger32RegisterName(register)
	e.AddTwoArgumentInstruction("movl", strconv.Itoa(integerConstant), registerName)

	return register
}

func NewIntelMicrosoftX64AssemblyEmitter() AssemblyEmitter {
	return &AMD64AssemblyEmitter{
		StringConstantCount: 0,
		LabelCount:          0,
		instructions:        []string{},
		registerAllocator:   NewAMD64SystemVRegisterAllocator(),
	}
}
