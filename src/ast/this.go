package ast

import "ghostlang.org/x/ghost/token"

type This struct {
	ExpressionNode
	Token token.Token
}
