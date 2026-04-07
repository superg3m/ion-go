package Codegen

// register allocator

type RegisterAllocator interface {
	// Allocate

	AcquireInteger32Register() int
	AcquireInteger64Register() int
	ReleaseInteger32Register() int
	ReleaseInteger64Register() int

	AcquireFloat32Register() int
	AcquireFloat64Register() int
	ReleaseFloat32Register() int
	ReleaseFloat64Register() int
}

type X64RegisterAllocator struct {
	/*
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
	*/

	CallerRegisters []int
	CalleeRegisters []int

	/*
		8-byte register | Bytes 0-3 | Bytes 0-1 | Byte 0
		%rax                %eax       %ax         %al
		%rcx                %ecx       %cx         %cl
		%rdx                %edx       %dx         %dl
		%rbx                %ebx       %bx         %bl
		%rsi                %esi       %si         %sil
		%rdi                %edi       %di         %dil
		%rsp                %esp       %sp         %spl
		%rbp                %ebp       %bp         %bpl
		%r8                 %r8d       %r8w        %r8b
		%r9                 %r9d       %r9w        %r9b
		%r10                %r10d      %r10w       %r10b
		%r11                %r11d      %r11w       %r11b
		%r12                %r12d      %r12w       %r12b
		%r13                %r13d      %r13w       %r13b
		%r14                %r14d      %r14w       %r14b
		%r15                %r15d      %r15w       %r15b
	*/

	AllocatedIntegerRegister []bool
	Integer32RegisterList    []string
	Integer64RegisterList    []string

	AllocatedFloatRegister []bool
	FloatRegisterList      []string

	// caller registers
	// callee registers
}

func (x X64RegisterAllocator) AcquireInteger32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AcquireInteger64Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) ReleaseInteger32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) ReleaseInteger64Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AcquireFloat32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AcquireFloat64Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) ReleaseFloat32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) ReleaseFloat64Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AllocateInteger32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AllocateInteger64Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AllocateFloat32Register() int {
	//TODO implement me
	panic("implement me")
}

func (x X64RegisterAllocator) AllocateFloat64Register() int {
	//TODO implement me
	panic("implement me")
}

const (
	RCX = 1
	R8  = 2
	R9  = 3
	R10 = 4
	R11 = 5

	RBX = 0
	R12 = 6
	R13 = 7
	R14 = 8
	R15 = 9
)

/*
int callerSavedRegisters[NUM_CALLER_SAVED] = {ECX, R8D, R9D, R10D, R11D};
int calleeSavedRegisters[NUM_CALLEE_SAVED] = {EBX, R12D, R13D, R14D, R15D};
*/

func NewX64RegisterAllocator() RegisterAllocator {
	// gonna have to add the floating point registers to caller registers
	// and callee registers
	// https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170
	// https://wiki.osdev.org/Calling_Conventions#System_V_ABI
	return &X64RegisterAllocator{
		[]int{
			RCX, R8, R9, R10, R11,
		},
		[]int{
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
