package ast

import "ghostlang.org/x/ghost/token"

type Identifier struct {
	ExpressionNode
	Name token.Token
}
