package ast

import "ghostlang.org/x/ghost/token"

type For struct {
	ExpressionNode
	Token       token.Token
	Identifier  *Identifier
	Initializer StatementNode
	Condition   ExpressionNode
	Increment   StatementNode
	Block       *Block
}

func (node *For) Accept(v Visitor) {
	v.visitFor(node)
}
