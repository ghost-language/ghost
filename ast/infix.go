package ast

import "ghostlang.org/x/ghost/token"

type Infix struct {
	ExpressionNode
	Token    token.Token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (node *Infix) Accept(v Visitor) {
	v.visitInfix(node)
}
