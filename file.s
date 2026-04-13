.section .rodata
.align 8
.type p1, @object
.global p1
.size p1, 24
p1:
	.zero 24
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
	popq %r15
	popq %r14
	popq %r13
	popq %r12
	popq %rdi
	leave
	ret
