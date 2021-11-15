package ast

import "ghostlang.org/x/ghost/token"

type Null struct {
	ExpressionNode
	Token token.Token
}

func (structure *Null) Accept(v Visitor) {
	v.visitNull(structure)
}
