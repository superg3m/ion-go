package Parser

import (
	"ion-go/AST"
	"ion-go/Token"
)

func (parser *Parser) parseParameters() []AST.Parameter {
	var args []AST.Parameter

	parser.expect(Token.LEFT_PAREN)
	for !parser.consumeOnMatch(Token.RIGHT_PAREN) {
		param := parser.expect(Token.IDENTIFIER)
		parser.expect(Token.COLON)
		dataType := parser.parseDataType()

		args = append(args, AST.Parameter{
			Name:     param.Lexeme,
			DeclType: dataType,
		})

		if parser.peekNthToken(0).Kind != Token.RIGHT_PAREN {
			parser.expect(Token.COMMA)
		}
	}

	return args
}

func (parser *Parser) parseVariableDeclaration() AST.Declaration {
	parser.expect(Token.VAR)
	ident := parser.expect(Token.IDENTIFIER)
	parser.expect(Token.COLON)
	var dataType AST.DataType
	if parser.peekNthToken(0).Kind == Token.EQUALS {
		parser.expect(Token.EQUALS)
	} else {
		dataType = parser.parseDataType()
		parser.expect(Token.EQUALS)
	}

	rhs := parser.parseExpression()
	parser.expect(Token.SEMI_COLON)

	return &AST.DeclarationVariable{
		Name:     ident.Lexeme,
		DeclType: dataType,
		RHS:      rhs,
	}
}

func (parser *Parser) parseFunctionDeclaration() AST.Declaration {
	parser.expect(Token.FN)
	ident := parser.expect(Token.IDENTIFIER)
	params := parser.parseParameters()
	parser.expect(Token.RIGHT_ARROW)
	dataType := parser.parseDataType()
	block := parser.parseStatementBlock().(*AST.StatementBlock)

	return &AST.DeclarationFunction{
		Name:       ident.Lexeme,
		Parameters: params,
		ReturnType: dataType,
		Block:      block,
	}
}

func (parser *Parser) parseDeclaration() AST.Declaration {
	current := parser.peekNthToken(0)

	if current.Kind == Token.VAR {
		return parser.parseVariableDeclaration()
	} else if current.Kind == Token.FN {
		return parser.parseFunctionDeclaration()
	}

	return nil
}
