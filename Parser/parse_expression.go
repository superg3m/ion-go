package Parser

import (
	"ion-go/AST"
	"ion-go/Token"
	"strconv"
)

// call site
func (parser *Parser) parseArguments() []AST.Expression {
	var ret []AST.Expression

	parser.expect(Token.LEFT_PAREN)
	for !parser.consumeOnMatch(Token.RIGHT_PAREN) {
		expression := parser.parseExpression()

		ret = append(ret, expression)

		if parser.peekNthToken(0).Kind != Token.RIGHT_PAREN {
			parser.expect(Token.COMMA)
		}
	}

	return ret
}

// <Primary>    ::= <integer> | <float> | <boolean> | '(' <Expression> ')'
func (parser *Parser) parsePrimary() AST.Expression {
	current := parser.peekNthToken(0)
	if parser.consumeOnMatch(Token.INTEGER_LITERAL) {
		num, _ := strconv.Atoi(current.Lexeme)
		return &AST.ExpressionInteger{Value: num}
	} else if parser.consumeOnMatch(Token.BOOLEAN_LITERAL) {
		b := current.Lexeme == "true"
		return &AST.ExpressionBoolean{Value: b}
	} else if parser.consumeOnMatch(Token.FLOAT_LITERAL) {
		num, _ := strconv.ParseFloat(current.Lexeme, 32)
		return &AST.ExpressionFloat{Value: float32(num)}
	} else if parser.consumeOnMatch(Token.BUILTIN_LEN) {
		parser.expect(Token.LEFT_PAREN)
		arr := parser.parseExpression()
		parser.expect(Token.RIGHT_PAREN)

		return &AST.ExpressionLen{
			Array: arr,
		}
	} else if parser.consumeOnMatch(Token.IDENTIFIER) {
		next := parser.peekNthToken(0)
		if next.Kind == Token.LEFT_BRACKET {
			parser.expect(Token.LEFT_BRACKET)
			index := parser.parseExpression()
			parser.expect(Token.RIGHT_BRACKET)
			return &AST.ExpressionArrayAccess{
				Name:  current.Lexeme,
				Index: index,
			}
		}

		if next.Kind == Token.LEFT_PAREN {
			arguments := parser.parseArguments()
			return &AST.ExpressionFunctionCall{
				Name:      current.Lexeme,
				Arguments: arguments,
			}
		}

		return &AST.ExpressionIdentifier{
			Name: current.Lexeme,
		}
	} else if parser.consumeOnMatch(Token.LEFT_PAREN) {
		ret := &AST.ExpressionGrouping{
			Expr: parser.parseExpression(),
		}
		parser.expect(Token.RIGHT_PAREN)

		return ret
	}

	return nil
}

// <Unary>      ::= ('+'|'-'|'!') <unary> | <Primary>
func (parser *Parser) parseUnaryExpression() AST.Expression {
	ret := &AST.ExpressionUnary{}

	if parser.consumeOnMatch(Token.NOT) || parser.consumeOnMatch(Token.MINUS) || parser.consumeOnMatch(Token.PLUS) {
		ret.Operator = parser.previousToken()
		ret.Operand = parser.parseUnaryExpression()

		return ret
	}

	return parser.parsePrimary()
}

// <multiplicative>     ::= <Unary> (('*'|'/') <Unary>)*
func (parser *Parser) parseMultiplicativeExpression() AST.Expression {
	expr := parser.parseUnaryExpression()

	for parser.consumeOnMatch(Token.STAR) || parser.consumeOnMatch(Token.DIVISION) {
		op := parser.previousToken()
		right := parser.parseUnaryExpression()
		expr = &AST.ExpressionBinary{
			Operator: op,
			Left:     expr,
			Right:    right,
		}
	}

	return expr
}

// <additive>       ::= <Factor> (('+'|'-') <Factor>)*
func (parser *Parser) parseAdditiveExpression() AST.Expression {
	expr := parser.parseMultiplicativeExpression()

	for parser.consumeOnMatch(Token.PLUS) || parser.consumeOnMatch(Token.MINUS) {
		op := parser.previousToken()
		right := parser.parseMultiplicativeExpression()
		expr = &AST.ExpressionBinary{
			Operator: op,
			Left:     expr,
			Right:    right,
		}
	}

	return expr
}

// <comparison> ::= <additive> (('<'|'<='|'>='|'>'} <additive>)*
func (parser *Parser) parseComparisonExpression() AST.Expression {
	expr := parser.parseAdditiveExpression()

	for parser.consumeOnMatch(Token.LESS_THAN) ||
		parser.consumeOnMatch(Token.LESS_THAN_EQUALS) ||
		parser.consumeOnMatch(Token.GREATER_THAN_EQUALS) ||
		parser.consumeOnMatch(Token.GREATER_THAN) {
		op := parser.previousToken()
		right := parser.parseAdditiveExpression()
		expr = &AST.ExpressionBinary{
			Operator: op,
			Left:     expr,
			Right:    right,
		}
	}

	return expr
}

// <logical> ::= <comparison> (('&&'|'||') <comparison>)*
func (parser *Parser) parseLogicalExpression() AST.Expression {
	expr := parser.parseComparisonExpression()

	for parser.consumeOnMatch(Token.LOGICAL_AND) || parser.consumeOnMatch(Token.LOGICAL_OR) {
		op := parser.previousToken()
		right := parser.parseComparisonExpression()
		expr = &AST.ExpressionBinary{
			Operator: op,
			Left:     expr,
			Right:    right,
		}
	}

	return expr
}

// <array> ::= <type>.[(<expression>,)*]
func (parser *Parser) parseArrayExpression() AST.Expression {
	var elements []AST.Expression

	dataType := parser.parseDataType()
	parser.expect(Token.DOT)
	parser.expect(Token.LEFT_BRACKET)
	for !parser.consumeOnMatch(Token.RIGHT_BRACKET) {
		expr := parser.parseExpression()
		elements = append(elements, expr)

		if parser.peekNthToken(0).Kind != Token.RIGHT_BRACKET {
			parser.expect(Token.COMMA)
		}
	}

	return &AST.ExpressionArray{
		Elements: elements,
		DeclType: dataType,
	}
}

// <Expression> ::= <additive>
func (parser *Parser) parseExpression() AST.Expression {
	current := parser.peekNthToken(0)
	next := parser.peekNthToken(1)
	next2 := parser.peekNthToken(2)
	if current.Kind == Token.IDENTIFIER &&
		next.Kind == Token.DOT && next2.Kind == Token.LEFT_BRACKET {
		return parser.parseArrayExpression()
	} else {
		return parser.parseLogicalExpression()
	}
}
