package TypeChecker

import (
	"fmt"
	"ion-go/AST"
	"ion-go/TS"
)

type StatementTypePair struct {
	stmt *AST.StatementReturn
	t    *TS.Type
}

var globalFunctions map[string]*AST.DeclarationFunction
var globalStruct map[string]*AST.DeclarationStruct
var globalReturnStatementStack []StatementTypePair

func typeCheckFunctionCall(v *AST.SE_FunctionCall, env *TypeEnv) *TS.Type {
	functionDeclaration, ok := globalFunctions[v.Tok.Lexeme]
	if !ok {
		panic("undefined function " + v.Tok.Lexeme)
	}

	argCount := len(v.Arguments)
	paramCount := len(functionDeclaration.DeclType.Parameters)

	if paramCount != argCount {
		panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
	}

	for i := 0; i < argCount; i++ {
		param := functionDeclaration.DeclType.Parameters[i]
		argType := typeCheckExpression(v.Arguments[i], env)

		if !TS.TypeCompare(param.DeclType, argType) {
			panic(fmt.Sprintf("Line %d | argument %d: expected %s, got %s", v.Tok.Line, i, param.DeclType.String(), argType.String()))
		}
	}

	return functionDeclaration.DeclType.GetReturnType()
}

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

	case *AST.SE_FunctionCall:
		return typeCheckFunctionCall(v, env)

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

	case *AST.ExpressionLen:
		switch ev := v.Iterable.(type) {
		case *AST.ExpressionArray:
		case *AST.ExpressionString:
		case *AST.ExpressionIdentifier:
			evType := env.get(ev.Tok)
			if evType.DeclType.Kind != TS.ARRAY && evType.DeclType.Kind != TS.STRING {
				panic(fmt.Sprintf("Builtin Len() argument is not iterable"))
			}
		default:
			panic(fmt.Sprintf("Builtin Len() argument is not iterable %T", v.Iterable))
		}

		return TS.NewType(TS.INTEGER, nil, nil)

	case *AST.ExpressionUnary:
		return typeCheckExpression(v.Operand, env)

	case *AST.ExpressionGrouping:
		return typeCheckExpression(v.Expr, env)

	case *AST.ExpressionTypeCast:
		exprType := typeCheckExpression(v.Expr, env)
		if TS.TypeCompare(v.CastType, exprType) {
			return v.CastType
		}

		if !TS.CanCastType(v.CastType, exprType) {
			panic(fmt.Sprintf("Typechecking error Line %d | Invalid cast to %s from %s", v.Tok.Line, v.CastType.String(), exprType.String()))
		}

		return v.CastType

	case *AST.ExpressionStruct:
		structDecl, ok := globalStruct[v.Tok.Lexeme]
		if !ok {
			panic("Undefined type: " + v.Tok.Lexeme)
		}

		argCount := len(v.MemberValues)
		memberCount := len(structDecl.Members)

		if memberCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, memberCount))
		}

		for i := 0; i < argCount; i++ {
			member := structDecl.Members[i]
			argType := typeCheckExpression(v.MemberValues[member.Tok.Lexeme], env)

			if !TS.TypeCompare(member.DeclType, argType) {
				panic(fmt.Sprintf("Line %d | argument %d: expected %s: %s, got %s", v.Tok.Line, i, member.Tok.Lexeme, member.DeclType.String(), argType.String()))
			}
		}

		return TS.NewType(TS.STRUCT, TS.NewType(TS.TypeKind(structDecl.Tok.Lexeme), nil, nil), nil)

	case *AST.ExpressionAccessChain:
		ident := env.get(v.Tok)
		decl := globalStruct[ident.DeclType.String()]

		accessType := ident.DeclType
		accessString := ident.Tok.Lexeme

		for i := 0; i < len(v.AccessKeys); i++ {
			switch ev := v.AccessKeys[i].(type) {
			case *AST.ExpressionIdentifier:
				memberName := ev.Tok
				accessString += "." + memberName.Lexeme
				if accessType == nil || !accessType.IsStruct() {
					panic(fmt.Sprintf("Line: %d | undefined struct access: %s", v.Tok.Line, accessString))
				}

				accessType = decl.MemberLookup[memberName.Lexeme].DeclType
				decl = globalStruct[accessType.String()]

			case *AST.ExpressionArrayAccess:
				index, ok := ev.Index.(*AST.ExpressionInteger)
				if ok {
					accessString += fmt.Sprintf("[%d]", index.Value)
				}

				identifier, ok := ev.Index.(*AST.ExpressionIdentifier)
				if ok {
					if !TS.TypeCompare(typeCheckExpression(identifier, env), TS.NewType(TS.INTEGER, nil, nil)) {
						panic("Array Index Access is not of type int")
					}

					accessString += fmt.Sprintf("[%d]", identifier.Tok.Lexeme)
				}

				acc, ok := ev.Index.(*AST.ExpressionAccessChain)
				if ok {
					if !TS.TypeCompare(typeCheckExpression(acc, env), TS.NewType(TS.INTEGER, nil, nil)) {
						panic("Array Index Access is not of type int")
					}

					accessString += fmt.Sprintf("[...]")
				}

				if accessType.IsArray() {
					accessType = accessType.RemoveArrayModifier()
					decl = globalStruct[accessType.String()]
				} else {
					panic(fmt.Sprintf("Line: %d | undefined array access: %s", v.Tok.Line, accessString))
				}
			}
		}

		return accessType

	default:
		panic(fmt.Sprintf("undefined statement: %T", v))
	}

	return TS.NewType(TS.INVALID_TYPE, nil, nil)
}

