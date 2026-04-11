package Codegen

type AssemblyDirective interface {
	ReadOnlyData() string
	Data() string
	BSS() string
	Text() string
	Global() string
}
