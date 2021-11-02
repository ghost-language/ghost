package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number) (object.Object, bool) {
	return &object.Number{Value: node.Value}, true
}

func evaluateNumberInfix(node *ast.Infix, right object.Object, left object.Object) (object.Object, bool) {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch node.Operator {
	case "+":
		return &object.Number{Value: leftValue.Add(rightValue)}, true
	case "-":
		return &object.Number{Value: leftValue.Sub(rightValue)}, true
	case "*":
		return &object.Number{Value: leftValue.Mul(rightValue)}, true
	case "/":
		return &object.Number{Value: leftValue.Div(rightValue)}, true
	case "%":
		return &object.Number{Value: leftValue.Mod(rightValue)}, true
	case "<":
		return &object.Boolean{Value: leftValue.LessThan(rightValue)}, true
	case "<=":
		return &object.Boolean{Value: leftValue.LessThanOrEqual(rightValue)}, true
	case ">":
		return &object.Boolean{Value: leftValue.GreaterThan(rightValue)}, true
	case ">=":
		return &object.Boolean{Value: leftValue.GreaterThanOrEqual(rightValue)}, true
	case "==":
		return &object.Boolean{Value: leftValue.Equal(rightValue)}, true
	case "!=":
		return &object.Boolean{Value: !leftValue.Equal(rightValue)}, true
	default:
		return nil, false
	}
}
