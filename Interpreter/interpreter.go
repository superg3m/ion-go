package Interpreter

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

var global_functions map[string]AST.DeclarationFunction

func interpretBinaryExpression(kind Token.TokenType, left, right AST.Expression) AST.Expression {
	switch kind {
	case Token.PLUS, Token.MINUS, Token.STAR, Token.DIVISION,
		Token.LESS_THAN, Token.LESS_THAN_EQUALS, Token.GREATER_THAN, Token.GREATER_THAN_EQUALS:
		// Try int + int
		if lhs, ok1 := left.(AST.ExpressionInteger); ok1 {
			if rhs, ok2 := right.(AST.ExpressionInteger); ok2 {
				return evaluateIntegers(kind, lhs, rhs)
			}
		}

		// Try float + float
		if lhs, ok1 := left.(AST.ExpressionFloat); ok1 {
			if rhs, ok2 := right.(AST.ExpressionFloat); ok2 {
				return evaluateFloats(kind, lhs, rhs)
			}
		}

		// Mixed int + float (promote int to float)
		if lhs, ok1 := left.(AST.ExpressionInteger); ok1 {
			if rhs, ok2 := right.(AST.ExpressionFloat); ok2 {
				return evaluateFloats(kind, AST.ExpressionFloat(lhs), rhs)
			}
		}

		// Mixed float + int (promote int to float)
		if lhs, ok1 := left.(AST.ExpressionFloat); ok1 {
			if rhs, ok2 := right.(AST.ExpressionInteger); ok2 {
				return evaluateFloats(kind, lhs, AST.ExpressionFloat(rhs))
			}
		}

		panic(fmt.Sprintf("invalid operands for %v: %T and %T", kind, left, right))

	case Token.LOGICAL_AND, Token.LOGICAL_OR:
		lhs, ok1 := left.(AST.ExpressionBoolean)
		rhs, ok2 := right.(AST.ExpressionBoolean)
		if !ok1 || !ok2 {
			panic(fmt.Sprintf("expected booleans for %v, got %T and %T", kind, left, right))
		}

		if kind == Token.LOGICAL_AND {
			return lhs && rhs
		} else {
			return lhs || rhs
		}

	default:
		panic(fmt.Sprintf("unhandled operator: %v", kind))
	}
}

func evaluateIntegers(kind Token.TokenType, lhs, rhs AST.ExpressionInteger) AST.Expression {
	switch kind {
	case Token.PLUS:
		return lhs + rhs
	case Token.MINUS:
		return lhs - rhs
	case Token.STAR:
		return lhs * rhs
	case Token.DIVISION:
		return lhs / rhs
	case Token.EQUALS_EQUALS:
		return AST.ExpressionBoolean(lhs == rhs)
	case Token.LESS_THAN:
		return AST.ExpressionBoolean(lhs < rhs)
	case Token.LESS_THAN_EQUALS:
		return AST.ExpressionBoolean(lhs <= rhs)
	case Token.GREATER_THAN:
		return AST.ExpressionBoolean(lhs > rhs)
	case Token.GREATER_THAN_EQUALS:
		return AST.ExpressionBoolean(lhs >= rhs)
	}

	panic("unreachable")
}

func evaluateFloats(kind Token.TokenType, lhs, rhs AST.ExpressionFloat) AST.Expression {
	switch kind {
	case Token.PLUS:
		return lhs + rhs
	case Token.MINUS:
		return lhs - rhs
	case Token.STAR:
		return lhs * rhs
	case Token.DIVISION:
		return lhs / rhs
	case Token.LESS_THAN:
		return AST.ExpressionBoolean(lhs < rhs)
	case Token.LESS_THAN_EQUALS:
		return AST.ExpressionBoolean(lhs <= rhs)
	case Token.GREATER_THAN:
		return AST.ExpressionBoolean(lhs > rhs)
	case Token.GREATER_THAN_EQUALS:
		return AST.ExpressionBoolean(lhs >= rhs)
	}
	panic("unreachable")
}

