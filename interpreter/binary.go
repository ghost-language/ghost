package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateBinary(node *ast.Binary) object.Object {
	left := Evaluate(node.Left)
	right := Evaluate(node.Right)

	switch node.Operator.Type {
	case token.MINUS:
		value := left.(*object.Number).Value.Sub(right.(*object.Number).Value)
		return &object.Number{Value: value}
	case token.PLUS:
		value := left.(*object.Number).Value.Add(right.(*object.Number).Value)
		return &object.Number{Value: value}
	case token.SLASH:
		value := left.(*object.Number).Value.Div(right.(*object.Number).Value)
		return &object.Number{Value: value}
	case token.STAR:
		value := left.(*object.Number).Value.Mul(right.(*object.Number).Value)
		return &object.Number{Value: value}
	case token.GREATER:
		value := left.(*object.Number).Value.GreaterThan(right.(*object.Number).Value)
		return &object.Boolean{Value: value}
	case token.GREATEREQUAL:
		value := left.(*object.Number).Value.GreaterThanOrEqual(right.(*object.Number).Value)
		return &object.Boolean{Value: value}
	case token.LESS:
		value := left.(*object.Number).Value.LessThan(right.(*object.Number).Value)
		return &object.Boolean{Value: value}
	case token.LESSEQUAL:
		value := left.(*object.Number).Value.LessThanOrEqual(right.(*object.Number).Value)
		return &object.Boolean{Value: value}
	case token.EQUALEQUAL:
		value := left.(*object.Number).Value.Equal(right.(*object.Number).Value)
		return &object.Boolean{Value: value}
	}

	panic("Fatal error")
}
