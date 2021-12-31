package ast

import "ghostlang.org/x/ghost/token"

type Method struct {
	ExpressionNode
	Token     token.Token
	Left      ExpressionNode
	Method    ExpressionNode
	Arguments []ExpressionNode
}

func (node *Method) Accept(v Visitor) {
	v.visitMethod(node)
}
