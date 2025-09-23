package JSON

import (
	"encoding/json"
	"fmt"
	"ion-go/AST"
)

func expressionToJson(e AST.Expression) any {
	switch v := e.(type) {
	case *AST.ExpressionInteger:
		return v.Value
	case *AST.ExpressionFloat:
		return v.Value
	case *AST.ExpressionBoolean:
		return v.Value

	case *AST.ExpressionIdentifier:
		return map[string]any{
			"Identifier": v.Name,
		}

	case *AST.ExpressionBinary:
		return map[string]any{
			"BinaryOp": map[string]any{
				"Op":    v.Operator.Lexeme,
				"Left":  expressionToJson(v.Left),
				"Right": expressionToJson(v.Right),
			},
		}
	case *AST.ExpressionArray:
		return map[string]any{
			"Elements": v.Elements,
			"DeclType": v.DeclType,
		}

	case *AST.ExpressionArrayAccess:
		return map[string]any{}

	default:
		panic(fmt.Sprintf("%T", v))
	}

	return nil
}

func statementToJson(s AST.Statement) map[string]any {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		return map[string]any{
			"AssignmentStatement": map[string]any{
				"name": v.Name,
				"rhs":  expressionToJson(v.RHS),
			},
		}

	case *AST.StatementReturn:
		return map[string]any{
			"ReturnStatement": expressionToJson(v.Expr),
		}

	case *AST.StatementPrint:
		return map[string]any{
			"PrintStatement": expressionToJson(v.Expr),
		}

		// TODO(Jovanni): Actually implement this
	case *AST.StatementFor:
		return map[string]any{}

	case *AST.StatementIfElse:
		return map[string]any{}

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
			"DeclType": v.DeclType.String(),
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
			"DeclType": v.ReturnType.String(),
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
