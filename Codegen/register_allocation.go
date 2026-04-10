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

// ---

/*
int callerSavedRegisters[NUM_CALLER_SAVED] = {ECX, R8D, R9D, R10D, R11D};
int calleeSavedRegisters[NUM_CALLEE_SAVED] = {EBX, R12D, R13D, R14D, R15D};

%rax 	Return value, caller-saved 	%eax 	%ax 	%al
%rdi 	1st argument, caller-saved 	%edi 	%di 	%dil
%rsi 	2nd argument, caller-saved 	%esi 	%si 	%sil
%rdx 	3rd argument, caller-saved 	%edx 	%dx 	%dl
%rcx 	4th argument, caller-saved 	%ecx 	%cx 	%cl
%r8 	5th argument, caller-saved 	%r8d 	%r8w 	%r8b
%r9 	6th argument, caller-saved 	%r9d 	%r9w 	%r9b
%r10 	Scratch/temporary, caller-saved 	%r10d 	%r10w 	%r10b
%r11 	Scratch/temporary, caller-saved

%rsp 	Stack pointer, callee-saved 	%esp 	%sp 	%spl
%rbx 	Local variable, callee-saved 	%ebx 	%bx 	%bl
%rbp 	Local variable, callee-saved 	%ebp 	%bp 	%bpl
%r12 	Local variable, callee-saved 	%r12d 	%r12w 	%r12b
%r13 	Local variable, callee-saved 	%r13d 	%r13w 	%r13b
%r14 	Local variable, callee-saved 	%r14d 	%r14w 	%r14b
%r15 	Local variable, callee-saved 	%r15d 	%r15w 	%r15b

// https://wiki.osdev.org/Calling_Conventions#System_V_ABI
*/

type AMD64SystemVRegisterAllocator struct {
	CallerRegisters []IntegerRegister
	CalleeRegisters []IntegerRegister

	IntegerParameterRegisterMap map[IntegerRegister]IntegerRegisterData
	IntegerRegisterMap          map[IntegerRegister]IntegerRegisterData
}

func (r *AMD64SystemVRegisterAllocator) GetCallerSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RCX, R8, R9, R10, R11}
}

func (r *AMD64SystemVRegisterAllocator) GetCalleeSavedIntegerRegister() []IntegerRegister {
	return []IntegerRegister{RBX, R12, R13, R14, R15}
}

func (r *AMD64SystemVRegisterAllocator) AcquireIntegerRegister() IntegerRegister {
	for register, registerData := range r.IntegerRegisterMap {
		if registerData.Allocated {
			continue
		}

		registerData.Allocated = true
		return register
	}

	panic("Failed to acquire integer register")
}

func (r *AMD64SystemVRegisterAllocator) ReleaseIntegerRegister(register IntegerRegister) {
	if !r.IntegerRegisterMap[register].Allocated {
		panic("Failed to release integer register, not allocated!")
	}

	d := r.IntegerRegisterMap[register]
	d.Allocated = false

	r.IntegerRegisterMap[register] = d
}

func (r *AMD64SystemVRegisterAllocator) IsIntegerRegisterAllocated(register IntegerRegister) bool {
	return r.IntegerRegisterMap[register].Allocated
}

func (r *AMD64SystemVRegisterAllocator) GetInteger32RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer32RegisterName
}

func (r *AMD64SystemVRegisterAllocator) GetInteger64RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer64RegisterName
}

// NewAMD64SystemVRegisterAllocator Syntax: AT&T | https://wiki.osdev.org/Calling_Conventions#System_V_ABI
func NewAMD64SystemVRegisterAllocator() RegisterAllocator {
	return &AMD64SystemVRegisterAllocator{
		[]IntegerRegister{
			RCX, R8, R9, R10, R11,
		},
		[]IntegerRegister{
			RBX, R12, R13, R14, R15,
		},

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

/*
// IntelMicrosoftX64RegisterAllocator https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170
type IntelMicrosoftX64RegisterAllocator struct {
	CallerRegisters []IntegerRegister
	CalleeRegisters []IntegerRegister

	IntegerParameterRegisterMap map[IntegerRegister]IntegerRegisterData
	IntegerRegisterMap          map[IntegerRegister]IntegerRegisterData
}

func (r *IntelMicrosoftX64RegisterAllocator) AcquireIntegerRegister() IntegerRegister {
	for register, registerData := range r.IntegerRegisterMap {
		if registerData.Allocated {
			continue
		}

		registerData.Allocated = true
		return register
	}

	panic("Failed to acquire integer register")
}

func (r *IntelMicrosoftX64RegisterAllocator) ReleaseIntegerRegister(register IntegerRegister) {
	if !r.IntegerRegisterMap[register].Allocated {
		panic("Failed to release integer register, not allocated!")
	}

	d := r.IntegerRegisterMap[register]
	d.Allocated = false

	r.IntegerRegisterMap[register] = d
}

func (r *IntelMicrosoftX64RegisterAllocator) IsIntegerRegisterAllocated(register IntegerRegister) bool {
	return r.IntegerRegisterMap[register].Allocated
}

func (r *IntelMicrosoftX64RegisterAllocator) GetInteger32RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer32RegisterName
}

func (r *IntelMicrosoftX64RegisterAllocator) GetInteger64RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer64RegisterName
}

// NewMicrosoft_X64_RegisterAllocator Syntax: Intel
func NewMicrosoft_X64_RegisterAllocator() RegisterAllocator {
	return &IntelMicrosoftX64RegisterAllocator{
		[]IntegerRegister{
			RAX, RCX, RDX, R8, R10, R11,
		},
		[]IntegerRegister{
			R12, R13, R14, R15, RDI, RSI, RBX, RBP, RSP,
		},

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

*/

// -----------
