package ast

import "ghostlang.org/x/ghost/token"

type Call struct {
	ExpressionNode
	Token     token.Token
	Callee    ExpressionNode
	Arguments []ExpressionNode
}
