package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Function struct {
	ExpressionNode
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Defaults   map[string]Expression
	Body       *Block
	// Environment *object.Environment
}

func (node *Function) Accept(v Visitor) {
	v.visitFunction(node)
}
