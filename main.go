package main

import (
	"fmt"
	"io/fs"
	"ion-go/Codegen"
	"ion-go/Lexer"
	"ion-go/Parser"
	"ion-go/TypeChecker"
	"os"
	"path"
)

func getIonFileList(dir string) []string {
	root := os.DirFS(dir)

	mdFiles, err := fs.Glob(root, "*.ion")

	if err != nil {
		panic(err)
	}

	var files []string
	for _, v := range mdFiles {
		files = append(files, path.Join(dir, v))
	}

	return files
}

func main() {
	// tokenStream := Lexer.GenerateTokenStream("./factorial.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fib.ion")
	// tokenStream := Lexer.GenerateTokenStream("./array.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fractal.ion")
	// tokenStream := Lexer.GenerateTokenStream("./test.ion")
	// tokenStream := Lexer.GenerateTokenStream("./struct.ion")

	// IonSource
	// iterate through these have like
	// 001_implicit_cast_binary_op.ion
	// 002_implicit_cast_function_arg.ion
	// 003_multi_dimension_static_array_.ion

	testDirectory := "./IonSource"

	for _, file := range getIonFileList(testDirectory) {
		fmt.Printf("\nProcessing: %s\n", file)

		tokenStream := Lexer.GenerateTokenStream(file)
		program := Parser.ParseProgram(tokenStream)
		TypeChecker.TypeCheckProgram(program)
		// JSON.PrettyPrint(program)
		// Interpreter.InterpretProgram(program)
	}

	e := Codegen.NewAMD64AssemblyEmitter(Codegen.INTEL, Codegen.SYSYEM_V)
	// a.EmitMainFunction()
	e.EmitLoadIntegerConstant(6)
	e.EmitInstructions("file.s")
}
