package AST

import "ion-go/Token"

type StatementExpression interface {
	Statement
	Expression
	isStatementExpression()
}

type SE_FunctionCall struct {
	Tok       Token.Token
	Arguments []Expression
}

func (*SE_FunctionCall) isNode()                {}
func (*SE_FunctionCall) isStatementExpression() {}
func (*SE_FunctionCall) isExpression()          {}
func (*SE_FunctionCall) isStatement()           {}
func (*SE_FunctionCall) isDeferrable()          {}
