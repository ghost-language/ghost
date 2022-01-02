package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Function struct {
	ExpressionNode
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Defaults   map[string]ExpressionNode
	Body       *Block
}
