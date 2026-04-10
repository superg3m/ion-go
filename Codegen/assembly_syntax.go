package Codegen

import (
	"fmt"
	"strconv"
)

type SyntaxType int

const (
	ATT SyntaxType = iota
	INTEL
)

const (
	IntelEffectiveAddressString = "QWORD PTR [%s + %d]" // register + offset (64 bit)
	ATTEffectiveAddressString   = "%s(%s)"              // register + offset
)

type Syntax interface {
	MOVL(destination, source IntegerRegister) string
	IMOVL(destination IntegerRegister, source string) string
	MOVQ(destination, source IntegerRegister) string

	PUSHL(source IntegerRegister) string
	PUSHQ(source IntegerRegister) string

	POPL(destination IntegerRegister) string
	POPQ(destination IntegerRegister) string

	EffectiveAddressString() string
	IntegerConstant(integerConstant int) string
}

type ATTSyntax struct {
	registerAllocator RegisterAllocator
}

func (s *ATTSyntax) IntegerConstant(integerConstant int) string {
	return "$" + strconv.Itoa(integerConstant)
}

func (s *ATTSyntax) EffectiveAddressString() string {
	return "%d(%s)"
}

func (s *ATTSyntax) MOVL(destination, source IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tmovl %s, %s", sourceName, destName)
}

func (s *ATTSyntax) IMOVL(destination IntegerRegister, source string) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	return fmt.Sprintf("\tmovl %s, %s", source, destName)
}

func (s *ATTSyntax) MOVQ(destination, source IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tmovl %s, %s", sourceName, destName)
}

func (s *ATTSyntax) PUSHL(source IntegerRegister) string {
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tpushl %s", sourceName)
}

func (s *ATTSyntax) PUSHQ(source IntegerRegister) string {
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tpushq %s", sourceName)
}

func (s *ATTSyntax) POPL(destination IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)

	return fmt.Sprintf("\tpopl %s", destName)
}

func (s *ATTSyntax) POPQ(destination IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)

	return fmt.Sprintf("\tpopq %s", destName)
}

func NewATTSyntax(registerAllocator RegisterAllocator) Syntax {
	return &ATTSyntax{
		registerAllocator,
	}
}

// --

type IntelSyntax struct {
	registerAllocator RegisterAllocator
}

func (s *IntelSyntax) IMOVL(destination IntegerRegister, source string) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	return fmt.Sprintf("\tmov %s, %s", destName, source)
}

func (s *IntelSyntax) IntegerConstant(integerConstant int) string {
	return strconv.Itoa(integerConstant)
}

func (s *IntelSyntax) EffectiveAddressString() string {
	return "QWORD PTR [%s + %d]"
}

func (s *IntelSyntax) MOVL(destination, source IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tmov %s, %s", destName, sourceName)
}

func (s *IntelSyntax) MOVQ(destination, source IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tmov %s, %s", destName, sourceName)
}

func (s *IntelSyntax) PUSHL(source IntegerRegister) string {
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tpush %s", sourceName)
}

func (s *IntelSyntax) PUSHQ(source IntegerRegister) string {
	sourceName := s.registerAllocator.GetInteger32RegisterName(source)

	return fmt.Sprintf("\tpush %s", sourceName)
}

func (s *IntelSyntax) POPL(destination IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)

	return fmt.Sprintf("\tpop %s", destName)
}

func (s *IntelSyntax) POPQ(destination IntegerRegister) string {
	destName := s.registerAllocator.GetInteger32RegisterName(destination)

	return fmt.Sprintf("\tpop %s", destName)
}

func NewIntelSyntax(registerAllocator RegisterAllocator) Syntax {
	return &IntelSyntax{
		registerAllocator,
	}
}
