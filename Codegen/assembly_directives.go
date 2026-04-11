package Codegen

type AssemblerType int

const (
	GAS AssemblerType = iota
	NASM
	MASM
)

type Directive interface {
	ReadOnlyData() string
	Data() string
	BSS() string
	Text() string
	Global() string
	Type(functionName string) string
}

type GASDirective struct{}

func (g *GASDirective) ReadOnlyData() string {
	return ".rodata"
}

func (g *GASDirective) Data() string {
	return ".data"
}

func (g *GASDirective) BSS() string {
	return ".bss"
}

func (g *GASDirective) Text() string {
	return ".text"
}

func (g *GASDirective) Global() string {
	return ".global"
}

func (g *GASDirective) Type(functionName string) string {
	return ".type " + functionName + ",@function"
}

func NewGASDirective() Directive {
	return &GASDirective{}
}
