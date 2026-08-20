package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluatePostfix(node *ast.Postfix, scope *object.Scope) object.Object {
	switch node.Operator {
	case token.PLUSPLUS:
		value, ok := scope.Environment.Get(node.Token.Lexeme)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		if value.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		newValue := value.(*object.Number).Increment()

		scope.Environment.Set(node.Token.Lexeme, newValue)

		return newValue
	case token.MINUSMINUS:
		value, ok := scope.Environment.Get(node.Token.Lexeme)

		if !ok {
			return newError("%d:%d:%s: runtime error: identifier not found: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		if value.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: identifier is not a number: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Token.Lexeme)
		}

		newValue := value.(*object.Number).Decrement()

		scope.Environment.Set(node.Token.Lexeme, newValue)

		return newValue
	default:
		return newError("%d:%d:%s: runtime error: unknown operator: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Operator)
	}
}
