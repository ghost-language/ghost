package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

func evaluatePostfix(node *ast.Postfix, scope *object.Scope) object.Object {
	switch node.Operator {
	case "++":
		value, ok := scope.Environment.Get(node.Token.Lexeme)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		if value.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		one := decimal.NewFromInt(1)

		newValue := &object.Number{
			Value: value.(*object.Number).Value.Add(one),
		}

		scope.Environment.Set(node.Token.Lexeme, newValue)

		return newValue
	case "--":
		value, ok := scope.Environment.Get(node.Token.Lexeme)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		if value.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		one := decimal.NewFromInt(1)

		newValue := &object.Number{
			Value: value.(*object.Number).Value.Sub(one),
		}

		scope.Environment.Set(node.Token.Lexeme, newValue)

		return newValue
	default:
		return newError("%d:%d:%s: runtime error: unknown operator: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Operator)
	}
}
