package ast

import (
	"ghostlang.org/x/ghost/token"
)

type String struct {
	ExpressionNode
	Token token.Token
	Value string
}
