package ast

import "ghostlang.org/x/ghost/token"

type Identifier struct {
	ExpressionNode
	Token token.Token
	Value string
}

func (structure *Identifier) Accept(v Visitor) {
	v.visitIdentifier(structure)
}
