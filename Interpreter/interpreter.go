package Interpreter

import (
	"fmt"
	"ion-go/AST"
	"ion-go/TS"
	"ion-go/Token"
	"strings"
)

var globalFunctions map[string]*AST.DeclarationFunction
var globalStructs map[string]*AST.DeclarationStruct
var globalScope Scope

// Returns either a struct or array and then their respective indices
func evaluateAccessChainExpression(chain *AST.ExpressionAccessChain, scope *Scope) (AST.Expression, AST.Expression) {
	ret := scope.get(chain.Tok)
	for i := 0; i < len(chain.AccessKeys)-1; i++ {
		switch ev := chain.AccessKeys[i].(type) {
		case *AST.ExpressionArrayAccess:
			temp := ret.(*AST.ExpressionArray)
			index := interpretExpression(ev.Index, scope).(*AST.ExpressionInteger).Value
			ret = interpretExpression(temp.Elements[index], scope)

		case *AST.ExpressionIdentifier:
			temp := ret.(*AST.ExpressionStruct)
			index := ev.Tok.Lexeme
			ret = interpretExpression(temp.MemberValues[index], scope)
		}
	}

	var index AST.Expression = nil
	switch ev := chain.AccessKeys[len(chain.AccessKeys)-1].(type) {
	case *AST.ExpressionArrayAccess:
		index = ev.Index

	case *AST.ExpressionIdentifier:
		index = ev
	}

	return ret, index
}

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
				return evaluateString(kind, lhs.Value, fmt.Sprintf("%.5g", rhs.Value))
			case *AST.ExpressionString:
				return evaluateString(kind, lhs.Value, rhs.Value)
			}
		}

		if rhs, ok1 := right.(*AST.ExpressionString); ok1 {
			switch lhs := left.(type) {
			case *AST.ExpressionInteger:
				return evaluateString(kind, fmt.Sprintf("%d", lhs.Value), rhs.Value)
			case *AST.ExpressionFloat:
				return evaluateString(kind, fmt.Sprintf("%.5g", lhs.Value), rhs.Value)
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
	case Token.EQUALS_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs == rhs}
	case Token.NOT_EQUALS:
		return &AST.ExpressionBoolean{Value: lhs != rhs}
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

func evaluateString(kind Token.TokenType, lhs, rhs string) AST.Expression {
	switch kind {
	case Token.PLUS:
		return &AST.ExpressionString{Value: lhs + rhs}
	}

	panic("unreachable")
}

func interpretUnaryExpression(kind Token.TokenType, operand AST.Expression) AST.Expression {
	switch kind {
	case Token.MINUS:
		switch v := operand.(type) {
		case *AST.ExpressionInteger:
			return &AST.ExpressionInteger{Value: -v.Value}
		case *AST.ExpressionFloat:
			return &AST.ExpressionFloat{Value: -v.Value}

		default:
			panic(fmt.Sprintf("unhandled operator: %v", kind))
		}
	default:
		panic(fmt.Sprintf("unhandled operator: %v", kind))
	}
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
	case *AST.SE_FunctionCall:
		functionDeclaration := globalFunctions[v.Tok.Lexeme]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.DeclType.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		functionScope := CreateScope(&globalScope)
		for i := 0; i < argCount; i++ {
			param := functionDeclaration.DeclType.Parameters[i]
			arg := v.Arguments[i]
			functionScope.set(param.Tok, interpretExpression(arg, scope))
		}

		return interpretExpression(interpretNodes(functionDeclaration.Block.Body, &functionScope), &functionScope)

	case *AST.ExpressionLen:
		v.Iterable = interpretExpression(v.Iterable, scope)
		switch ve := v.Iterable.(type) {
		case *AST.ExpressionArray:
			return &AST.ExpressionInteger{Value: len(ve.Elements)}

		case *AST.ExpressionString:
			return &AST.ExpressionInteger{Value: len(ve.Value)}

		default:
			panic(fmt.Sprintf("unhandled operator: %T", v.Iterable))
		}

	case *AST.ExpressionGrouping:
		return interpretExpression(v.Expr, scope)

	case *AST.ExpressionBinary:
		leftExpression := interpretExpression(v.Left, scope)
		if v.Operator.Kind == Token.LOGICAL_OR && leftExpression.(*AST.ExpressionBoolean).Value {
			return &AST.ExpressionBoolean{Value: true}
		} else if v.Operator.Kind == Token.LOGICAL_AND && !leftExpression.(*AST.ExpressionBoolean).Value {
			return &AST.ExpressionBoolean{Value: false}
		}

		rightExpression := interpretExpression(v.Right, scope)
		return interpretBinaryExpression(v.Operator.Kind, leftExpression, rightExpression)

	case *AST.ExpressionArray:
		for i, element := range v.Elements {
			v.Elements[i] = interpretExpression(element, scope)
		}

		return v

	case *AST.ExpressionPseudo:
		return interpretExpression(v.Expr, scope)

	case *AST.ExpressionUnary:
		operand := interpretExpression(v.Operand, scope)
		return interpretUnaryExpression(v.Operator.Kind, operand)

	case *AST.ExpressionStruct:
		for i, element := range v.MemberValues {
			v.MemberValues[i] = interpretExpression(element, scope)
		}

		return v

	case *AST.ExpressionAccessChain:
		ret, index := evaluateAccessChainExpression(v, scope)
		switch ev := ret.(type) {
		case *AST.ExpressionArray:
			return ev.Elements[interpretExpression(index, scope).(*AST.ExpressionInteger).Value]
		case *AST.ExpressionStruct:
			return ev.MemberValues[index.(*AST.ExpressionIdentifier).Tok.Lexeme]

		default:
			panic("unreachable")
		}

	case *AST.ExpressionTypeCast:
		v.Expr = interpretExpression(v.Expr, scope)
		switch ev := v.Expr.(type) {
		case *AST.ExpressionInteger:
			if v.CastType.Kind == TS.STRING {
				v.Expr = &AST.ExpressionString{
					Value: fmt.Sprintf("%d", ev.Value),
				}
			}

			if v.CastType.Kind == TS.FLOAT {
				v.Expr = &AST.ExpressionFloat{
					Value: float32(ev.Value),
				}
			}

		case *AST.ExpressionFloat:
			if v.CastType.Kind == TS.STRING {
				v.Expr = &AST.ExpressionString{
					Value: fmt.Sprintf("%.5g", ev.Value),
				}
			}

			if v.CastType.Kind == TS.INTEGER {
				v.Expr = &AST.ExpressionInteger{
					Value: int(ev.Value),
				}
			}

		default:
			panic(fmt.Sprintf("undefined expression %T", v.Expr))
		}

		return v.Expr

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

		temp := interpretExpression(v.RHS, scope)
		if temp == nil {
			panic("Attempting to assign void to variable: " + v.Tok.Lexeme)
		}

		scope.set(v.Tok, temp)

	case *AST.DeclarationFunction:
		globalFunctions[v.Tok.Lexeme] = v

	case *AST.DeclarationStruct:
		globalStructs[v.Tok.Lexeme] = v

	default:
		panic(fmt.Sprintf("unhandled declaration: %T", decl))
	}
}

