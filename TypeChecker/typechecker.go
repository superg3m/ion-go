package TypeChecker

import (
	"fmt"
	"ion-go/AST"
)

var globalFunctions map[string]*AST.DeclarationFunction

func typeCheckExpression(e AST.Expression, env *TypeEnv) AST.DataType {
	switch v := e.(type) {
	case *AST.ExpressionInteger:
		return AST.DataType{Name: "int"}
	case *AST.ExpressionFloat:
		return AST.DataType{Name: "float"}
	case *AST.ExpressionBoolean:
		return AST.DataType{Name: "bool"}
	case *AST.ExpressionIdentifier:
		decl := env.get(v.Name)

		return decl.DeclType
	case *AST.ExpressionFunctionCall:
		decl, ok := globalFunctions[v.Name]
		if !ok {
			panic("undefined function " + v.Name)
		}

		return decl.ReturnType
	}

	panic("undefined expression")
	return AST.DataType{}
}

func typeCheckStatement(s AST.Statement, env *TypeEnv) {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		decl := env.get(v.Name)
		rhsType := typeCheckExpression(v.RHS, env)

		if decl.DeclType != rhsType {
			panic(fmt.Sprintf("Can't assign type %s to type %s", rhsType.Name, decl.DeclType.Name))
		}

	}
}

func typeCheckDeclaration(decl AST.Declaration, env *TypeEnv) {
	switch v := decl.(type) {
	case *AST.DeclarationVariable:
		rhsType := typeCheckExpression(v.RHS, env)
		if v.DeclType.Name == "" {
			v.DeclType = rhsType
		}

		env.set(v.Name, v)

		if v.DeclType != rhsType {
			panic(fmt.Sprintf("Can't assign type %s to type %s", rhsType.Name, v.DeclType.Name))
		}

	case *AST.DeclarationFunction:
		if _, ok := globalFunctions[v.Name]; ok {
			panic("Attempting to redeclare function " + v.Name)
		} else {
			globalFunctions[v.Name] = v
		}

		_, ok := v.Block.Body[len(v.Block.Body)-1].(*AST.StatementReturn)
		if !ok && v.ReturnType.Name != "void" {
			panic(fmt.Sprintf("Missing return type in %s() body", v.Name))
		}

		for _, node := range v.Block.Body {
			if ret, ok := node.(*AST.StatementReturn); ok {
				if v.ReturnType.Name == "void" && ret.Expr != nil {
					panic(fmt.Sprintf("Attempting to return expression in %s() with return type void", v.Name))
				}
			}

			typeCheckNode(node, env)
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
