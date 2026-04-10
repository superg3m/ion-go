package Codegen

import (
	"os"
	"strconv"
)

type AssemblyEmitter interface {
	EmitInstructions(filepath string)
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

type ATTSystemVAssemblyEmitter struct {
	StringConstantCount int
	LabelCount          int

	instructions      []string
	registerAllocator *ATTSystemVRegisterAllocator
}

func (e *ATTSystemVAssemblyEmitter) EmitInstructions(filepath string) {
	f, _ := os.Create(filepath)
	for _, inst := range e.instructions {
		_, err := f.WriteString(inst + "\n")
		if err != nil {
			panic("Could not emit instruction: " + inst)
		}
	}
}

// EmitLoadIntegerConstant this MOV instruction should just be a constant probably
func (e *ATTSystemVAssemblyEmitter) EmitLoadIntegerConstant(integerConstant int) IntegerRegister {
	register := e.registerAllocator.AcquireIntegerRegister()
	registerName := e.registerAllocator.GetInteger32RegisterName(register)
	e.AddInstruction("\tmovl " + "$" + strconv.Itoa(integerConstant) + ", " + registerName)

	return register
}

func (e *ATTSystemVAssemblyEmitter) AddInstruction(instruction string) {
	e.instructions = append(e.instructions, instruction)
}

func NewATTAssemblyEmitter() AssemblyEmitter {
	return &ATTSystemVAssemblyEmitter{
		StringConstantCount: 0,
		LabelCount:          0,
		instructions:        []string{},
		registerAllocator:   NewAMD64SystemVRegisterAllocator().(*ATTSystemVRegisterAllocator),
	}
}
