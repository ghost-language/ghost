package optimizer

import (
	"math"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// The fold rules below mirror evaluator/number.go, evaluator/string.go,
// evaluator/boolean.go and evaluator/prefix.go exactly. Where the evaluator
// would report a runtime error the expression is deliberately left unfolded so
// that the error still happens, with its original position, at run time.

// foldInfix returns the literal an infix expression collapses to, or nil when
// it cannot be folded.
func foldInfix(node *ast.Infix) ast.ExpressionNode {
	switch left := node.Left.(type) {
	case *ast.Number:
		right, ok := node.Right.(*ast.Number)
		if !ok {
			return nil
		}

		return foldNumberInfix(node, left, right)

	case *ast.String:
		right, ok := node.Right.(*ast.String)
		if !ok {
			return nil
		}

		return foldStringInfix(node, left, right)

	case *ast.Boolean:
		right, ok := node.Right.(*ast.Boolean)
		if !ok {
			return nil
		}

		return foldBooleanInfix(node, left, right)
	}

	return nil
}

func foldNumberInfix(node *ast.Infix, left, right *ast.Number) ast.ExpressionNode {
	// Integer operands stay integral; a float on either side promotes the
	// result, matching object.Number.
	isFloat := left.IsFloat || right.IsFloat
	leftFloat, rightFloat := floatOf(left), floatOf(right)

	switch node.Operator {
	case token.PLUS:
		if isFloat {
			return floatNode(node, leftFloat+rightFloat)
		}

		return intNode(node, left.IntValue+right.IntValue)

	case token.MINUS:
		if isFloat {
			return floatNode(node, leftFloat-rightFloat)
		}

		return intNode(node, left.IntValue-right.IntValue)

	case token.STAR:
		if isFloat {
			return floatNode(node, leftFloat*rightFloat)
		}

		return intNode(node, left.IntValue*right.IntValue)

	case token.SLASH:
		// Division by zero is a runtime error, so leave it for the evaluator.
		// Division always promotes to float, matching object.Number.Div.
		if isZero(right) {
			return nil
		}

		return floatNode(node, leftFloat/rightFloat)

	case token.PERCENT:
		if isZero(right) {
			return nil
		}

		if isFloat {
			return floatNode(node, math.Mod(leftFloat, rightFloat))
		}

		return intNode(node, left.IntValue%right.IntValue)

	case token.LESS:
		if isFloat {
			return booleanNode(node, leftFloat < rightFloat)
		}

		return booleanNode(node, left.IntValue < right.IntValue)

	case token.LESSEQUAL:
		if isFloat {
			return booleanNode(node, leftFloat <= rightFloat)
		}

		return booleanNode(node, left.IntValue <= right.IntValue)

	case token.GREATER:
		if isFloat {
			return booleanNode(node, leftFloat > rightFloat)
		}

		return booleanNode(node, left.IntValue > right.IntValue)

	case token.GREATEREQUAL:
		if isFloat {
			return booleanNode(node, leftFloat >= rightFloat)
		}

		return booleanNode(node, left.IntValue >= right.IntValue)

	case token.EQUALEQUAL:
		if isFloat {
			return booleanNode(node, leftFloat == rightFloat)
		}

		return booleanNode(node, left.IntValue == right.IntValue)

	case token.BANGEQUAL:
		if isFloat {
			return booleanNode(node, leftFloat != rightFloat)
		}

		return booleanNode(node, left.IntValue != right.IntValue)
	}

	// token.DOTDOT is deliberately absent: a range builds a list object, and a
	// folded literal would have to be a shared mutable value.
	return nil
}

func foldStringInfix(node *ast.Infix, left, right *ast.String) ast.ExpressionNode {
	switch node.Operator {
	case token.PLUS:
		return &ast.String{Token: node.Token, Value: left.Value + right.Value}
	case token.LESS:
		return booleanNode(node, left.Value < right.Value)
	case token.LESSEQUAL:
		return booleanNode(node, left.Value <= right.Value)
	case token.GREATER:
		return booleanNode(node, left.Value > right.Value)
	case token.GREATEREQUAL:
		return booleanNode(node, left.Value >= right.Value)
	case token.EQUALEQUAL:
		return booleanNode(node, left.Value == right.Value)
	case token.BANGEQUAL:
		return booleanNode(node, left.Value != right.Value)
	}

	return nil
}

func foldBooleanInfix(node *ast.Infix, left, right *ast.Boolean) ast.ExpressionNode {
	// Ghost evaluates both sides of and/or before dispatching, so folding two
	// literal operands cannot skip a side effect.
	switch node.Operator {
	case token.AND:
		return booleanNode(node, left.Value && right.Value)
	case token.OR:
		return booleanNode(node, left.Value || right.Value)
	case token.EQUALEQUAL:
		return booleanNode(node, left.Value == right.Value)
	case token.BANGEQUAL:
		return booleanNode(node, left.Value != right.Value)
	}

	return nil
}

// foldPrefix returns the literal a prefix expression collapses to, or nil when
// it cannot be folded.
func foldPrefix(node *ast.Prefix) ast.ExpressionNode {
	switch node.Operator {
	case token.MINUS:
		// Negation is only defined for numbers; anything else is a runtime
		// error and is left alone.
		number, ok := node.Right.(*ast.Number)
		if !ok {
			return nil
		}

		if number.IsFloat {
			return floatNode(node, -number.FloatValue)
		}

		return intNode(node, -number.IntValue)

	case token.BANG:
		switch right := node.Right.(type) {
		case *ast.Boolean:
			return booleanNode(node, !right.Value)
		case *ast.Null:
			return booleanNode(node, true)
		case *ast.Number, *ast.String:
			// Matches the evaluator, where bang over a non-boolean,
			// non-null operand is false.
			return booleanNode(node, false)
		}
	}

	return nil
}

// =============================================================================
// Literal node constructors. Each keeps the original expression's token so that
// positions reported for the surrounding code stay meaningful.

// nodeToken recovers the token to attach to a folded literal.
func nodeToken(node ast.Node) token.Token {
	switch node := node.(type) {
	case *ast.Infix:
		return node.Token
	case *ast.Prefix:
		return node.Token
	}

	return token.Token{}
}

func intNode(node ast.Node, value int64) *ast.Number {
	return &ast.Number{Token: nodeToken(node), IntValue: value}
}

func floatNode(node ast.Node, value float64) *ast.Number {
	return &ast.Number{Token: nodeToken(node), FloatValue: value, IsFloat: true}
}

func booleanNode(node ast.Node, value bool) *ast.Boolean {
	return &ast.Boolean{Token: nodeToken(node), Value: value}
}

func floatOf(number *ast.Number) float64 {
	if number.IsFloat {
		return number.FloatValue
	}

	return float64(number.IntValue)
}

func isZero(number *ast.Number) bool {
	if number.IsFloat {
		return number.FloatValue == 0
	}

	return number.IntValue == 0
}
