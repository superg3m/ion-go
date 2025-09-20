package Parser

import (
	"ion-go/AST"
	"ion-go/Token"
)

func (parser *Parser) parsePrintStatement() AST.Statement {
	parser.expect(Token.PRINT)
	parser.expect(Token.LEFT_PAREN)
	expr := parser.parseExpression()
	parser.expect(Token.RIGHT_PAREN)
	parser.expect(Token.SEMI_COLON)

	return &AST.StatementPrint{
		Expr: expr,
	}
}
func (parser *Parser) parseStatementBlock() AST.Statement {
	var body []AST.Node
	parser.expect(Token.LEFT_CURLY)
	for !parser.consumeOnMatch(Token.RIGHT_CURLY) {
		if decl := parser.parseDeclaration(); decl != nil {
			body = append(body, decl)
			continue
		}

		if stmt := parser.parseStatement(); stmt != nil {
			body = append(body, stmt)
			continue
		}
	}

	return &AST.StatementBlock{
		Body: body,
	}
}

func (parser *Parser) parseAssignmentStatement() AST.Statement {
	ident := parser.expect(Token.IDENTIFIER)
	parser.expect(Token.EQUALS)
	rhs := parser.parseExpression()
	parser.expect(Token.SEMI_COLON)

	return &AST.StatementAssignment{
		Name: ident.Lexeme,
		RHS:  rhs,
	}
}
func (parser *Parser) parseForStatement() AST.Statement {
	parser.expect(Token.FOR)
	parser.expect(Token.LEFT_PAREN)
	initializer := parser.parseVariableDeclaration()
	parser.expect(Token.SEMI_COLON)
	condition := parser.parseExpression()
	parser.expect(Token.SEMI_COLON)
	increment := parser.parseAssignmentStatement()
	parser.expect(Token.RIGHT_PAREN)
	body := parser.parseStatement()

	return &AST.StatementFor{
		Initializer: initializer.(*AST.DeclarationVariable),
		Condition:   condition,
		Increment:   increment.(*AST.StatementAssignment),
		Body:        body.(*AST.StatementBlock),
	}
}

func (parser *Parser) parseStatement() AST.Statement {
	current := parser.peekNthToken(0)

	if current.Kind == Token.LEFT_CURLY {
		return parser.parseStatementBlock()
	} else if current.Kind == Token.IDENTIFIER {
		return parser.parseAssignmentStatement()
	} else if current.Kind == Token.PRINT {
		return parser.parsePrintStatement()
	} else if current.Kind == Token.RETURN {
		parser.expect(Token.RETURN)
		expr := parser.parseExpression()
		parser.expect(Token.SEMI_COLON)

		return &AST.StatementReturn{
			Expr: expr,
		}
	} else if current.Kind == Token.FOR {
		return parser.parseForStatement()
	}

	panic("INVALID STATEMENT!")
	return nil
}
