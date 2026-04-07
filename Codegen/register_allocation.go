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
	Float32RegisterList    []string
	Float64RegisterList    []string
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

func NewX64RegisterAllocator() RegisterAllocator {
	return &X64RegisterAllocator{
		[]bool{false},
		[]string{""},
		[]string{""},
		[]bool{false},
		[]string{""},
		[]string{""},
	}
}
