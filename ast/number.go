package ast

import (
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

type Number struct {
	ExpressionNode
	Token token.Token
	Value decimal.Decimal
}

func (node *Number) Accept(v Visitor) {
	v.visitNumber(node)
}