func fixNewLineCode(s string) string {
	var ret []byte

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && s[i+1] == 'n' {
			ret = append(ret, '\n')
			i += 1
		} else {
			ret = append(ret, s[i])
		}
	}

	return string(ret)
}

func generateIndent(level int) string {
	if level <= 0 {
		return ""
	}

	return strings.Repeat(" ", level*4)
}

func printExpression(expr AST.Expression, scope *Scope, indentLevel int, newLine bool) {
	nl := ""
	indentForMembers := ""
	indentForCloser := generateIndent(indentLevel)

	if newLine {
		nl = "\n"
		indentForMembers = generateIndent(indentLevel + 1)
	} else {
		indentForCloser = ""
	}

	switch v := expr.(type) {
	case *AST.ExpressionInteger:
		fmt.Print(v.Value)

	case *AST.ExpressionFloat:
		fmt.Printf("%.5g", v.Value)

	case *AST.ExpressionBoolean:
		fmt.Print(v.Value)

	case *AST.ExpressionString:
		fmt.Print(fixNewLineCode(v.Value))

	case *AST.ExpressionIdentifier:
		printExpression(scope.get(v.Tok), scope, indentLevel, newLine)

	case *AST.ExpressionArray:
		fmt.Printf("[")

		nextLevel := indentLevel + 1
		for i, elem := range v.Elements {
			printExpression(interpretExpression(elem, scope), scope, nextLevel, false)

			if i < len(v.Elements)-1 {
				fmt.Printf(", ")
			}
		}

		fmt.Printf("%s]", indentForCloser)

	case *AST.ExpressionStruct:
		structDecl := globalStructs[v.Tok.Lexeme]

		fmt.Printf("{")

		nextLevel := indentLevel + 1

		for i := 0; i < len(structDecl.Members); i++ {
			name := structDecl.Members[i]
			value := v.MemberValues[name.Tok.Lexeme]

			fmt.Printf("%s%s", nl, indentForMembers)

			fmt.Printf("%s: %s = ", name.Tok.Lexeme, name.DeclType.String())

			printExpression(interpretExpression(value, scope), scope, nextLevel, false)

			if i < len(structDecl.Members)-1 {
				fmt.Printf(", ")
			}
		}

		fmt.Printf("%s%s}", nl, indentForCloser)

	default:
		panic(fmt.Sprintf("unprintable type: %T", v))
	}
}

