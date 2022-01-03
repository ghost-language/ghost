package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number, env *object.Environment) object.Object {
	return &object.Number{Value: node.Value}
}

func evaluateNumberInfix(node *ast.Infix, right object.Object, left object.Object) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch node.Operator {
	case "+":
		return &object.Number{Value: rightValue.Add(leftValue)}
	case "-":
		return &object.Number{Value: rightValue.Sub(leftValue)}
	case "*":
		return &object.Number{Value: rightValue.Mul(leftValue)}
	case "/":
		return &object.Number{Value: rightValue.Div(leftValue)}
	case "%":
		return &object.Number{Value: rightValue.Mod(leftValue)}
	case "<":
		return toBooleanValue(rightValue.LessThan(leftValue))
	case "<=":
		return toBooleanValue(rightValue.LessThanOrEqual(leftValue))
	case ">":
		return toBooleanValue(rightValue.GreaterThan(leftValue))
	case ">=":
		return toBooleanValue(rightValue.GreaterThanOrEqual(leftValue))
	case "==":
		return toBooleanValue(rightValue.Equal(leftValue))
	case "!=":
		return toBooleanValue(!rightValue.Equal(leftValue))
	default:
		return newError("unknown operator: %s %s %s", right.Type(), node.Operator, left.Type())
	}
}
