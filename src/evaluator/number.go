package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number, scope *object.Scope) object.Object {
	return &object.Number{Value: node.Value}
}

func evaluateNumberInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Number).Value
	rightValue := right.(*object.Number).Value

	switch node.Operator {
	case "+":
		return &object.Number{Value: leftValue.Add(rightValue)}
	case "-":
		return &object.Number{Value: leftValue.Sub(rightValue)}
	case "*":
		return &object.Number{Value: leftValue.Mul(rightValue)}
	case "/":
		return &object.Number{Value: leftValue.Div(rightValue)}
	case "%":
		return &object.Number{Value: leftValue.Mod(rightValue)}
	case "<":
		return toBooleanValue(leftValue.LessThan(rightValue))
	case "<=":
		return toBooleanValue(leftValue.LessThanOrEqual(rightValue))
	case ">":
		return toBooleanValue(leftValue.GreaterThan(rightValue))
	case ">=":
		return toBooleanValue(leftValue.GreaterThanOrEqual(rightValue))
	case "==":
		return toBooleanValue(leftValue.Equal(rightValue))
	case "!=":
		return toBooleanValue(!leftValue.Equal(rightValue))
	}

	return newError("%d:%d: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, right.Type(), node.Operator, left.Type())
}
