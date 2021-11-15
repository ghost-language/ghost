package ast

import (
	"ghostlang.org/x/ghost/token"
)

type String struct {
	ExpressionNode
	Token token.Token
	Value string
}

func (structure *String) Accept(v Visitor) {
	v.visitString(structure)
}
