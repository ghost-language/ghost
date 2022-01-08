package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

func evaluateInfix(node *ast.Infix, env *object.Environment) object.Object {
	left := Evaluate(node.Left, env)

	if isError(left) {
		return left
	}

	right := Evaluate(node.Right, env)

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		return evaluateNumberInfix(node, left, right)
	case left.Type() == object.STRING:
		return evaluateStringInfix(node, left, right)
	case left.Type() != right.Type():
		return newError("%d:%d: runtime error: type mismatch: %s %s %s", node.Token.Line, node.Token.Column, left.Type(), node.Operator, right.Type())
	default:
		log.Debug("default reporting")
		return newError("%d:%d: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, left.Type(), node.Operator, right.Type())
	}
}
