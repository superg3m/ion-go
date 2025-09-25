package TypeChecker

import (
	"fmt"
	"ion-go/AST"
	"ion-go/TS"
)

var globalFunctions map[string]*AST.DeclarationFunction

func typeCheckExpression(e AST.Expression, env *TypeEnv) *TS.Type {
	switch v := e.(type) {
	case *AST.ExpressionInteger:
		return TS.NewType(TS.INTEGER, nil, nil)

	case *AST.ExpressionFloat:
		return TS.NewType(TS.FLOAT, nil, nil)

	case *AST.ExpressionBoolean:
		return TS.NewType(TS.BOOL, nil, nil)

	case *AST.ExpressionString:
		return TS.NewType(TS.STRING, nil, nil)

	case *AST.ExpressionIdentifier:
		decl := env.get(v.Tok)
		return decl.DeclType

	case *AST.ExpressionBinary:
		lt := typeCheckExpression(v.Left, env)
		rt := typeCheckExpression(v.Right, env)

		promotedType := TS.GetPromotedType(v.Operator, lt, rt)
		if promotedType == TS.INVALID_TYPE {
			panic(fmt.Sprintf("Typechecking error Line %d | Operation %s not supported on Left: %s | Right: %s", v.Operator.Line, v.Operator.Lexeme, lt.String(), rt.String()))
		}

		return TS.NewType(promotedType, nil, nil)

	case *AST.ExpressionFunctionCall:
		decl, ok := globalFunctions[v.Tok.Lexeme]
		if !ok {
			panic("undefined function " + v.Tok.Lexeme)
		}

		functionDeclaration := globalFunctions[v.Tok.Lexeme]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.DeclType.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		for i := 0; i < argCount; i++ {
			param := functionDeclaration.DeclType.Parameters[i]
			argType := typeCheckExpression(v.Arguments[i], env)

			if !TS.TypeCompare(param.DeclType, argType) {
				panic(fmt.Sprintf("Line %d | argument %d: expected %s, got %s", v.Tok.Line, i, argType.String(), param.DeclType.String()))
			}
		}

		return decl.DeclType.GetReturnType()

	case *AST.ExpressionArray:
		for i, element := range v.Elements {
			if ref, ok := element.(*AST.ExpressionArray); ok {
				ref.DeclType = v.DeclType.RemoveArrayModifier()
			}

			elementType := typeCheckExpression(element, env)
			if !TS.TypeCompare(elementType, v.DeclType.RemoveArrayModifier()) {
				panic(fmt.Sprintf("Element %d: expected %s, got %s", i, v.DeclType.RemoveArrayModifier().String(), elementType.String()))
			}
		}

		return v.DeclType

	case *AST.ExpressionArrayAccess:
		decl := env.get(v.Tok)
		accessType := decl.DeclType
		for i := 0; i < len(v.Indices); i++ {
			if !accessType.IsArray() {
				panic(fmt.Sprintf("Line: %d | undefined array access: %s", v.Tok.Line, v.Tok.Lexeme))
			}

			accessType = accessType.RemoveArrayModifier()
		}

		return accessType

	case *AST.ExpressionLen:
		if _, ok := v.Array.(*AST.ExpressionArray); ok {
			panic("Builtin Len() argument is not iterable")
		}

		return TS.NewType(TS.INTEGER, nil, nil)

	case *AST.ExpressionUnary:
		return typeCheckExpression(v.Operand, env)

	case *AST.ExpressionGrouping:
		return typeCheckExpression(v.Expr, env)

	default:
		panic(fmt.Sprintf("undefined statement: %T", v))
	}

	return TS.NewType(TS.INVALID_TYPE, nil, nil)
}

