package ast

import "ghostlang.org/x/ghost/token"

type Use struct {
	ExpressionNode
	Token  token.Token
	Traits []*Identifier
}
