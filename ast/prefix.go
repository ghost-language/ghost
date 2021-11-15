package ast

import "ghostlang.org/x/ghost/token"

type Prefix struct {
	ExpressionNode
	Token    token.Token
	Operator string
	Right    ExpressionNode
}

func (structure *Prefix) Accept(v Visitor) {
	v.visitPrefix(structure)
}
