package ast

import "github.com/shopspring/decimal"

// Number structures are for number expressions.
type Number struct {
	ExpressionNode
	Value decimal.Decimal
}
