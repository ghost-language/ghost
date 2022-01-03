package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number, env *object.Environment) (object.Object, bool) {
	return &object.Number{Value: node.Value}, true
}

func evaluateNumberInfix(node *ast.Infix, right object.Object, left object.Object) (object.Object, bool) {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch node.Operator {
	case "+":
		return &object.Number{Value: rightValue.Add(leftValue)}, true
	case "-":
		return &object.Number{Value: rightValue.Sub(leftValue)}, true
	case "*":
		return &object.Number{Value: rightValue.Mul(leftValue)}, true
	case "/":
		return &object.Number{Value: rightValue.Div(leftValue)}, true
	case "%":
		return &object.Number{Value: rightValue.Mod(leftValue)}, true
	case "<":
		return toBooleanValue(rightValue.LessThan(leftValue)), true
	case "<=":
		return toBooleanValue(rightValue.LessThanOrEqual(leftValue)), true
	case ">":
		return toBooleanValue(rightValue.GreaterThan(leftValue)), true
	case ">=":
		return toBooleanValue(rightValue.GreaterThanOrEqual(leftValue)), true
	case "==":
		return toBooleanValue(rightValue.Equal(leftValue)), true
	case "!=":
		return toBooleanValue(!rightValue.Equal(leftValue)), true
	default:
		err := newError("unknown operator: %s %s %s", right.Type(), node.Operator, left.Type())

		return err, false
	}
}