func typeCheckStatement(s AST.Statement, env *TypeEnv) {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		lhsType := typeCheckExpression(v.LHS, env)
		rhsType := typeCheckExpression(v.RHS, env)

		if !TS.TypeCompare(lhsType, rhsType) {
			panic(fmt.Sprintf("Line %d | Can't assign type %s to type %s", v.Tok.Line, rhsType.String(), lhsType.String()))
		}

	case *AST.StatementPrint:
		typeCheckExpression(v.Expr, env)

	case *AST.StatementReturn:
		if v.Expr != nil {
			globalReturnStatementStack = append(globalReturnStatementStack,
				StatementTypePair{
					stmt: v,
					t:    typeCheckExpression(v.Expr, env),
				},
			)
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
		typeCheckStatement(v.Block, env)
		env.CurrentStatus = NORMAL

	case *AST.StatementIfElse:
		condition := typeCheckExpression(v.Condition, env)
		if condition.Kind != TS.BOOL {
			panic("For statement condition doesn't resolve to a bool it resolves to: " + condition.String())
		}

		typeCheckStatement(v.IfBlock, env)

		if v.ElseBlock != nil {
			typeCheckStatement(v.ElseBlock, env)
		}

	case *AST.StatementDefer:
		typeCheckNode(v.DeferredNode.(AST.Node), env)

	case *AST.StatementBlock:
		for _, node := range v.Body {
			typeCheckNode(node, env)
		}

	case *AST.SE_FunctionCall:
		typeCheckFunctionCall(v, env)

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
			panic(fmt.Sprintf("Line: %d | Can't assign type %s to type %s", v.Tok.Line, rhsType.String(), v.DeclType.String()))
		}

	case *AST.DeclarationFunction:
		if _, ok := globalFunctions[v.Tok.Lexeme]; ok {
			panic("Attempting to redeclare function " + v.Tok.Lexeme)
		} else {
			globalFunctions[v.Tok.Lexeme] = v
		}

		_, ok := v.Block.Body[len(v.Block.Body)-1].(*AST.StatementReturn)
		if !ok && v.DeclType.GetReturnType().Kind != TS.VOID {
			panic(fmt.Sprintf("%s() body is missing a return statement or it is not the last statement in the body", v.Tok.Lexeme))
		}

		funcEnv := NewTypeEnv(env)
		for _, param := range v.DeclType.Parameters {
			funcEnv.set(param.Tok, &AST.DeclarationVariable{
				Tok:      param.Tok,
				DeclType: param.DeclType,
			})
		}

		for _, node := range v.Block.Body {
			typeCheckNode(node, funcEnv)
			for _, pair := range globalReturnStatementStack {
				if v.DeclType.GetReturnType().Kind == TS.VOID {
					panic(fmt.Sprintf("Attempting to return expression in %s() with return type void", v.Tok.Lexeme))
				}

				if !TS.TypeCompare(v.DeclType.GetReturnType(), pair.t) {
					panic(fmt.Sprintf("Line %d | %s() has a return type of %s but returns a %s", pair.stmt.Tok.Line, v.Tok.Lexeme, v.DeclType.GetReturnType().String(), pair.t.String()))
				}
			}
			globalReturnStatementStack = nil
		}

	case *AST.DeclarationStruct:
		if _, ok := globalStruct[v.Tok.Lexeme]; ok {
			panic("Attempting to redeclare type: " + v.Tok.Lexeme)
		} else {
			globalStruct[v.Tok.Lexeme] = v
		}

	default:
		panic(fmt.Sprintf("undefined declaration: %T", v))
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
	globalStruct = make(map[string]*AST.DeclarationStruct)

	for _, decl := range program.Declarations {
		typeCheckDeclaration(decl, globalEnv)
	}
}
