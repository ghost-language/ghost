package ast

import "ghostlang.org/x/ghost/token"

type Break struct {
	ExpressionNode
	Token token.Token
}
