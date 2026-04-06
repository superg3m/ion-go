package main

import (
	"fmt"
	"ion-go/Interpreter"
	"ion-go/JSON"
	"ion-go/Lexer"
	"ion-go/Parser"
	"ion-go/TypeChecker"
)

func main() {
	// tokenStream := Lexer.GenerateTokenStream("./factorial.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fib.ion")
	tokenStream := Lexer.GenerateTokenStream("./array.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fractal.ion")
	// tokenStream := Lexer.GenerateTokenStream("./test.ion")
	// tokenStream := Lexer.GenerateTokenStream("./struct.ion")

	// SPL_Test
	// iterate through these have like
	// 001_implicit_cast_binary_op.spl
	// 002_implicit_cast_function_arg.spl
	// 003_multi_dimension_static_array_.spl

	for i := 0; i < len(tokenStream); i++ {
		token := tokenStream[i]
		tokenType, tokenValue := token.Kind, token.Lexeme
		fmt.Print("Type: ", tokenType, "(", tokenValue, ") | Line:", token.Line, "\n")
	}

	program := Parser.ParseProgram(tokenStream)
	//fmt.Printf("%+v\n", program)

	TypeChecker.TypeCheckProgram(program)
	JSON.PrettyPrint(program)
	Interpreter.InterpretProgram(program)
}
