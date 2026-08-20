package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrefix(node *ast.Prefix, scope *object.Scope) object.Object {
	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case token.BANG:
		switch right {
		case value.TRUE:
			return value.FALSE
		case value.FALSE:
			return value.TRUE
		case value.NULL:
			return value.TRUE
		default:
			return value.FALSE
		}
	case token.MINUS:
		// Only works with number objects
		if right.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: unknown operator: -%s", node.Token.Line, node.Token.Column, node.Token.File, right.Type())
		}

		return right.(*object.Number).Neg()
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s%s", node.Token.Line, node.Token.Column, node.Token.File, node.Operator, right.Type())
}