func typeCheckStatement(s AST.Statement, env *TypeEnv) {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		decl := env.get(v.Tok)
		rhsType := typeCheckExpression(v.RHS, env)

		if !TS.TypeCompare(decl.DeclType, rhsType) {
			panic(fmt.Sprintf("Line %d | Can't assign type %s to type %s", v.Tok.Line, rhsType.String(), decl.DeclType.String()))
		}

	case *AST.StatementPrint:
		typeCheckExpression(v.Expr, env)

	case *AST.StatementReturn:
		if v.Expr != nil {
			typeCheckExpression(v.Expr, env)
		}

	case *AST.StatementBreak, *AST.StatementContinue:
		if env.CurrentStatus != IN_LOOP {
			panic("break statement is not in loop")
		}

	case *AST.StatementFor:
		typeCheckDeclaration(v.Initializer, env)
		condition := typeCheckExpression(v.Condition, env)
		if condition.Kind != TS.BOOL {
			panic("For statement condition doesn't resolve to a bool it resolves to: " + condition.String())
		}

		typeCheckStatement(v.Increment, env)

		env.CurrentStatus = IN_LOOP
		for _, node := range v.Block.Body {
			typeCheckNode(node, env)
		}
		env.CurrentStatus = NORMAL

	case *AST.StatementWhile:
		condition := typeCheckExpression(v.Condition, env)
		if condition.Kind != TS.BOOL {
			panic("For statement condition doesn't resolve to a bool it resolves to: " + condition.String())
		}

		env.CurrentStatus = IN_LOOP
		for _, node := range v.Block.Body {
			typeCheckNode(node, env)
		}
		env.CurrentStatus = NORMAL

	case *AST.StatementIfElse:
		condition := typeCheckExpression(v.Condition, env)
		if condition.Kind != TS.BOOL {
			panic("For statement condition doesn't resolve to a bool it resolves to: " + condition.String())
		}

		for _, node := range v.IfBlock.Body {
			typeCheckNode(node, env)
		}

		if v.ElseBlock != nil {
			for _, node := range v.ElseBlock.Body {
				typeCheckNode(node, env)
			}
		}

	case *AST.StatementDefer:
		typeCheckNode(v.DeferredNode.(AST.Node), env)

	default:
		panic(fmt.Sprintf("undefined statement: %T", v))

	}
}

func typeCheckDeclaration(decl AST.Declaration, env *TypeEnv) {
	switch v := decl.(type) {
	case *AST.DeclarationVariable:
		rhsType := typeCheckExpression(v.RHS, env)
		if v.DeclType == nil || v.DeclType.Kind == TS.INVALID_TYPE {
			v.DeclType = rhsType
		}

		env.set(v.Tok, v)

		if !TS.TypeCompare(v.DeclType, rhsType) {
			panic(fmt.Sprintf("Line: %d |  Can't assign type %s to type %s", v.Tok.Line, rhsType.String(), v.DeclType.String()))
		}

	case *AST.DeclarationFunction:
		if _, ok := globalFunctions[v.Tok.Lexeme]; ok {
			panic("Attempting to redeclare function " + v.Tok.Lexeme)
		} else {
			globalFunctions[v.Tok.Lexeme] = v
		}

		_, ok := v.Block.Body[len(v.Block.Body)-1].(*AST.StatementReturn)
		if !ok && v.DeclType.GetReturnType().Kind != TS.VOID {
			panic(fmt.Sprintf("Missing return type in %s() body", v.Tok.Lexeme))
		}

		funcEnv := NewTypeEnv(env)
		for _, param := range v.DeclType.Parameters {
			funcEnv.set(param.Tok, &AST.DeclarationVariable{
				Tok:      param.Tok,
				DeclType: param.DeclType,
			})
		}

		for _, node := range v.Block.Body {
			if ret, ok := node.(*AST.StatementReturn); ok {
				if v.DeclType.GetReturnType().Kind == TS.VOID && ret.Expr != nil {
					panic(fmt.Sprintf("Attempting to return expression in %s() with return type void", v.Tok.Lexeme))
				}

				retType := typeCheckExpression(ret.Expr, funcEnv)
				if TS.TypeCompare(v.DeclType, retType) {
					panic(fmt.Sprintf("%s() has a return type of %s but returns a %s", v.Tok.Lexeme, v.DeclType.GetReturnType().String(), retType.String()))
				}

				continue
			}

			typeCheckNode(node, funcEnv)
		}
	}
}

func typeCheckNode(node AST.Node, env *TypeEnv) {
	switch v := node.(type) {
	case AST.Statement:
		typeCheckStatement(v, env)
	case AST.Expression:
		typeCheckExpression(v, env)
	case AST.Declaration:
		typeCheckDeclaration(v, env)
	}
}

func TypeCheckProgram(program AST.Program) {
	globalEnv := NewTypeEnv(nil)
	globalFunctions = make(map[string]*AST.DeclarationFunction)

	for _, decl := range program.Declarations {
		typeCheckDeclaration(decl, globalEnv)
	}
}