func interpretExpression(e AST.Expression, scope *Scope) AST.Expression {
	if e == nil {
		return nil
	}

	switch v := e.(type) {
	case AST.ExpressionInteger, AST.ExpressionFloat, AST.ExpressionBoolean:
		return v
	case AST.ExpressionIdentifier:
		return scope.get(v.Name)
	case AST.ExpressionFunctionCall:
		functionDeclaration := global_functions[v.Name]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		functionScope := CreateScope(scope)
		for i := 0; i < argCount; i++ {
			param := functionDeclaration.Parameters[i]
			arg := v.Arguments[i]
			functionScope.set(param.Name, interpretExpression(arg, scope))
		}

		return interpretNodes(functionDeclaration.Block.Body, &functionScope)

	case AST.ExpressionLen:
		v.Array = interpretExpression(v.Array, scope)
		return AST.ExpressionInteger(len(v.Array.(AST.ExpressionArray).Elements))
	case AST.ExpressionGrouping:
		return interpretExpression(v.Expr, scope)
	case AST.ExpressionBinary:
		leftExpression := interpretExpression(v.Left, scope)
		rightExpression := interpretExpression(v.Right, scope)
		return interpretBinaryExpression(v.Operator.Kind, leftExpression, rightExpression)
	case AST.ExpressionArray:
		return v
	case AST.ExpressionArrayAccess:
		arr := scope.get(v.Name).(AST.ExpressionArray)
		index := interpretExpression(v.Index, scope).(AST.ExpressionInteger)
		return arr.Elements[index]

	default:
		fmt.Printf("Type: %T\n", e)
		panic("unreachable")
	}
}

func interpretDeclaration(decl AST.Declaration, scope *Scope) {
	switch v := decl.(type) {
	case AST.DeclarationVariable:
		if scope.has(v.Name) {
			panic("Attempting to redeclare: " + v.Name)
		}

		v.RHS = interpretExpression(v.RHS, scope)
		if v.RHS == nil {
			panic("Attempting to assign void to variable: " + v.Name)
		}

		scope.set(v.Name, v.RHS)

	case AST.DeclarationFunction:
		global_functions[v.Name] = v

	}
}

func interpretStatement(s AST.Statement, scope *Scope) AST.Expression {
	switch v := s.(type) {
	case AST.StatementPrint:
		fmt.Println(interpretExpression(v.Expr, scope))

	case AST.StatementAssignment:
		if !scope.has(v.Name) {
			panic("Attempting to assign to undeclared identifier: " + v.Name)
		}

		v.RHS = interpretExpression(v.RHS, scope)
		if v.RHS == nil {
			panic("Attempting to assign void to variable: " + v.Name)
		}

		scope.set(v.Name, v.RHS)

	case AST.StatementBlock:
		blockScope := CreateScope(scope)
		interpretNodes(v.Body, &blockScope)

	case AST.StatementFor:
		forScope := CreateScope(scope)
		interpretDeclaration(v.Initializer, &forScope)
		for interpretExpression(v.Condition, &forScope).(AST.ExpressionBoolean) {
			interpretStatement(v.Body, &forScope)
			interpretStatement(v.Increment, &forScope)
		}

	case AST.StatementReturn:
		return interpretExpression(v.Expr, scope)

	default:
		panic("unreachable")
	}

	return nil
}

func interpretNodes(nodes []AST.Node, scope *Scope) AST.Expression {
	for _, node := range nodes {
		switch v := node.(type) {
		case AST.Statement:
			ret := interpretStatement(v, scope)
			if ret != nil {
				return ret
			}
		case AST.Declaration:
			interpretDeclaration(v, scope)
		}

	}

	return nil
}

func InterpretProgram(program AST.Program) {
	globalScope := CreateScope(nil)
	global_functions = make(map[string]AST.DeclarationFunction)

	for _, decl := range program.Declarations {
		interpretDeclaration(decl, &globalScope)
	}

	if main_decl, ok := global_functions["main"]; ok {
		main_call := AST.ExpressionFunctionCall{
			Name: main_decl.Name,
		}

		interpretExpression(main_call, &globalScope)
	} else {
		panic("main function not found")
	}
}
