package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Number struct {
	ExpressionNode
	Token      token.Token
	IntValue   int64
	FloatValue float64
	IsFloat    bool
}
