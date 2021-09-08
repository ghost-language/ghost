package ast

import "ghostlang.org/x/ghost/token"

type Get struct {
	ExpressionNode
	Name token.Token
	Expression ExpressionNode
}