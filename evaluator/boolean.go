package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateBoolean(node *ast.Boolean, scope *object.Scope) object.Object {
	return toBooleanValue(node.Value)
}

func evaluateBooleanInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch node.Operator {
	case token.AND:
		return toBooleanValue(leftValue && rightValue)
	case token.OR:
		return toBooleanValue(leftValue || rightValue)
	case token.EQUALEQUAL:
		return toBooleanValue(leftValue == rightValue)
	case token.BANGEQUAL:
		return toBooleanValue(leftValue != rightValue)
	}

	return object.NewError(fault.Type, node.Token, "cannot use `%s` between two booleans", node.Operator)
}
