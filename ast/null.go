package ast

import "ghostlang.org/x/ghost/token"

type Null struct {
	ExpressionNode
	Token token.Token
}

func (node *Null) Accept(v Visitor) {
	v.visitNull(node)
}
