package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/helper"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

func evaluateBinary(node *ast.Binary) object.Object {
	left := Evaluate(node.Left)
	right := Evaluate(node.Right)

	switch node.Operator.Type {
	case token.MINUS:
		value := left.(*object.Number).Value.Sub(right.(*object.Number).Value)
		return &object.Number{Value: value}
	case token.PLUS:
		switch left.(type) {
		case *object.Number:
			switch right.(type) {
			case *object.String:
				number, _ := decimal.NewFromString(right.(*object.String).Value)
				value := left.(*object.Number).Value.Add(number)
				return &object.Number{Value: value}
			case *object.Number:
				value := left.(*object.Number).Value.Add(right.(*object.Number).Value)
				return &object.Number{Value: value}
			}

		case *object.String:
			switch right.(type) {
			case *object.String:
				value := left.(*object.String).Value + right.(*object.String).Value
				return &object.String{Value: value}
			case *object.Number:
				value := left.(*object.String).Value + right.(*object.Number).String()
				return &object.String{Value: value}
			}
		}
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
		return helper.NativeBooleanToObject(helper.IsEqual(left, right))
	}

	panic("Fatal error")
}
