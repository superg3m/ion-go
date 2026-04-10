package Codegen

import (
	"os"
)

type AssemblyEmitter interface {
	EmitInstructions(filepath string)
	EmitLoadIntegerConstant(integerConstant int) IntegerRegister
	AddInstruction(inst string)

	GetRegisterAllocator() RegisterAllocator
	GetCallingConvention() CallingConvention
	GetSyntax() Syntax

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

	instructions []string

	syntax            Syntax
	registerAllocator RegisterAllocator
	callingConvention CallingConvention

	// instructions
}

func (e *AMD64AssemblyEmitter) GetCallingConvention() CallingConvention {
	return e.callingConvention
}

func (e *AMD64AssemblyEmitter) GetSyntax() Syntax {
	return e.syntax
}

func (e *AMD64AssemblyEmitter) AddInstruction(inst string) {
	e.instructions = append(e.instructions, inst)
}

func (e *AMD64AssemblyEmitter) GetRegisterAllocator() RegisterAllocator {
	return e.registerAllocator
}

func (e *AMD64AssemblyEmitter) EmitInstructions(filepath string) {
	f, _ := os.Create(filepath)
	for _, inst := range e.instructions {
		_, err := f.WriteString(inst + "\n")
		if err != nil {
			panic("Could not emit instruction: " + inst)
		}
	}
}

// EmitLoadIntegerConstant this MOV instruction should just be a constant probably
func (e *AMD64AssemblyEmitter) EmitLoadIntegerConstant(integerConstant int) IntegerRegister {
	register := e.registerAllocator.AcquireIntegerRegister()
	e.AddInstruction(e.syntax.IMOVL(register, e.syntax.IntegerConstant(integerConstant)))

	return register
}

func NewAMD64AssemblyEmitter(syntaxType SyntaxType, cc CallingConventionType) AssemblyEmitter {
	var registerAllocator RegisterAllocator
	var callingConvention CallingConvention
	var syntax Syntax
	if cc == SYSYEM_V {
		callingConvention = NewSystemVCallingConvention()
	} else if cc == MICROSOFT_X64 {
		callingConvention = NewMicrosoftX64CallingConvention()
	}

	if syntaxType == ATT {
		registerAllocator = NewATTRegisterAllocator(callingConvention)
		syntax = NewATTSyntax(registerAllocator)
	} else if syntaxType == INTEL {
		registerAllocator = NewIntelRegisterAllocator(callingConvention)
		syntax = NewIntelSyntax(registerAllocator)
	}

	return &AMD64AssemblyEmitter{
		StringConstantCount: 0,
		LabelCount:          0,
		instructions:        []string{},
		registerAllocator:   registerAllocator,
		callingConvention:   callingConvention,
		syntax:              syntax,
	}
}
