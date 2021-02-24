package ast

import "ghostlang.org/x/ghost/token"

type Variable struct {
	ExpressionNode
	Name token.Token
}
