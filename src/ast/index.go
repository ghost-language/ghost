package ast

import "ghostlang.org/x/ghost/token"

type Index struct {
	ExpressionNode
	Token token.Token
	Left  ExpressionNode
	Index ExpressionNode
}

func (node *Index) Accept(v Visitor) {
	v.visitIndex(node)
}
