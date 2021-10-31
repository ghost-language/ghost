package ast

import "ghostlang.org/x/ghost/token"

type Null struct {
	ExpressionNode
	Token token.Token
}
