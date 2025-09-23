package Interpreter

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

var globalFunctions map[string]*AST.DeclarationFunction
var globalScope Scope

func interpretBinaryExpression(kind Token.TokenType, left, right AST.Expression) AST.Expression {
	switch kind {
	case Token.PLUS, Token.MINUS, Token.STAR, Token.DIVISION,
		Token.LESS_THAN, Token.LESS_THAN_EQUALS, Token.GREATER_THAN, Token.GREATER_THAN_EQUALS,
		Token.EQUALS_EQUALS, Token.NOT_EQUALS:
		// Try int + int
		if lhs, ok1 := left.(*AST.ExpressionInteger); ok1 {
			if rhs, ok2 := right.(*AST.ExpressionInteger); ok2 {
				return evaluateIntegers(kind, lhs.Value, rhs.Value)
			}
		}

		// Try float + float
		if lhs, ok1 := left.(*AST.ExpressionFloat); ok1 {
			if rhs, ok2 := right.(*AST.ExpressionFloat); ok2 {
				return evaluateFloats(kind, lhs.Value, rhs.Value)
			}
		}

		// Mixed int + float (promote int to float)
		if lhs, ok1 := left.(*AST.ExpressionInteger); ok1 {
			if rhs, ok2 := right.(*AST.ExpressionFloat); ok2 {
				return evaluateFloats(kind, float32(lhs.Value), rhs.Value)
			}
		}

		// Mixed float + int (promote int to float)
		if lhs, ok1 := left.(*AST.ExpressionFloat); ok1 {
			if rhs, ok2 := right.(*AST.ExpressionInteger); ok2 {
				return evaluateFloats(kind, lhs.Value, float32(rhs.Value))
			}
		}

		if lhs, ok1 := left.(*AST.ExpressionString); ok1 {
			switch rhs := right.(type) {
			case *AST.ExpressionInteger:
				return evaluateString(kind, lhs.Value, fmt.Sprintf("%d", rhs.Value))
			case *AST.ExpressionFloat:
				return evaluateString(kind, lhs.Value, fmt.Sprintf("%f", rhs.Value))
			case *AST.ExpressionString:
				return evaluateString(kind, lhs.Value, rhs.Value)
			}
		}

		if rhs, ok1 := right.(*AST.ExpressionString); ok1 {
			switch lhs := right.(type) {
			case *AST.ExpressionInteger:
				return evaluateString(kind, fmt.Sprintf("%d", lhs.Value), rhs.Value)
			case *AST.ExpressionFloat:
				return evaluateString(kind, fmt.Sprintf("%f", lhs.Value), rhs.Value)
			case *AST.ExpressionString:
				return evaluateString(kind, lhs.Value, rhs.Value)
			}
		}

		panic(fmt.Sprintf("invalid operands for %v: %T and %T", kind, left, right))

	case Token.LOGICAL_AND, Token.LOGICAL_OR:
		lhs, ok1 := left.(*AST.ExpressionBoolean)
		rhs, ok2 := right.(*AST.ExpressionBoolean)
		if !ok1 || !ok2 {
			panic(fmt.Sprintf("expected booleans for %v, got %T and %T", kind, left, right))
		}

		if kind == Token.LOGICAL_AND {
			return &AST.ExpressionBoolean{Value: lhs.Value && rhs.Value}
		} else {
			return &AST.ExpressionBoolean{Value: lhs.Value || rhs.Value}
		}

	default:
		panic(fmt.Sprintf("unhandled operator: %v", kind))
	}
}

func evaluateIntegers(kind Token.TokenType, lhs, rhs int) AST.Expression {
	switch kind {
	case Token.PLUS:
		return &AST.ExpressionInteger{Value: lhs + rhs}
	case Token.MINUS:
		return &AST.ExpressionInteger{Value: lhs - rhs}
	case Token.STAR:
		return &AST.ExpressionInteger{Value: lhs * rhs}
	case Token.DIVISION:
		return &AST.ExpressionInteger{Value: lhs / rhs}
	case Token.EQUALS_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs == rhs}
	case Token.LESS_THAN:
		return &AST.ExpressionBoolean{Value: lhs < rhs}
	case Token.LESS_THAN_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs <= rhs}
	case Token.GREATER_THAN:
		return &AST.ExpressionBoolean{Value: lhs > rhs}
	case Token.GREATER_THAN_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs >= rhs}
	}

	panic("unreachable")
}

