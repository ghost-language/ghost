package ast

import "ghostlang.org/x/ghost/token"

type Identifier struct {
	ExpressionNode
	Token token.Token
	Value string
}
