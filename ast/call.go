package ast

import "ghostlang.org/x/ghost/token"

type Call struct {
	Callee    ExpressionNode
	Paren     token.Token
	Arguments []ExpressionNode
}
