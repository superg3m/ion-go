package JSON

import (
	"encoding/json"
	"fmt"
	"ion-go/AST"
)

func expressionToJson(e AST.Expression) any {
	switch v := e.(type) {
	case *AST.ExpressionInteger, *AST.ExpressionFloat, *AST.ExpressionBoolean:
		return v

	case *AST.ExpressionIdentifier:
		return map[string]any{
			"Identifier": v.Name,
		}

	default:
		panic(fmt.Sprintf("%T", v))
	}

	return nil
}

func statementToJson(s AST.Statement) map[string]any {
	switch v := s.(type) {
	case *AST.StatementReturn:
		return map[string]any{
			"ReturnStatement": expressionToJson(v.Expr),
		}
	case *AST.StatementPrint:
		return map[string]any{
			"PrintStatement": expressionToJson(v.Expr),
		}

	default:
		panic(fmt.Sprintf("%T", v))
	}
	return nil
}

func declarationToJson(decl AST.Declaration) map[string]any {
	switch v := decl.(type) {
	case *AST.DeclarationVariable:
		desc := map[string]any{
			"Name":     v.Name,
			"DeclType": v.DeclType.Name,
		}
		return map[string]any{
			"VariableDeclaration": desc,
		}

	case *AST.DeclarationFunction:
		var body []any
		for _, node := range v.Block.Body {
			body = append(body, nodeToJson(node))
		}
		desc := map[string]any{
			"Name":     v.Name,
			"DeclType": v.ReturnType.Name,
			"Body":     body,
		}
		return map[string]any{
			"FunctionDeclaration": desc,
		}

	default:
		panic(fmt.Sprintf("%T", v))
	}
	return nil
}

func nodeToJson(node AST.Node) any {
	switch v := node.(type) {
	case AST.Expression:
		return expressionToJson(v)

	case AST.Statement:
		return statementToJson(v)

	case AST.Declaration:
		return declarationToJson(v)

	default:
		panic(fmt.Sprintf("%T", v))
	}
	return nil
}

func PrettyPrint(program AST.Program) {
	var declarations []any
	for _, decl := range program.Declarations {
		declarations = append(declarations, declarationToJson(decl))
	}

	indent, err := json.MarshalIndent(map[string]any{"Declarations": declarations}, "", "    ")
	if err != nil {
		return
	}

	fmt.Println(string(indent))
}
