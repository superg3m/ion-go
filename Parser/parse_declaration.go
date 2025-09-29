package Parser

import (
	"ion-go/AST"
	"ion-go/TS"
	"ion-go/Token"
)

func (parser *Parser) parseParameters() []TS.Parameter {
	var params []TS.Parameter

	parser.expect(Token.LEFT_PAREN)
	for !parser.consumeOnMatch(Token.RIGHT_PAREN) {
		param := parser.expect(Token.IDENTIFIER)
		parser.expect(Token.COLON)
		dataType := parser.parseType()

		params = append(params, TS.Parameter{
			Tok:      param,
			DeclType: dataType,
		})

		if parser.peekNthToken(0).Kind != Token.RIGHT_PAREN {
			parser.expect(Token.COMMA)
		}
	}

	return params
}

func (parser *Parser) parseMembers() []AST.Member {
	var params []AST.Member

	parser.expect(Token.LEFT_CURLY)
	for !parser.consumeOnMatch(Token.RIGHT_CURLY) {
		member := parser.expect(Token.IDENTIFIER)
		parser.expect(Token.COLON)
		dataType := parser.parseType()

		params = append(params, AST.Member{
			Tok:      member,
			DeclType: dataType,
		})

		if parser.peekNthToken(0).Kind != Token.RIGHT_CURLY {
			parser.expect(Token.COMMA)
		}
	}

	return params
}

func (parser *Parser) parseVariableDeclaration() AST.Declaration {
	parser.expect(Token.VAR)
	ident := parser.expect(Token.IDENTIFIER)
	parser.expect(Token.COLON)
	var dataType *TS.Type
	if parser.peekNthToken(0).Kind == Token.EQUALS {
		parser.expect(Token.EQUALS)
	} else {
		dataType = parser.parseType()
		parser.expect(Token.EQUALS)
	}

	rhs := parser.parseExpression()
	parser.expect(Token.SEMI_COLON)

	return &AST.DeclarationVariable{
		Tok:      ident,
		DeclType: dataType,
		RHS:      rhs,
	}
}

func (parser *Parser) parseFunctionDeclaration() AST.Declaration {
	parser.expect(Token.FN)
	ident := parser.expect(Token.IDENTIFIER)
	params := parser.parseParameters()
	parser.expect(Token.RIGHT_ARROW)
	returnType := parser.parseType()
	block := parser.parseStatementBlock().(*AST.StatementBlock)

	declType := TS.NewType(TS.FUNCTION, returnType, params)

	return &AST.DeclarationFunction{
		Tok:      ident,
		DeclType: declType,
		Block:    block,
	}
}

func (parser *Parser) parseStructDeclaration() AST.Declaration {
	parser.expect(Token.STRUCT)
	typeName := parser.expect(Token.IDENTIFIER)
	members := parser.parseMembers()

	memberLookup := make(map[string]AST.Member)
	for _, member := range members {
		memberLookup[member.Tok.Lexeme] = member
	}

	parser.ctx.ParsedStructDeclaration[typeName.Lexeme] = &AST.DeclarationStruct{
		Tok:          typeName,
		Members:      members,
		MemberLookup: memberLookup,
	}

	return parser.ctx.ParsedStructDeclaration[typeName.Lexeme]
}

func (parser *Parser) parseDeclaration() AST.Declaration {
	current := parser.peekNthToken(0)

	if current.Kind == Token.VAR {
		return parser.parseVariableDeclaration()
	} else if current.Kind == Token.FN {
		return parser.parseFunctionDeclaration()
	} else if current.Kind == Token.STRUCT {
		return parser.parseStructDeclaration()
	}

	return nil
}
