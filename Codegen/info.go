package Codegen

// https://web.mit.edu/gnu/doc/html/as_7.html#SEC89
// NASM, MASM, and GAS

// Its Gnu Assembler (GAS)

// Intel Registers
// ATT Registers

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
