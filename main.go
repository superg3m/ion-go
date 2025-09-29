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
	// tokenStream := Lexer.GenerateTokenStream("./test.ion")
	// tokenStream := Lexer.GenerateTokenStream("./factorial.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fractal.ion")
	// tokenStream := Lexer.GenerateTokenStream("./fib.ion")
	// tokenStream := Lexer.GenerateTokenStream("./array.ion")
	tokenStream := Lexer.GenerateTokenStream("./struct.ion")

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
