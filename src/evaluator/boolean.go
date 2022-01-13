package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBoolean(node *ast.Boolean, scope *object.Scope) object.Object {
	return toBooleanValue(node.Value)
}

func evaluateBooleanInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch node.Operator {
	case "and":
		return toBooleanValue(leftValue && rightValue)
	case "or":
		return toBooleanValue(leftValue || rightValue)
	case "==":
		return toBooleanValue(leftValue == rightValue)
	case "!=":
		return toBooleanValue(leftValue != rightValue)
	}

	return newError("%d:%d: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, right.Type(), node.Operator, left.Type())
}
