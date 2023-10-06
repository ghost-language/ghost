package ast

import "ghostlang.org/x/ghost/token"

type Import struct {
	ExpressionNode
	Token token.Token
	Path  *String
}
