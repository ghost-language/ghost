package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateNumber(node *ast.Number, scope *object.Scope) object.Object {
	if node.IsFloat {
		return object.NewFloat(node.FloatValue)
	}
	return object.NewInt(node.IntValue)
}

func evaluateNumberInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftNum := left.(*object.Number)
	rightNum := right.(*object.Number)

	switch node.Operator {
	case "+":
		return leftNum.Add(rightNum)
	case "-":
		return leftNum.Sub(rightNum)
	case "*":
		return leftNum.Mul(rightNum)
	case "/":
		return leftNum.Div(rightNum)
	case "%":
		return leftNum.Mod(rightNum)
	case "<":
		return toBooleanValue(leftNum.LessThan(rightNum))
	case "<=":
		return toBooleanValue(leftNum.LessThanOrEqual(rightNum))
	case ">":
		return toBooleanValue(leftNum.GreaterThan(rightNum))
	case ">=":
		return toBooleanValue(leftNum.GreaterThanOrEqual(rightNum))
	case "==":
		return toBooleanValue(leftNum.Equal(rightNum))
	case "!=":
		return toBooleanValue(!leftNum.Equal(rightNum))
	case "..":
		start := leftNum.Int64()
		end := rightNum.Int64()

		if start > end {
			return &object.List{Elements: []object.Object{}}
		}

		numbers := make([]object.Object, 0, end-start+1)

		for i := start; i <= end; i++ {
			numbers = append(numbers, object.NewInt(i))
		}

		return &object.List{Elements: numbers}
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, right.Type(), node.Operator, left.Type())
}
