package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateString(node *ast.String, scope *object.Scope) object.Object {
	return &object.String{Value: node.Value}
}

func evaluateStringInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.String()
	rightValue := right.String()

	switch node.Operator {
	case token.PLUS:
		return &object.String{Value: leftValue + rightValue}
	case token.LESS:
		return toBooleanValue(leftValue < rightValue)
	case token.LESSEQUAL:
		return toBooleanValue(leftValue <= rightValue)
	case token.GREATER:
		return toBooleanValue(leftValue > rightValue)
	case token.GREATEREQUAL:
		return toBooleanValue(leftValue >= rightValue)
	case token.EQUALEQUAL:
		return toBooleanValue(leftValue == rightValue)
	case token.BANGEQUAL:
		return toBooleanValue(leftValue != rightValue)
	}

	return object.NewError(fault.Type, node.Token, "cannot use `%s` between two strings", node.Operator)
}
