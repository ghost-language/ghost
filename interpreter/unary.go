package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateUnary(node *ast.Unary) object.Object {
	right := Evaluate(node.Right)

	if node.Operator.Type == token.MINUS && right.Type() == object.NUMBER {
		value := right.(*object.Number).Value.Neg()
		return &object.Number{Value: value}
	} else if node.Operator.Type == token.BANG {
		return helper.NativeBooleanToObject(!helper.IsTruthy(right))
	}

	panic("Unexpected error.")
}
