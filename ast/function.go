package ast

import "ghostlang.org/x/ghost/token"

type Function struct {
	Token      token.Token
	Name       string
	Parameters []*Variable
	Defaults   map[string]ExpressionNode
	Body       *Block
}
