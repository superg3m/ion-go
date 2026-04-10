package Codegen

// register allocator

// Volatile is caller saved
// RAX 	Volatile 	Return value register
// RCX 	Volatile 	First integer argument
// RDX 	Volatile 	Second integer argument
// R8 	Volatile 	Third integer argument
// R9 	Volatile 	Fourth integer argument
// R10:R11 	Volatile 	Must be preserved as needed by caller; used in syscall/sysret instructions

// Non-Volatile is callee saved
// R12:R15 	Nonvolatile 	Must be preserved by callee
// RDI 	Nonvolatile 	Must be preserved by callee
// RSI 	Nonvolatile 	Must be preserved by callee
// RBX 	Nonvolatile 	Must be preserved by callee
// RBP 	Nonvolatile 	May be used as a frame pointer; must be preserved by callee
// RSP 	Nonvolatile 	Stack pointer

type IntegerRegister int

const (
	RAX IntegerRegister = iota
	RBP
	RSP
	RDI
	RSI
	RDX
	RBX
	RCX
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
)

// RegisterAllocator Don't allow you to allocate special register
type RegisterAllocator interface {
	AcquireIntegerRegister() IntegerRegister
	ReleaseIntegerRegister(register IntegerRegister)
	IsIntegerRegisterAllocated(register IntegerRegister) bool

	GetInteger32RegisterName(register IntegerRegister) string
	GetInteger64RegisterName(register IntegerRegister) string

	GetCallerSavedIntegerRegister() []IntegerRegister
	GetCalleeSavedIntegerRegister() []IntegerRegister
}

type IntegerRegisterData struct {
	Allocated             bool
	Integer32RegisterName string
	Integer64RegisterName string
}

type ATTRegisterAllocator struct {
	IntegerParameterRegisterMap map[IntegerRegister]IntegerRegisterData
	IntegerRegisterMap          map[IntegerRegister]IntegerRegisterData
}

func (r *ATTRegisterAllocator) GetCallerSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RCX, R8, R9, R10, R11}
}

func (r *ATTRegisterAllocator) GetCalleeSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RBX, R12, R13, R14, R15}
}

func (r *ATTRegisterAllocator) AcquireIntegerRegister() IntegerRegister {
	for register, registerData := range r.IntegerRegisterMap {
		if registerData.Allocated {
			continue
		}

		registerData.Allocated = true
		r.IntegerRegisterMap[register] = registerData

		return register
	}

	panic("Failed to acquire integer register")
}

func (r *ATTRegisterAllocator) ReleaseIntegerRegister(register IntegerRegister) {
	if !r.IntegerRegisterMap[register].Allocated {
		panic("Failed to release integer register, not allocated!")
	}

	d := r.IntegerRegisterMap[register]
	d.Allocated = false

	r.IntegerRegisterMap[register] = d
}

func (r *ATTRegisterAllocator) IsIntegerRegisterAllocated(register IntegerRegister) bool {
	return r.IntegerRegisterMap[register].Allocated
}

func (r *ATTRegisterAllocator) GetInteger32RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer32RegisterName
}

func (r *ATTRegisterAllocator) GetInteger64RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer64RegisterName
}

// NewAMD64SystemVRegisterAllocator Syntax: AT&T | https://wiki.osdev.org/Calling_Conventions#System_V_ABI
func NewATTRegisterAllocator(convention CallingConvention) RegisterAllocator {
	// convention.GetCalleeList
	// convention.GetCallerList

	return &ATTRegisterAllocator{
		// PARAMETERS: rdi, rsi, rdx, rcx, r8, and r9
		map[IntegerRegister]IntegerRegisterData{
			RDI: {false, "%edi", "%rdi"},
			RSI: {false, "%esi", "%rsi"},
			RDX: {false, "%edx", "%rdx"},
			RCX: {false, "%ecx", "%rcx"},
			R8:  {false, "%r8d", "%r8"},
			R9:  {false, "%r9d", "%r9"},
		},

		map[IntegerRegister]IntegerRegisterData{
			RBX: {false, "%edi", "%rdi"},
			R11: {false, "%r11d", "%r11"},
			R12: {false, "%r12d", "%r12"},
			R13: {false, "%r13d", "%r13"},
			R14: {false, "%r14d", "%r14"},
			R15: {false, "%r15d", "%r15"},
		},
	}
}

// ---

// IntelMicrosoftX64RegisterAllocator https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170
type IntelRegisterAllocator struct {
	IntegerParameterRegisterMap map[IntegerRegister]IntegerRegisterData
	IntegerRegisterMap          map[IntegerRegister]IntegerRegisterData
}

func (r *IntelRegisterAllocator) GetCallerSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{
		RAX, RCX, RDX, R8, R10, R11,
	}
}

func (r *IntelRegisterAllocator) GetCalleeSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{
		R12, R13, R14, R15, RDI, RSI, RBX, RBP, RSP,
	}
}

func (r *IntelRegisterAllocator) AcquireIntegerRegister() IntegerRegister {
	for register, registerData := range r.IntegerRegisterMap {
		if registerData.Allocated {
			continue
		}

		registerData.Allocated = true
		r.IntegerRegisterMap[register] = registerData

		return register
	}

	panic("Failed to acquire integer register")
}

func (r *IntelRegisterAllocator) ReleaseIntegerRegister(register IntegerRegister) {
	if !r.IntegerRegisterMap[register].Allocated {
		panic("Failed to release integer register, not allocated!")
	}

	d := r.IntegerRegisterMap[register]
	d.Allocated = false

	r.IntegerRegisterMap[register] = d
}

func (r *IntelRegisterAllocator) IsIntegerRegisterAllocated(register IntegerRegister) bool {
	return r.IntegerRegisterMap[register].Allocated
}

func (r *IntelRegisterAllocator) GetInteger32RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer32RegisterName
}

func (r *IntelRegisterAllocator) GetInteger64RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer64RegisterName
}

// NewMicrosoft_X64_RegisterAllocator Syntax: Intel
func NewIntelRegisterAllocator(cc CallingConvention) RegisterAllocator {
	return &IntelRegisterAllocator{
		// PARAMETERS: 1st: RCX | 2nd: RDX | 3rd: R8 | 4th: R9 | STACK...
		map[IntegerRegister]IntegerRegisterData{
			RCX: {false, "ecx", "rcx"},
			RDX: {false, "esi", "rdx"},
			R8:  {false, "r8d", "r8"},
			R9:  {false, "r9d", "r9"},
		},
		map[IntegerRegister]IntegerRegisterData{
			RDI: {false, "edi", "rdi"},
			RSI: {false, "esi", "rsi"},
			RBX: {false, "ebx", "rbx"},
			R10: {false, "r10d", "r10"},
			R11: {false, "r11d", "r11"},
			R12: {false, "r12d", "r12"},
			R13: {false, "r13d", "r13"},
			R14: {false, "r14d", "r14"},
			R15: {false, "r15d", "r15"},
		},
	}
}
