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
}

type IntegerRegisterData struct {
	Allocated             bool
	Integer32RegisterName string
	Integer64RegisterName string
}

// IntelMicrosoftRegisterAllocator https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170
type IntelMicrosoftRegisterAllocator struct {
	CallerRegisters []IntegerRegister
	CalleeRegisters []IntegerRegister

	IntegerParameterRegisterMap map[IntegerRegister]IntegerRegisterData
	IntegerRegisterMap          map[IntegerRegister]IntegerRegisterData
}

func (r *IntelMicrosoftRegisterAllocator) AcquireIntegerRegister() IntegerRegister {
	for register, registerData := range r.IntegerRegisterMap {
		if registerData.Allocated {
			continue
		}

		registerData.Allocated = true
		return register
	}

	panic("Failed to acquire integer register")
}

func (r *IntelMicrosoftRegisterAllocator) ReleaseIntegerRegister(register IntegerRegister) {
	if !r.IntegerRegisterMap[register].Allocated {
		panic("Failed to release integer register, not allocated!")
	}

	d := r.IntegerRegisterMap[register]
	d.Allocated = false

	r.IntegerRegisterMap[register] = d
}

func (r *IntelMicrosoftRegisterAllocator) IsIntegerRegisterAllocated(register IntegerRegister) bool {
	return r.IntegerRegisterMap[register].Allocated
}

func (r *IntelMicrosoftRegisterAllocator) GetInteger32RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer32RegisterName
}

func (r *IntelMicrosoftRegisterAllocator) GetInteger64RegisterName(register IntegerRegister) string {
	if !r.IsIntegerRegisterAllocated(register) {
		panic("register not allocated")
	}

	return r.IntegerRegisterMap[register].Integer64RegisterName
}

// NewMicrosoft_X64_RegisterAllocator Syntax: Intel
func NewMicrosoft_X64_RegisterAllocator() RegisterAllocator {
	return &IntelMicrosoftRegisterAllocator{
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

// -----------

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

// Syntax: AT&T
/*
func NewSystemV_X64_RegisterAllocator() RegisterAllocator {
	// gonna have to add the floating point registers to caller registers
	// and callee registers
	// https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170
	// https://wiki.osdev.org/Calling_Conventions#System_V_ABI
	return &X64RegisterAllocator{
		[]IntegerRegister{
			RCX, R8, R9, R10, R11,
		},
		[]IntegerRegister{
			RBX, R12, R13, R14, R15,
		},

		[]bool{
			false, false, false, false,
			false, false, false, false,
			false, false,
		},
		[]string{
			"%ebx", "%ecx", "%r8d", "%r9d",
			"%r10d", "%r11d", "%r12d", "%r13d",
			"%r14d", "%r15d",
		},
		[]string{
			"%rbx", "%rcx", "%r8", "%r9",
			"%r10", "%r11", "%r12", "%r13",
			"%r14", "%r15",
		},

		[]bool{
			false, false, false, false,
			false, false, false, false,
		},
		[]string{
			"%xmm0", "%xmm1", "%xmm2", "%xmm3",
			"%xmm4", "%xmm5", "%xmm6", "%xmm7",
		},
	}
}
*/
