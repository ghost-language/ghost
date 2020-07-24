package ast

import (
	"ghostlang.org/x/ghost/decimal"
	"ghostlang.org/x/ghost/token"
)

type NumberLiteral struct {
	Token token.Token
	Value decimal.Decimal
}

func (nl *NumberLiteral) expressionNode() {}

func (nl *NumberLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NumberLiteral) String() string {
	return nl.Token.Literal
}
