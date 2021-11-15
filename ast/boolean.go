package ast

import "ghostlang.org/x/ghost/token"

type Boolean struct {
	ExpressionNode
	Token token.Token
	Value bool
}

func (structure *Boolean) Accept(v Visitor) {
	v.visitBoolean(structure)
}
