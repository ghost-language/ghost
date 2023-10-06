package ast

import "ghostlang.org/x/ghost/token"

type Continue struct {
	ExpressionNode
	Token token.Token
}
