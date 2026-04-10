package Codegen

import "ion-go/TS"

type CallingConventionType int

const (
	SYSYEM_V CallingConventionType = iota
	MICROSOFT_X64
)

type CallingConvention interface {
	EmitFunctionPreCall()                  // push caller registers
	EmitFunctionCall(functionType TS.Type) // handle params

	EmitFunctionPrologue(e AssemblyEmitter, functionName string) // handle params, push callee registers, 16 byte alignment
	EmitFunctionEpilogue()                                       // handle params, pop callee registers, ret, leave

	EmitFunctionExit()                 // handle params
	EmitPostCall(functionType TS.Type) // pop caller saved registers

	GetCallerSavedIntegerRegister() []IntegerRegister
	GetCalleeSavedIntegerRegister() []IntegerRegister
}

type CallingConventionSystemV struct{}
type CallingConventionMicrosoftX64 struct{}

func (c *CallingConventionSystemV) EmitFunctionPreCall() {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionSystemV) EmitFunctionCall(functionType TS.Type) {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionSystemV) EmitFunctionPrologue(e AssemblyEmitter, functionName string) {
	s := e.GetSyntax()

	e.AddInstruction("\t.glob " + functionName)
	e.AddInstruction("\t.type " + functionName + ",@function")
	e.AddInstruction("\t" + functionName + ":")
	e.AddInstruction(s.PUSHQ(RBP))
	e.AddInstruction(s.MOVQ(RBP, RSP))
}

func (c *CallingConventionSystemV) EmitFunctionEpilogue() {

	// pop callee

	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionSystemV) EmitFunctionExit() {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionSystemV) EmitPostCall(functionType TS.Type) {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionSystemV) GetCallerSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RCX, R8, R9, R10, R11}
}

func (c *CallingConventionSystemV) GetCalleeSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RBX, R12, R13, R14, R15}
}

func NewSystemVCallingConvention() CallingConvention {
	return &CallingConventionSystemV{}
}

// --

func (c *CallingConventionMicrosoftX64) EmitFunctionPreCall() {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionMicrosoftX64) EmitFunctionCall(functionType TS.Type) {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionMicrosoftX64) EmitFunctionPrologue(e AssemblyEmitter, functionName string) {
	s := e.GetSyntax()

	e.AddInstruction("\t.glob " + functionName)
	e.AddInstruction("\t.type " + functionName + ",@function")
	e.AddInstruction("\t" + functionName + ":")
	e.AddInstruction(s.PUSHQ(RBP))
	e.AddInstruction(s.MOVQ(RBP, RSP))
}

func (c *CallingConventionMicrosoftX64) EmitFunctionEpilogue() {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionMicrosoftX64) EmitFunctionExit() {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionMicrosoftX64) EmitPostCall(functionType TS.Type) {
	//TODO implement me
	panic("implement me")
}

func (c *CallingConventionMicrosoftX64) GetCallerSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{}
}

func (c *CallingConventionMicrosoftX64) GetCalleeSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{}
}

func NewMicrosoftX64CallingConvention() CallingConvention {
	return &CallingConventionMicrosoftX64{}
}
