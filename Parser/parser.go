package Parser

import (
	"fmt"
	"ion-go/AST"
	"ion-go/TS"
	"ion-go/Token"
)

type Context struct {
	ParsingForIncrement     bool
	ParsingArrayLiteral     int
	ParsedStructDeclaration map[string]*AST.DeclarationStruct
}

type Parser struct {
	tokens  []Token.Token
	current int

	ctx Context
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

func (parser *Parser) parseType() TS.Type {
	var countArray []int
	for parser.peekNthToken(0).Kind == Token.LEFT_BRACKET {
		parser.consumeOnMatch(Token.LEFT_BRACKET)
		count := parser.parseExpression().(*AST.ExpressionInteger).Value // this should be like a [4]int
		countArray = append(countArray, count)
		parser.consumeOnMatch(Token.RIGHT_BRACKET)
	}

	// TODO(Jovanni): handle pointer parsing here

	next := parser.peekNthToken(0)
	if next.Kind != Token.IDENTIFIER {
		return nil
	}

	dataTypeToken := parser.expect(Token.IDENTIFIER)
	var retType TS.Type
	if v, ok := TS.GetBuiltin(dataTypeToken.Lexeme); ok {
		retType = v
	} else if _, ok2 := parser.ctx.ParsedStructDeclaration[dataTypeToken.Lexeme]; ok2 {
		retType = TS.NewTypeStruct(dataTypeToken.Lexeme, nil)
	} else {
		parser.reportError(fmt.Sprintf("Line: %d, Unrecognized type: %s", dataTypeToken.Line, dataTypeToken.Lexeme))
	}

	for _, count := range countArray {
		retType = retType.AddStaticArray(count)
	}

	return retType
}

func ParseProgram(tokens []Token.Token) AST.Program {
	parser := Parser{}
	parser.current = 0
	parser.tokens = tokens
	parser.ctx.ParsedStructDeclaration = make(map[string]*AST.DeclarationStruct)

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