func interpretStatement(s AST.Statement, scope *Scope) AST.Expression {
	switch v := s.(type) {
	case *AST.StatementPrint:
		printExpression(interpretExpression(v.Expr, scope), scope, 0, true)
		if v.IsNewLine {
			fmt.Println("")
		}

		return nil

	case *AST.StatementAssignment:
		if !scope.has(v.Tok) {
			panic(fmt.Sprintf("Line %d | Attempting to assign to undeclared identifier: %s", v.Tok.Line, v.Tok.Lexeme))
		}

		rhs := interpretExpression(v.RHS, scope)
		if rhs == nil {
			panic(fmt.Sprintf("Line %d | Attempting to assign void to variable: %s", v.Tok.Line, v.Tok.Lexeme))
		}

		switch ev := v.LHS.(type) {
		case *AST.ExpressionIdentifier:
			scope.set(ev.Tok, rhs)

		case *AST.ExpressionAccessChain:
			ret, index := evaluateAccessChainExpression(ev, scope)
			switch lv := ret.(type) {
			case *AST.ExpressionArray:
				lv.Elements[interpretExpression(index, scope).(*AST.ExpressionInteger).Value] = rhs
			case *AST.ExpressionStruct:
				lv.MemberValues[index.(*AST.ExpressionIdentifier).Tok.Lexeme] = rhs

			default:
				panic("unreachable")
			}

		default:
			panic("unreachable")
		}

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

	case *AST.StatementWhile:
		for interpretExpression(v.Condition, scope).(*AST.ExpressionBoolean).Value {
			blockRet := interpretStatement(v.Block, scope)
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
		}

		return nil

	case *AST.StatementReturn:
		return &AST.ExpressionPseudo{
			Expr:     interpretExpression(v.Expr, scope),
			Behavior: AST.RETURN,
		}

	case *AST.StatementDefer:
		scope.AddDeferStatement(v)
		return nil

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

	case *AST.SE_FunctionCall:
		functionDeclaration := globalFunctions[v.Tok.Lexeme]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.DeclType.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		functionScope := CreateScope(&globalScope)
		for i := 0; i < argCount; i++ {
			param := functionDeclaration.DeclType.Parameters[i]
			arg := v.Arguments[i]
			functionScope.set(param.Tok, interpretExpression(arg, scope))
		}

		return interpretExpression(interpretNodes(functionDeclaration.Block.Body, &functionScope), &functionScope)

	default:
		fmt.Printf("Type: %T\n", v)
		panic("unreachable")
	}

	return nil
}

func interpretNode(node AST.Node, scope *Scope) AST.Expression {
	switch v := node.(type) {
	case AST.Statement:
		ret := interpretStatement(v, scope)
		if pseudo, ok := ret.(*AST.ExpressionPseudo); ok {
			return pseudo
		}

	case AST.Declaration:
		interpretDeclaration(v, scope)
	}

	return nil
}

func interpretNodes(nodes []AST.Node, scope *Scope) AST.Expression {
	defer scope.ResolveDeferStack()
	for _, node := range nodes {
		ret := interpretNode(node, scope)
		if pseudo, ok := ret.(*AST.ExpressionPseudo); ok {
			return pseudo
		}
	}

	return nil
}

func InterpretProgram(program AST.Program) {
	globalScope = CreateScope(nil)
	globalFunctions = make(map[string]*AST.DeclarationFunction)
	globalStructs = make(map[string]*AST.DeclarationStruct)

	for _, decl := range program.Declarations {
		interpretDeclaration(decl, &globalScope)
	}

	if mainDecl, ok := globalFunctions["main"]; ok {
		mainCall := &AST.SE_FunctionCall{
			Tok:       mainDecl.Tok,
			Arguments: nil,
		}

		interpretStatement(mainCall, &globalScope)
	} else {
		panic("main function not found")
	}
}
