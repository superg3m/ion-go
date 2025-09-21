package Lexer

import (
	"fmt"
	"ion-go/Token"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	leftPos  int
	rightPos int
	line     int
	c        byte
	source   []byte
	tokens   []Token.Token
}

func createLexer() Lexer {
	return Lexer{
		leftPos:  0,
		rightPos: 0,
		line:     1,
		c:        0,
		source:   []byte{},
		tokens:   []Token.Token{},
	}
}

func (lexer *Lexer) isEOF() bool {
	return lexer.rightPos >= len(lexer.source)
}

func (lexer *Lexer) getScratchBuffer() string {
	return strings.TrimSpace(string(lexer.source[lexer.leftPos:lexer.rightPos]))
}

func (lexer *Lexer) consumeUntilNewLine() {
	for lexer.c != '\n' && !lexer.isEOF() {
		lexer.consumeNextChar()
	}
}

func (lexer *Lexer) reportError(format string, args ...interface{}) {
	fmt.Println("String:", string(lexer.getScratchBuffer()))
	fmt.Println("Data:", []byte(lexer.getScratchBuffer()))
	fmt.Println("Error Line:", lexer.line, "|", fmt.Sprintf(format, args...))

	panic("")
}

func (lexer *Lexer) consumeNextChar() {
	lexer.c = lexer.source[lexer.rightPos]
	lexer.rightPos += 1
}

func (lexer *Lexer) peekNthChar(n int) byte {
	if lexer.rightPos+n >= len(lexer.source) {
		return 0
	}

	return lexer.source[lexer.rightPos+n]
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == 0
}

func (lexer *Lexer) consumeWhitespace() bool {
	if !isWhitespace(lexer.c) {
		return false
	}

	if lexer.c == '\n' {
		lexer.line += 1
	}

	return true
}

func (lexer *Lexer) consumeWord() bool {
	if !unicode.IsLetter(rune(lexer.c)) {
		return false
	}

	isAlphaNumeric := func(c rune) bool {
		return unicode.IsDigit(c) || unicode.IsLetter(c) || c == '_'
	}

	for isAlphaNumeric(rune(lexer.peekNthChar(0))) {
		lexer.consumeNextChar()
	}

	return true
}

func (lexer *Lexer) consumeNextToken() {
	lexer.leftPos = lexer.rightPos
	lexer.consumeNextChar()

	if lexer.consumeWhitespace() {
	} else if lexer.consumeLiteral() {
	} else if lexer.consumeIdentifier() {
	} else if lexer.consumeSyntax() {
	} else {
		lexer.reportError("Illegal token found")
	}
}

func (lexer *Lexer) consumeOnMatch(expected byte) bool {
	if lexer.peekNthChar(0) != expected {
		return false
	}

	lexer.consumeNextChar()
	return true
}

func (lexer *Lexer) consumeSyntax() bool {
	switch lexer.c {
	case '&':
		if lexer.consumeOnMatch('&') {
		}

	case '[':
		if lexer.consumeOnMatch('.') {
			lexer.consumeNextChar() // should be .
			lexer.consumeNextChar() // should be ]
		}

	case '<':
		if lexer.consumeOnMatch('<') {
			lexer.consumeOnMatch('=')
		} else if lexer.consumeOnMatch('=') {
		}

	case '>':
		if lexer.consumeOnMatch('>') {
			lexer.consumeOnMatch('=')
		} else if lexer.consumeOnMatch('=') {
		}

	case '-':
		if lexer.consumeOnMatch('=') {
		} else if lexer.consumeOnMatch('-') {
		} else if lexer.consumeOnMatch('>') {
		}

	case '+':
		if lexer.consumeOnMatch('=') {
		} else {
			lexer.consumeOnMatch('+')
		}

	case '/':
		if lexer.consumeOnMatch('=') {
		} else if lexer.consumeOnMatch('/') {
			lexer.consumeUntilNewLine()
			// lexer.addToken(Token.COMMENT)
			return true
		} else if lexer.consumeOnMatch('*') {
			for !(lexer.peekNthChar(0) == '*' && lexer.peekNthChar(1) == '/') {
				if lexer.isEOF() {
					panic("Multiline comment doesn't terminate\n")
				}

				lexer.consumeNextChar()
			}

			lexer.consumeNextChar() // consume '*'
			lexer.consumeNextChar() // consume '/'
			// lexer_add_token(lexer, SPL_TOKEN_COMMENT);
			return true
		}

	case '!', '*', '=':
		lexer.consumeOnMatch('=')
	}

	kind, ok := Token.GetSyntaxToken(lexer.getScratchBuffer())
	if ok {
		lexer.addToken(kind)
		return true
	}

	return false
}

func (lexer *Lexer) consumeIdentifier() bool {
	if !lexer.consumeWord() {
		return false
	}

	kind, ok := Token.GetKeywordToken(lexer.getScratchBuffer())
	if ok {
		lexer.addToken(kind)
		return true
	}

	kind, ok = Token.GetBuiltinToken(lexer.getScratchBuffer())
	if ok {
		lexer.addToken(kind)
		return true
	}

	lexer.addToken(Token.IDENTIFIER)

	return true
}

func (lexer *Lexer) addToken(kind Token.TokenType) {
	lexer.tokens = append(lexer.tokens, Token.CreateToken(kind, lexer.getScratchBuffer(), lexer.line))
}

/*
func (lexer *Lexer) tryConsumeStringLiteral() {
	for !lexer.consumeOnMatch('"') {
		if lexer.isEOF() {
			lexer.reportError("String literal doesn't have a closing double quote!")
		}
	}

	lexer.consumeNextChar()
	lexer.addToken(Token.STRING_LITERAL)
}


func (lexer *Lexer) tryConsumeCharacterLiteral() {
	if lexer.consumeOnMatch('\'') {
		lexer.reportError("Character literal doesn't have any ascii data in between")
	}

	for !lexer.consumeOnMatch('\'') {
		if lexer.isEOF() {
			lexer.reportError("character literal doesn't have a closing quote!")
		}

		lexer.consumeNextChar()
	}

	lexer.addToken(Token.CHARACTER_LITERAL)
}
*/

func (lexer *Lexer) tryConsumeDigitLiteral() {
	var kind Token.TokenType = Token.INTEGER_LITERAL

	for unicode.IsDigit(rune(lexer.peekNthChar(0))) || lexer.peekNthChar(0) == '.' {
		if lexer.c == '.' {
			kind = Token.FLOAT_LITERAL
		}

		lexer.consumeNextChar()
	}

	lexer.addToken(kind)
}

func (lexer *Lexer) consumeLiteral() bool {
	if unicode.IsDigit(rune(lexer.c)) {
		lexer.tryConsumeDigitLiteral()
		return true
	} else if lexer.c == '"' {
		// lexer.tryConsumeStringLiteral()
		panic("NOT IMPLEMENTED STRING LITERALS")

		return true
	} else if lexer.c == '\'' {
		panic("NOT IMPLEMENTED CHARACTER LITERALS")

		// lexer.tryConsumeCharacterLiteral()
		return true
	} else {
		return false
	}
}

func GenerateTokenStream(filePath string) []Token.Token {
	lexer := createLexer()

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic("File not found")
	}

	lexer.source = data

	for !lexer.isEOF() {
		lexer.consumeNextToken()
	}

	lexer.tokens = append(lexer.tokens, Token.CreateToken(Token.EOF, "", lexer.line))

	return lexer.tokens
}
