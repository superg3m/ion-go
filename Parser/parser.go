package Parser

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

type Parser struct {
	tokens  []Token.Token
	current int
}

func (parser *Parser) consumeNextToken() Token.Token {
	ret := parser.tokens[parser.current]
	parser.current += 1

	return ret
}

func (parser *Parser) peekNthToken(n int) Token.Token {
	return parser.tokens[parser.current+n]
}

func (parser *Parser) reportError(msg string) {
	nextToken := parser.peekNthToken(0)
	fmt.Printf("Parser Error: %s | Line: %d\n", nextToken.Lexeme, nextToken.Line)
	fmt.Printf("Msg: %s\n", msg)

	panic("")
}

func (parser *Parser) expect(expectedType Token.TokenType) Token.Token {
	if parser.peekNthToken(0).Kind != expectedType {
		msg := fmt.Sprintf("Expected: %s | Got: %s", string(expectedType), parser.peekNthToken(0).Lexeme)
		parser.reportError(msg)
	}

	return parser.consumeNextToken()
}

func (parser *Parser) consumeOnMatch(expectedType Token.TokenType) bool {
	if parser.peekNthToken(0).Kind == expectedType {
		parser.consumeNextToken()
		return true
	}

	return false
}

func (parser *Parser) previousToken() Token.Token {
	return parser.tokens[parser.current-1]
}

func (parser *Parser) parseDataType() AST.DataType {
	dataTypeToken := parser.expect(Token.IDENTIFIER)

	return AST.DataType{
		Name: dataTypeToken.Lexeme,
	}
}

func ParseProgram(tokens []Token.Token) AST.Program {
	parser := Parser{}
	parser.current = 0
	parser.tokens = tokens

	var program AST.Program
	for parser.current < (len(parser.tokens) - 1) {
		decl := parser.parseDeclaration()
		if decl == nil {
			parser.reportError("Unable to parse declaration")
		}

		program.Declarations = append(program.Declarations, decl)
	}

	return program
}
