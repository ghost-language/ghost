package ast

import "ghostlang.org/x/ghost/token"

type Function struct {
	Token      token.Token
	Name       string
	Parameters []*Identifier
	Defaults   map[string]ExpressionNode
	Body       *Block
}
