package Token

type TokenType string

const (
	ILLEGAL_TOKEN TokenType = "ILLEGAL_TOKEN"
	EOF                     = "EOF"

	// SYNTAX SINGLE CHARACTER
	NOT           = "NOT"           // "!"
	LEFT_PAREN    = "LEFT_PAREN"    // "("
	RIGHT_PAREN   = "RIGHT_PAREN"   // ")"
	PLUS          = "PLUS"          // "+"
	MINUS         = "MINUS"         // "-"
	STAR          = "STAR"          // "*"
	DIVISION      = "DIVISION"      // "/"
	MODULUS       = "MODULUS"       // "%"
	COMMA         = "COMMA"         // ","
	COLON         = "COLON"         // ":"
	DOT           = "DOT"           // "."
	SEMI_COLON    = "SEMI_COLON"    // ";"
	LESS_THAN     = "LESS_THAN"     // "<"
	EQUALS        = "EQUAL"         // "="
	GREATER_THAN  = "GREATER_THAN"  // ">"
	LEFT_BRACKET  = "LEFT_BRACKET"  // "["
	RIGHT_BRACKET = "RIGHT_BRACKET" // "]"
	LEFT_CURLY    = "LEFT_CURLY"    // "{"
	RIGHT_CURLY   = "RIGHT_CURLY"   // "}"

	// SYNTAX MULTIPLE CHARACTERS
	EQUALS_EQUALS       = "EQUALS_EQUALS"       // "=="
	NOT_EQUALS          = "NOT_EQUALS"          // "!="
	GREATER_THAN_EQUALS = "GREATER_THAN_EQUALS" // ">="
	LESS_THAN_EQUALS    = "LESS_THAN_EQUALS"    // "<="
	LOGICAL_AND         = "LOGICAL_AND"         // "&&"
	LOGICAL_OR          = "LOGICAL_OR"          // "||"
	RIGHT_ARROW         = "RIGHT_ARROW"         // "->"

	IDENTIFIER        = "IDENTIFIER"
	INTEGER_LITERAL   = "INTEGER_LITERAL"
	FLOAT_LITERAL     = "FLOAT_LITERAL"
	BOOLEAN_LITERAL   = "BOOLEAN_LITERAL"
	STRING_LITERAL    = "STRING_LITERAL"
	CHARACTER_LITERAL = "CHARACTER_LITERAL"

	// Keywords
	FN       = "FN"
	CAST     = "CAST"
	VAR      = "VAR"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	WHILE    = "WHILE"
	NULLPTR  = "NULLPTR"
	RETURN   = "RETURN"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	PRINT    = "PRINT"
	PRINTLN  = "PRINTLN"
	DEFER    = "DEFER"

	// Builtin
	BUILTIN_LEN = "BUILTIN_LEN"
)

type Token struct {
	Kind   TokenType
	Lexeme string
	Line   int
}

func CreateToken(kind TokenType, lexeme string, line int) Token {
	return Token{kind, lexeme, line}
}

func GetKeywordToken(input string) (TokenType, bool) {
	var m = map[string]TokenType{
		"fn":       FN,
		"cast":     CAST,
		"var":      VAR,
		"if":       IF,
		"else":     ELSE,
		"for":      FOR,
		"while":    WHILE,
		"nullptr":  NULLPTR,
		"return":   RETURN,
		"break":    BREAK,
		"continue": CONTINUE,
		"print":    PRINT,
		"println":  PRINTLN,
		"defer":    DEFER,
		"true":     BOOLEAN_LITERAL,
		"false":    BOOLEAN_LITERAL,
	}

	token, ok := m[input]

	return token, ok
}

func GetBuiltinToken(input string) (TokenType, bool) {
	var m = map[string]TokenType{
		"len": BUILTIN_LEN,
	}

	token, ok := m[input]

	return token, ok
}

func GetSyntaxToken(input string) (TokenType, bool) {
	var m = map[string]TokenType{
		"=":  EQUALS,
		"%":  MODULUS,
		"+":  PLUS,
		"-":  MINUS,
		"/":  DIVISION,
		"*":  STAR,
		"<":  LESS_THAN,
		">":  GREATER_THAN,
		"!":  NOT,
		"(":  LEFT_PAREN,
		")":  RIGHT_PAREN,
		",":  COMMA,
		":":  COLON,
		".":  DOT,
		";":  SEMI_COLON,
		"[":  LEFT_BRACKET,
		"]":  RIGHT_BRACKET,
		"{":  LEFT_CURLY,
		"}":  RIGHT_CURLY,
		"==": EQUALS_EQUALS,
		"!=": NOT_EQUALS,
		">=": GREATER_THAN_EQUALS,
		"<=": LESS_THAN_EQUALS,
		"&&": LOGICAL_AND,
		"||": LOGICAL_OR,
		"->": RIGHT_ARROW,
	}

	token, ok := m[input]

	return token, ok
}
