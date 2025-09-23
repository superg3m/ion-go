package TypeChecker

import (
	"fmt"
	"ion-go/AST"
	"ion-go/Token"
)

var globalFunctions map[string]*AST.DeclarationFunction

func compatibleTypes(lt, rt AST.DataType) (AST.DataType, bool) {
	if lt.String() == rt.String() {
		return lt, true
	} else if lt.String() == "float" && rt.String() == "int" {
		return lt, true
	} else if lt.String() == "int" && rt.String() == "float" {
		return rt, true
	}

	return AST.CreateDataType(""), false
}

func typeCheckExpression(e AST.Expression, env *TypeEnv) AST.DataType {
	switch v := e.(type) {
	case *AST.ExpressionInteger:
		return AST.CreateDataType("int")

	case *AST.ExpressionFloat:
		return AST.CreateDataType("float")

	case *AST.ExpressionBoolean:
		return AST.CreateDataType("bool")

	case *AST.ExpressionIdentifier:
		decl := env.get(v.Name)
		return decl.DeclType

	case *AST.ExpressionBinary:
		lt := typeCheckExpression(v.Left, env)
		rt := typeCheckExpression(v.Right, env)

		dataType, ok := compatibleTypes(lt, rt)
		if !ok {
			panic(fmt.Sprintf("type check failed: lt %v rt %v", lt.String(), rt.String()))
		}

		switch v.Operator.Kind {
		case Token.EQUALS_EQUALS, Token.NOT_EQUALS, Token.LESS_THAN,
			Token.LESS_THAN_EQUALS, Token.GREATER_THAN_EQUALS, Token.GREATER_THAN:
			return AST.CreateDataType("bool")

		case Token.PLUS, Token.MINUS, Token.STAR, Token.DIVISION, Token.MODULUS:
			return dataType

		default:
			panic(fmt.Sprintf("Undefined binary: %s", v.Operator.Kind))

		}

	case *AST.ExpressionFunctionCall:
		decl, ok := globalFunctions[v.Name]
		if !ok {
			panic("undefined function " + v.Name)
		}

		functionDeclaration := globalFunctions[v.Name]
		argCount := len(v.Arguments)
		paramCount := len(functionDeclaration.Parameters)

		if paramCount != argCount {
			panic(fmt.Sprintf("expected %d parameter(s), got %d", argCount, paramCount))
		}

		for i := 0; i < argCount; i++ {
			param := functionDeclaration.Parameters[i]
			argType := typeCheckExpression(v.Arguments[i], env)

			if param.DeclType.String() != argType.String() {
				panic(fmt.Sprintf("argument %d: expected %s, got %s", i, argType.String(), param.DeclType.String()))
			}
		}

		return decl.ReturnType

	case *AST.ExpressionArray:
		firstElementType := AST.CreateDataType("")
		for i, element := range v.Elements {
			elementType := typeCheckExpression(element, env)
			if firstElementType.String() == "" {
				firstElementType = elementType
			}

			if elementType.String() != firstElementType.String() {
				panic(fmt.Sprintf("Element %d: expected %s, got %s", i, firstElementType.String(), elementType.String()))
			}
		}

		v.DeclType = firstElementType
		return AST.CreateDataType(v.DeclType.String() + AST.ARRAY)

	case *AST.ExpressionArrayAccess:
		decl := env.get(v.Name)
		accessType := decl.DeclType.String()[:len(decl.DeclType.String())-2]

		return AST.CreateDataType(accessType)

	case *AST.ExpressionLen:
		if _, ok := v.Array.(*AST.ExpressionArray); ok {
			panic("Builtin Len() argument is not iterable")
		}

		return AST.CreateDataType("int")

	default:
		panic(fmt.Sprintf("undefined statement: %T", v))
	}

	return AST.DataType{}
}

func typeCheckStatement(s AST.Statement, env *TypeEnv) {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		decl := env.get(v.Name)
		rhsType := typeCheckExpression(v.RHS, env)

		if decl.DeclType.String() != rhsType.String() {
			panic(fmt.Sprintf("Can't assign type %s to type %s", rhsType.String(), decl.DeclType.String()))
		}

	case *AST.StatementPrint:
		typeCheckExpression(v.Expr, env)

	case *AST.StatementReturn:
		typeCheckExpression(v.Expr, env)

	case *AST.StatementFor:
		typeCheckDeclaration(v.Initializer, env)
		condition := typeCheckExpression(v.Condition, env)
		if condition.String() != "bool" {
			panic("For statement condition doesn't resolve to a bool it resolves to: " + condition.String())
		}

		typeCheckStatement(v.Increment, env)

		for _, node := range v.Block.Body {
			typeCheckNode(node, env)
		}

	case *AST.StatementIfElse:
		condition := typeCheckExpression(v.Condition, env)
		if condition.String() != "bool" {
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

	default:
		panic(fmt.Sprintf("undefined statement: %T", v))

	}
}

func typeCheckDeclaration(decl AST.Declaration, env *TypeEnv) {
	switch v := decl.(type) {
	case *AST.DeclarationVariable:
		rhsType := typeCheckExpression(v.RHS, env)
		if v.DeclType.String() == "" {
			v.DeclType = rhsType
		}

		env.set(v.Name, v)

		if v.DeclType.String() != rhsType.String() {
			panic(fmt.Sprintf("Can't assign type %s to type %s", rhsType.String(), v.DeclType.String()))
		}

	case *AST.DeclarationFunction:
		if _, ok := globalFunctions[v.Name]; ok {
			panic("Attempting to redeclare function " + v.Name)
		} else {
			globalFunctions[v.Name] = v
		}

		_, ok := v.Block.Body[len(v.Block.Body)-1].(*AST.StatementReturn)
		if !ok && v.ReturnType.String() != "void" {
			panic(fmt.Sprintf("Missing return type in %s() body", v.Name))
		}

		funcEnv := NewTypeEnv(env)
		for _, param := range v.Parameters {
			funcEnv.set(param.Name, &AST.DeclarationVariable{
				Name:     param.Name,
				DeclType: param.DeclType,
			})
		}

		for _, node := range v.Block.Body {
			if ret, ok := node.(*AST.StatementReturn); ok {
				if v.ReturnType.String() == "void" && ret.Expr != nil {
					panic(fmt.Sprintf("Attempting to return expression in %s() with return type void", v.Name))
				}

				retType := typeCheckExpression(ret.Expr, funcEnv)
				if v.ReturnType.String() != retType.String() {
					panic(fmt.Sprintf("%s() has a return type of %s but returns a %s", v.Name, v.ReturnType.String(), retType.String()))
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
