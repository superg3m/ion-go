.section .rodata
.type test, @object
.align 4
.global test
.size test, 4
test: .zero 4
.section .text
.type main,@function
.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	pushq %rdi
	pushq %r12
	pushq %r13
	pushq %r14
	pushq %r15
	movl $6, %r11d
	movl %r11d, %eax
	popq %rdi
	popq %r12
	popq %r13
	popq %r14
	popq %r15
	leave
	ret