func evaluateFloats(kind Token.TokenType, lhs, rhs float32) AST.Expression {
	switch kind {
	case Token.PLUS:
		return &AST.ExpressionFloat{Value: lhs + rhs}
	case Token.MINUS:
		return &AST.ExpressionFloat{Value: lhs - rhs}
	case Token.STAR:
		return &AST.ExpressionFloat{Value: lhs * rhs}
	case Token.DIVISION:
		return &AST.ExpressionFloat{Value: lhs / rhs}
	case Token.LESS_THAN:
		return &AST.ExpressionBoolean{Value: lhs < rhs}
	case Token.LESS_THAN_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs <= rhs}
	case Token.GREATER_THAN:
		return &AST.ExpressionBoolean{Value: lhs > rhs}
	case Token.GREATER_THAN_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs >= rhs}
	case Token.EQUALS_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs >= rhs}
	case Token.NOT_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs >= rhs}
	}
	panic("unreachable")
}

func evaluateString(kind Token.TokenType, lhs, rhs string) AST.Expression {
	switch kind {
	case Token.PLUS:
		return &AST.ExpressionString{Value: lhs + rhs}
	}

	panic("unreachable")
}

func interpretExpression(e AST.Expression, scope *Scope) AST.Expression {
	if e == nil {
		return nil
	}

	switch v := e.(type) {
	case *AST.ExpressionInteger, *AST.ExpressionFloat, *AST.ExpressionBoolean, *AST.ExpressionString:
		return v
	case *AST.ExpressionIdentifier:
		return scope.get(v.Tok)
	case *AST.ExpressionFunctionCall:
		functionDeclaration := globalFunctions[v.Tok.Lexeme]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		functionScope := CreateScope(&globalScope)
		for i := 0; i < argCount; i++ {
			param := functionDeclaration.Parameters[i]
			arg := v.Arguments[i]
			functionScope.set(param.Tok, interpretExpression(arg, scope))
		}

		return interpretExpression(interpretNodes(functionDeclaration.Block.Body, &functionScope), &functionScope)

	case *AST.ExpressionLen:
		v.Array = interpretExpression(v.Array, scope)
		return &AST.ExpressionInteger{Value: len(v.Array.(*AST.ExpressionArray).Elements)}
	case *AST.ExpressionGrouping:
		return interpretExpression(v.Expr, scope)
	case *AST.ExpressionBinary:
		leftExpression := interpretExpression(v.Left, scope)
		rightExpression := interpretExpression(v.Right, scope)
		return interpretBinaryExpression(v.Operator.Kind, leftExpression, rightExpression)
	case *AST.ExpressionArray:
		for i, element := range v.Elements {
			v.Elements[i] = interpretExpression(element, scope)
		}

		return v

	case *AST.ExpressionArrayAccess:
		var ret AST.Expression
		arr := scope.get(v.Tok).(*AST.ExpressionArray)
		for i := 0; i < len(v.Indices); i++ {
			index := interpretExpression(v.Indices[i], scope).(*AST.ExpressionInteger).Value

			if i < len(v.Indices)-1 {
				arr = interpretExpression(arr.Elements[index], scope).(*AST.ExpressionArray)
			} else {
				ret = interpretExpression(arr.Elements[index], scope)
			}
		}

		return ret

	case *AST.ExpressionPseudo:
		return interpretExpression(v.Expr, scope)

	default:
		fmt.Printf("Type: %T\n", e)
		panic("unreachable")
	}
}

func interpretDeclaration(decl AST.Declaration, scope *Scope) {
	switch v := decl.(type) {
	case *AST.DeclarationVariable:
		if scope.has(v.Tok) {
			panic("Attempting to redeclare: " + v.Tok.Lexeme)
		}

		v.RHS = interpretExpression(v.RHS, scope)
		if v.RHS == nil {
			panic("Attempting to assign void to variable: " + v.Tok.Lexeme)
		}

		scope.set(v.Tok, v.RHS)

	case *AST.DeclarationFunction:
		globalFunctions[v.Tok.Lexeme] = v

	}
}

