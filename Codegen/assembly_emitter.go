package Codegen

type AssemblyEmitter interface {
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
	// EmitLoadIntegerConstant()

	// for like strings .ascii
	// CreateDataDeclaration

	// EmitAssignment()

	// EmitFunctionPrologue()
	// EmitFunctionEpilogue()

	// GetNextStringLabel() string
	// GetNextGenericLabel() string
}
