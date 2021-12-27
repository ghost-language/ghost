package ast

import "ghostlang.org/x/ghost/token"

type List struct {
	ExpressionNode
	Token    token.Token
	Elements []ExpressionNode
}

func (node *List) Accept(v Visitor) {
	v.visitList(node)
}
