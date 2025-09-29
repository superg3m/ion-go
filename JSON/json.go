package JSON

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ion-go/AST"
)

func expressionToJson(e AST.Expression) any {
	if e == nil {
		return nil
	}

	switch v := e.(type) {
	case *AST.ExpressionInteger:
		return v.Value
	case *AST.ExpressionFloat:
		return v.Value
	case *AST.ExpressionBoolean:
		return v.Value
	case *AST.ExpressionString:
		return v.Value

	case *AST.ExpressionIdentifier:
		return map[string]any{
			"Identifier": v.Tok.Lexeme,
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
		return map[string]any{
			"ArrayAccess": v.Tok.Lexeme,
		}

	case *AST.ExpressionLen:
		return map[string]any{
			"ExpressionLen": expressionToJson(v.Iterable),
		}

	case *AST.ExpressionUnary:
		return map[string]any{
			"ExpressionUnary": expressionToJson(v.Operand),
		}

	case *AST.ExpressionGrouping:
		return map[string]any{
			"ExpressionGrouping": expressionToJson(v.Expr),
		}

	case *AST.SE_FunctionCall:
		return map[string]any{
			"FunctionCall": nil,
		}

	case *AST.ExpressionStructMemberAccess:
		return map[string]any{
			"ExpressionStructMemberAccess": nil,
		}

	default:
		panic(fmt.Sprintf("%T", v))
	}

	return nil
}

func statementToJson(s AST.Statement) any {
	switch v := s.(type) {
	case *AST.StatementAssignment:
		return map[string]any{
			"AssignmentStatement": map[string]any{
				"lhs": expressionToJson(v.LHS),
				"rhs": expressionToJson(v.RHS),
			},
		}

	case *AST.StatementReturn:
		return map[string]any{
			"ReturnStatement": expressionToJson(v.Expr),
		}

	case *AST.StatementDefer:
		return map[string]any{
			"DeferStatement": nodeToJson(v.DeferredNode),
		}

	case *AST.StatementContinue:
		return "ContinueStatement"

	case *AST.StatementBreak:
		return "BreakStatement"

	case *AST.StatementPrint:
		return map[string]any{
			"PrintStatement": expressionToJson(v.Expr),
		}

	case *AST.StatementBlock:
		var body []any
		for _, node := range v.Body {
			body = append(body, nodeToJson(node))
		}

		return map[string]any{
			"StatementBlock": body,
		}

	case *AST.StatementFor:
		return map[string]any{
			"ForStatement": map[string]any{
				"Initializer": declarationToJson(v.Initializer),
				"Condition":   expressionToJson(v.Condition),
				"Increment":   statementToJson(v.Increment),
				"Block":       statementToJson(v.Block),
			},
		}

	case *AST.StatementWhile:
		return map[string]any{
			"WhileStatement": map[string]any{
				"Condition": expressionToJson(v.Condition),
				"Block":     statementToJson(v.Block),
			},
		}

	case *AST.StatementIfElse:
		var elseBlockJson any = nil
		if v.ElseBlock != nil {
			elseBlockJson = statementToJson(v.ElseBlock)
		}

		return map[string]any{
			"IfElseStatement": map[string]any{
				"Condition": expressionToJson(v.Condition),
				"IfBlock":   statementToJson(v.IfBlock),
				"ElseBlock": elseBlockJson,
			},
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
			"Name":     v.Tok.Lexeme,
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
			"Name":     v.Tok.Lexeme,
			"DeclType": v.DeclType.String(),
			"Body":     body,
		}
		return map[string]any{
			"FunctionDeclaration": desc,
		}

	case *AST.DeclarationStruct:
		desc := map[string]any{
			"Name": v.Tok.Lexeme,
		}
		return map[string]any{
			"StructDeclaration": desc,
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

func MarshalIndentNoEscape(v any, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// Encode always appends a newline; strip it if you don’t want it
	out := buf.Bytes()
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out, nil
}

func PrettyPrint(program AST.Program) {
	var declarations []any
	for _, decl := range program.Declarations {
		declarations = append(declarations, declarationToJson(decl))
	}

	indent, err := MarshalIndentNoEscape(map[string]any{"Declarations": declarations}, "", "    ")
	if err != nil {
		return
	}

	fmt.Println(string(indent))
}
