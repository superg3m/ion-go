package Parser

import (
	"fmt"
	"ion-go/AST"
	"ion-go/TS"
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

func (parser *Parser) parseAccessChainExpression(token Token.Token) *AST.ExpressionAccessChain {
	var keys []AST.Expression

	inital_token := token

	next := parser.peekNthToken(0).Kind
	for next == Token.LEFT_BRACKET || next == Token.DOT {
		next = parser.peekNthToken(0).Kind

		if parser.consumeOnMatch(Token.DOT) {
			token = parser.expect(Token.IDENTIFIER)
			keys = append(keys, &AST.ExpressionIdentifier{
				Tok: token,
			})
		}

		if parser.consumeOnMatch(Token.LEFT_BRACKET) {
			keys = append(keys, &AST.ExpressionArrayAccess{
				Tok:   token,
				Index: parser.parseExpression(),
			})
			parser.expect(Token.RIGHT_BRACKET)
		}
	}

	return &AST.ExpressionAccessChain{
		Tok:        inital_token,
		AccessKeys: keys,
	}
}

// <Primary>    ::= <integer> | <float> | <boolean> | <string> | '(' <Expression> ')'
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
	} else if parser.consumeOnMatch(Token.STRING_LITERAL) {
		return &AST.ExpressionString{Value: current.Lexeme[1 : len(current.Lexeme)-1]}
	} else if parser.consumeOnMatch(Token.BUILTIN_LEN) {
		parser.expect(Token.LEFT_PAREN)
		iterable := parser.parseExpression()
		parser.expect(Token.RIGHT_PAREN)

		return &AST.ExpressionLen{
			Iterable: iterable,
		}
	} else if parser.consumeOnMatch(Token.IDENTIFIER) {
		next := parser.peekNthToken(0)
		if next.Kind == Token.DOT || next.Kind == Token.LEFT_BRACKET {
			return parser.parseAccessChainExpression(current)
		}

		if next.Kind == Token.LEFT_PAREN {
			arguments := parser.parseArguments()
			return &AST.SE_FunctionCall{
				Tok:       current,
				Arguments: arguments,
			}
		}

		return &AST.ExpressionIdentifier{
			Tok: current,
		}
	} else if parser.consumeOnMatch(Token.LEFT_PAREN) {
		expr := parser.parseExpression()
		if expr != nil {
			parser.expect(Token.RIGHT_PAREN)
			return &AST.ExpressionGrouping{
				Expr: expr,
			}
		}
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

// <comparison> ::= <additive> (('=='|'!='|<'|'<='|'>='|'>'} <additive>)*
func (parser *Parser) parseComparisonExpression() AST.Expression {
	expr := parser.parseAdditiveExpression()

	for parser.consumeOnMatch(Token.EQUALS_EQUALS) ||
		parser.consumeOnMatch(Token.NOT_EQUALS) ||
		parser.consumeOnMatch(Token.LESS_THAN) ||
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

	declType := TS.NewType(TS.ARRAY, nil, nil)
	if parser.ctx.ParsingArrayLiteral == 0 {
		declType = parser.parseType()
		parser.expect(Token.DOT)
	}

	parser.ctx.ParsingArrayLiteral += 1
	parser.expect(Token.LEFT_BRACKET)
	for !parser.consumeOnMatch(Token.RIGHT_BRACKET) {
		expr := parser.parseExpression()
		elements = append(elements, expr)

		if parser.peekNthToken(0).Kind != Token.RIGHT_BRACKET {
			parser.expect(Token.COMMA)
		}
	}
	parser.ctx.ParsingArrayLiteral -= 1

	return &AST.ExpressionArray{
		Elements: elements,
		DeclType: declType,
	}
}

// <struct> ::= <type>.{(<expression>,)*}
func (parser *Parser) parseStructExpression() AST.Expression {
	values := make(map[string]AST.Expression)

	typeName := parser.expect(Token.IDENTIFIER)
	structDecl, ok := parser.ctx.ParsedStructDeclaration[typeName.Lexeme]
	if !ok {
		panic(fmt.Sprintf("Line %d | Type %s is not defined", typeName.Line, typeName.Lexeme))
	}

	parser.expect(Token.DOT)
	parser.expect(Token.LEFT_CURLY)

	memberCount := 0

	for !parser.consumeOnMatch(Token.RIGHT_CURLY) {
		expr := parser.parseExpression()
		member := structDecl.Members[memberCount]
		values[member.Tok.Lexeme] = expr

		if parser.peekNthToken(0).Kind != Token.RIGHT_CURLY {
			parser.expect(Token.COMMA)
		}

		memberCount += 1
	}

	if memberCount != len(structDecl.Members) {
		panic(fmt.Sprintf("Line: %d | Expected members count to be: %d | Got: %d", typeName.Line, len(structDecl.Members), memberCount))
	}

	return &AST.ExpressionStruct{
		Tok:          typeName,
		MemberValues: values,
	}
}

// <Expression> ::= <additive>
func (parser *Parser) parseExpression() AST.Expression {
	current := parser.peekNthToken(0)
	next := parser.peekNthToken(1)
	next2 := parser.peekNthToken(2)

	if current.Kind == Token.LEFT_BRACKET {
		return parser.parseArrayExpression()
	} else if current.Kind == Token.IDENTIFIER && next.Kind == Token.DOT && next2.Kind == Token.LEFT_CURLY {
		return parser.parseStructExpression()
	} else if current.Kind == Token.CAST {
		cast := parser.expect(Token.CAST)
		parser.expect(Token.LEFT_PAREN)
		castType := parser.parseType()
		parser.expect(Token.RIGHT_PAREN)

		return &AST.ExpressionTypeCast{
			Tok:      cast,
			CastType: castType,
			Expr:     parser.parseExpression(),
		}
	} else {
		return parser.parseLogicalExpression()
	}
}
