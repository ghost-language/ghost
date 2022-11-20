package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateInfix(node *ast.Infix, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evaluateBooleanInfix(node, left, right)
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		return evaluateNumberInfix(node, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evaluateStringInfix(node, left, right)
	case left.Type() != right.Type():
		return newError("%d:%d:%s: runtime error: type mismatch: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type())
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type())
}
