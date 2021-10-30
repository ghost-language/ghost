package ast

import (
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

type Number struct {
	Token token.Token
	Value decimal.Decimal
}

func (number *Number) expressionNode() {}