func printExpression(expr AST.Expression, scope *Scope) {
	switch v := expr.(type) {
	case *AST.ExpressionInteger:
		fmt.Print(v.Value)

	case *AST.ExpressionFloat:
		fmt.Print(v.Value)

	case *AST.ExpressionBoolean:
		fmt.Print(v.Value)

	case *AST.ExpressionString:
		fmt.Print(v.Value)

	case *AST.ExpressionIdentifier:
		printExpression(scope.get(v.Tok), scope)

	case *AST.ExpressionArray:
		fmt.Print("[")
		for i, elem := range v.Elements {
			if i > 0 {
				fmt.Print(" ")
			}
			printExpression(interpretExpression(elem, scope), scope)
		}
		fmt.Print("]")

	default:
		panic(fmt.Sprintf("unprintable type: %T", v))
	}
}

func interpretStatement(s AST.Statement, scope *Scope) AST.Expression {
	switch v := s.(type) {
	case *AST.StatementPrint:
		printExpression(interpretExpression(v.Expr, scope), scope)
		fmt.Println("")

		return nil

	case *AST.StatementAssignment:
		if !scope.has(v.Tok) {
			panic(fmt.Sprintf("Line %d | Attempting to assign to undeclared identifier: %s", v.Tok.Line, v.Tok.Lexeme))
		}

		temp := interpretExpression(v.RHS, scope)
		if temp == nil {
			panic(fmt.Sprintf("Line %d | Attempting to assign void to variable: %s", v.Tok.Line, v.Tok.Lexeme))
		}

		scope.set(v.Tok, temp)

		return nil

	case *AST.StatementBlock:
		blockScope := CreateScope(scope)
		return interpretNodes(v.Body, &blockScope)

	case *AST.StatementFor:
		forScope := CreateScope(scope)
		interpretDeclaration(v.Initializer, &forScope)
		for interpretExpression(v.Condition, &forScope).(*AST.ExpressionBoolean).Value {
			blockRet := interpretStatement(v.Block, &forScope)
			if pseudo, ok := blockRet.(*AST.ExpressionPseudo); ok {
				if pseudo.Behavior == AST.BREAK {
					break
				} else if pseudo.Behavior == AST.RETURN {
					return pseudo.Expr
				} else if pseudo.Behavior == AST.CONTINUE {
				} else {
					panic("unreachable")
				}
			}

			interpretStatement(v.Increment, &forScope)
		}

		return nil

	case *AST.StatementReturn:
		return &AST.ExpressionPseudo{
			Expr:     interpretExpression(v.Expr, scope),
			Behavior: AST.RETURN,
		}

	case *AST.StatementBreak:
		return &AST.ExpressionPseudo{
			Expr:     nil,
			Behavior: AST.BREAK,
		}

	case *AST.StatementContinue:
		return &AST.ExpressionPseudo{
			Expr:     nil,
			Behavior: AST.CONTINUE,
		}

	case *AST.StatementIfElse:
		cond := interpretExpression(v.Condition, scope).(*AST.ExpressionBoolean)
		if cond.Value {
			return interpretStatement(v.IfBlock, scope)
		} else {
			if v.ElseBlock != nil {
				return interpretStatement(v.ElseBlock, scope)
			}
		}

	default:
		fmt.Printf("Type: %T\n", v)
		panic("unreachable")
	}

	return nil
}

func interpretNodes(nodes []AST.Node, scope *Scope) AST.Expression {
	for _, node := range nodes {
		switch v := node.(type) {
		case AST.Statement:
			ret := interpretStatement(v, scope)
			if pseudo, ok := ret.(*AST.ExpressionPseudo); ok {
				return pseudo
			}

		case AST.Declaration:
			interpretDeclaration(v, scope)
		}
	}

	return nil
}

func InterpretProgram(program AST.Program) {
	globalScope = CreateScope(nil)
	globalFunctions = make(map[string]*AST.DeclarationFunction)

	for _, decl := range program.Declarations {
		interpretDeclaration(decl, &globalScope)
	}

	if mainDecl, ok := globalFunctions["main"]; ok {
		mainCall := &AST.ExpressionFunctionCall{
			Tok: mainDecl.Tok,
		}

		interpretExpression(mainCall, &globalScope)
	} else {
		panic("main function not found")
	}
}
