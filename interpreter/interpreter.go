package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(expression ast.ExpressionNode) interface{} {
	switch node := expression.(type) {
	case *ast.Binary:
		left := Evaluate(node.Left)
		right := Evaluate(node.Right)

		switch node.Operator.Type {
		case token.MINUS:
			return left.(decimal.Decimal).Sub(right.(decimal.Decimal))
		case token.PLUS:
			switch left := left.(type) {
			case decimal.Decimal:
				return left.Add(right.(decimal.Decimal))
			case string:
				return left + right.(string)
			}
		case token.SLASH:
			return left.(decimal.Decimal).Div(right.(decimal.Decimal))
		case token.STAR:
			return left.(decimal.Decimal).Mul(right.(decimal.Decimal))
		case token.GREATER:
			return left.(decimal.Decimal).GreaterThan(right.(decimal.Decimal))
		}
	case *ast.Grouping:
		return Evaluate(node.Expression)
	case *ast.Literal:
		return node.Value
	}

	panic("Fatal error")
}
