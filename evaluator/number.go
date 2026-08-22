package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
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
	case token.PLUS:
		return leftNum.Add(rightNum)
	case token.MINUS:
		return leftNum.Sub(rightNum)
	case token.STAR:
		return leftNum.Mul(rightNum)
	case token.SLASH:
		if rightNum.IsZero() {
			return object.NewError(fault.Value, node.Token, "cannot divide by zero")
		}

		return leftNum.Div(rightNum)
	case token.PERCENT:
		if rightNum.IsZero() {
			return object.NewError(fault.Value, node.Token, "cannot take the remainder of a division by zero")
		}

		return leftNum.Mod(rightNum)
	case token.LESS:
		return toBooleanValue(leftNum.LessThan(rightNum))
	case token.LESSEQUAL:
		return toBooleanValue(leftNum.LessThanOrEqual(rightNum))
	case token.GREATER:
		return toBooleanValue(leftNum.GreaterThan(rightNum))
	case token.GREATEREQUAL:
		return toBooleanValue(leftNum.GreaterThanOrEqual(rightNum))
	case token.EQUALEQUAL:
		return toBooleanValue(leftNum.Equal(rightNum))
	case token.BANGEQUAL:
		return toBooleanValue(!leftNum.Equal(rightNum))
	case token.DOTDOT:
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

	return object.NewError(fault.Type, node.Token, "cannot use `%s` between two numbers", node.Operator)
}
